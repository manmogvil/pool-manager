package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
)

type Store interface {
	GetPendingDraws() ([]models.Draw, error)
	GetLotteryGameByID(id int) (models.LotteryGame, error)
	GetDrawsByGame(gameID int) ([]models.Draw, error)
	UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error)
	UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error)
}

type Scheduler struct {
	store     Store
	api       services.LoteriaAPI
	cronExpr  string
	cron      *cron.Cron
	entryID   cron.EntryID
	mu        sync.Mutex
	running   bool
}

func New(store Store, api services.LoteriaAPI, cronExpr string) *Scheduler {
	return &Scheduler{
		store:    store,
		api:      api,
		cronExpr: cronExpr,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	log.Printf("[Scheduler] Starting with cron: %s", s.cronExpr)

	s.cron = cron.New()

	var err error
	s.entryID, err = s.cron.AddFunc(s.cronExpr, func() {
		s.execute()
	})
	if err != nil {
		log.Printf("[Scheduler] Invalid cron expression '%s': %v", s.cronExpr, err)
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return
	}

	s.cron.Start()

	<-ctx.Done()

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	ctxStop := s.cron.Stop()
	<-ctxStop.Done()

	log.Println("[Scheduler] Stopped")
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running && s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.running = false
	}
}

func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Scheduler) execute() {
	log.Println("[Scheduler] Checking for pending draws...")

	draws, err := s.store.GetPendingDraws()
	if err != nil {
		log.Printf("[Scheduler] Error fetching pending draws: %v", err)
		return
	}

	if len(draws) == 0 {
		log.Println("[Scheduler] No pending draws found")
		return
	}

	log.Printf("[Scheduler] Found %d pending draws", len(draws))

	for _, draw := range draws {
		game, err := s.store.GetLotteryGameByID(draw.GameID)
		if err != nil {
			log.Printf("[Scheduler] Error fetching game %d: %v", draw.GameID, err)
			continue
		}

		slug := services.MapNameToGameSlug(game.Name)
		if slug == "" {
			log.Printf("[Scheduler] No slug found for game: %s", game.Name)
			continue
		}

		from := draw.DrawDate.Format("2006-01-02")
		to := draw.DrawDate.Format("2006-01-02")

		result, err := s.api.GetResults(slug, from, to)
		if err != nil {
			log.Printf("[Scheduler] Error fetching results for %s: %v", game.Name, err)
			continue
		}

		existingDraws, _ := s.store.GetDrawsByGame(draw.GameID)
		existingByDate := make(map[string]models.Draw)
		for _, d := range existingDraws {
			existingByDate[d.DrawDate.Format("2006-01-02")] = d
		}

		for _, item := range result.Data {
			drawDate, err := time.Parse("2006-01-02", item.DrawDate)
			if err != nil {
				continue
			}

			resultNumbers := services.FormatCombination(item.Combination)
			resultStars := services.FormatStars(item.ResultData.Estrellas)
			dateKey := drawDate.Format("2006-01-02")

			if existing, ok := existingByDate[dateKey]; ok {
				_, err := s.store.UpdateDrawResults(existing.ID, &resultNumbers, &resultStars)
				if err != nil {
					log.Printf("[Scheduler] Error updating draw %d: %v", existing.ID, err)
					continue
				}
				if item.DrawId != "" {
					_, err = s.store.UpdateDrawDrawIDAPI(existing.ID, item.DrawId)
					if err != nil {
						log.Printf("[Scheduler] Error updating draw API ID %d: %v", existing.ID, err)
					}
				}
				log.Printf("[Scheduler] Updated draw %d for %s on %s", existing.ID, game.Name, dateKey)
			}
		}
	}

	log.Println("[Scheduler] Execution completed")
}
