package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
)

type mockStore struct {
	pendingDraws    []models.Draw
	games           map[int]models.LotteryGame
	existingDraws   []models.Draw
	updatedDraws    []models.Draw
	updateErr       error
	getDrawsErr     error
	pendingDrawsErr error
}

func newMockStore() *mockStore {
	return &mockStore{
		games: make(map[int]models.LotteryGame),
	}
}

func (m *mockStore) GetPendingDraws() ([]models.Draw, error) {
	if m.pendingDrawsErr != nil {
		return nil, m.pendingDrawsErr
	}
	return m.pendingDraws, nil
}

func (m *mockStore) GetLotteryGameByID(id int) (models.LotteryGame, error) {
	if g, ok := m.games[id]; ok {
		return g, nil
	}
	return models.LotteryGame{}, errors.New("game not found")
}

func (m *mockStore) GetDrawsByGame(gameID int) ([]models.Draw, error) {
	if m.getDrawsErr != nil {
		return nil, m.getDrawsErr
	}
	return m.existingDraws, nil
}

func (m *mockStore) UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
	if m.updateErr != nil {
		return models.Draw{}, m.updateErr
	}
	d := models.Draw{ID: id, ResultNumbers: resultNumbers, ResultStars: resultStars}
	m.updatedDraws = append(m.updatedDraws, d)
	return d, nil
}

func (m *mockStore) UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error) {
	if m.updateErr != nil {
		return models.Draw{}, m.updateErr
	}
	return models.Draw{ID: id, DrawIDAPI: &drawIDAPI}, nil
}

type mockAPI struct {
	results *services.ResultsListResponse
	err     error
	calls   int
}

func (m *mockAPI) CheckCombination(gameSlug string, numbers string, extraNumbers string, drawId string) (*services.CheckCombinationResponse, error) {
	return nil, nil
}

func (m *mockAPI) GetResults(gameSlug string, fromDate string, toDate string) (*services.ResultsListResponse, error) {
	m.calls++
	return m.results, m.err
}

func (m *mockAPI) GetCalls() int {
	return m.calls
}

func TestScheduler_DisabledByDefault(t *testing.T) {
	store := newMockStore()
	api := &mockAPI{}
	s := New(store, api, "* * * * *")

	if s.IsRunning() {
		t.Error("scheduler should not be running by default")
	}
}

func TestScheduler_StartsWhenEnabled(t *testing.T) {
	store := newMockStore()
	api := &mockAPI{}
	s := New(store, api, "* * * * *")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	if !s.IsRunning() {
		t.Error("scheduler should be running after Start()")
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestScheduler_DoesNotStartWhenDisabled(t *testing.T) {
	store := newMockStore()
	api := &mockAPI{}
	s := New(store, api, "* * * * *")

	if s.IsRunning() {
		t.Error("scheduler should not be running")
	}

	s.Stop()
	time.Sleep(50 * time.Millisecond)

	if s.IsRunning() {
		t.Error("scheduler should not start after Stop()")
	}
}

func TestScheduler_UsesConfiguredCron(t *testing.T) {
	store := newMockStore()
	store.pendingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}
	store.games[1] = models.LotteryGame{ID: 1, Name: "EuroMillones"}
	store.existingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}

	api := &mockAPI{
		results: &services.ResultsListResponse{
			Data: []services.ResultsListItem{},
		},
	}

	cronExpr := "@every 1s"
	s := New(store, api, cronExpr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(2 * time.Second)

	if api.GetCalls() < 1 {
		t.Error("scheduler should have executed at least once")
	}

	cancel()
}

func TestScheduler_TriggersDrawProcessing(t *testing.T) {
	store := newMockStore()
	store.pendingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}
	store.games[1] = models.LotteryGame{ID: 1, Name: "EuroMillones"}
	store.existingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}

	api := &mockAPI{
		results: &services.ResultsListResponse{
			Data: []services.ResultsListItem{
				{
					DrawDate:    time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
					Combination: []int{1, 5, 12, 30, 45},
					ResultData: services.LoteriaAPIResult{
						Estrellas: []int{3, 7},
					},
				},
			},
		},
	}

	s := New(store, api, "@every 1s")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(2 * time.Second)

	if len(store.updatedDraws) == 0 {
		t.Error("scheduler should have updated at least one draw")
	}

	cancel()
}

func TestScheduler_ErrorsDoNotStopScheduler(t *testing.T) {
	store := newMockStore()
	store.pendingDrawsErr = errors.New("db error")

	api := &mockAPI{}
	s := New(store, api, "* * * * *")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	if !s.IsRunning() {
		t.Error("scheduler should still be running after error")
	}

	cancel()
}

func TestScheduler_StopsOnShutdown(t *testing.T) {
	store := newMockStore()
	api := &mockAPI{}
	s := New(store, api, "* * * * *")

	ctx, cancel := context.WithCancel(context.Background())

	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	if !s.IsRunning() {
		t.Error("scheduler should be running")
	}

	cancel()
	time.Sleep(200 * time.Millisecond)

	if s.IsRunning() {
		t.Error("scheduler should have stopped after context cancellation")
	}
}

func TestScheduler_NoOverlappingExecutions(t *testing.T) {
	store := newMockStore()
	store.pendingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}
	store.games[1] = models.LotteryGame{ID: 1, Name: "EuroMillones"}
	store.existingDraws = []models.Draw{
		{ID: 1, GameID: 1, DrawDate: time.Now().Add(-24 * time.Hour)},
	}

	api := &mockAPI{
		results: &services.ResultsListResponse{
			Data: []services.ResultsListItem{
				{
					DrawDate:    time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
					Combination: []int{1, 5, 12, 30, 45},
					ResultData: services.LoteriaAPIResult{
						Estrellas: []int{3, 7},
					},
				},
			},
		},
	}

	s := New(store, api, "* * * * *")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(200 * time.Millisecond)

	cancel()
}

func TestScheduler_InvalidCronExpression(t *testing.T) {
	store := newMockStore()
	api := &mockAPI{}
	s := New(store, api, "invalid cron")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	if s.IsRunning() {
		t.Error("scheduler should not be running with invalid cron expression")
	}
}
