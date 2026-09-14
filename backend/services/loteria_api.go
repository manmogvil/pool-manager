package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type LoteriaAPIClient struct {
	BaseURL    string
	APIKeys    []string
	currentKey int
	HTTPClient *http.Client
}

type LoteriaAPI interface {
	CheckCombination(gameSlug string, numbers string, extraNumbers string, drawId string) (*CheckCombinationResponse, error)
	GetResults(gameSlug string, fromDate string, toDate string) (*ResultsListResponse, error)
}

type LoteriaAPIGame struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type LoteriaAPIResult struct {
	Estrellas []int `json:"estrellas"`
}

type LoteriaAPIPrize struct {
	Category       interface{} `json:"category"`
	Winners        int         `json:"winners"`
	FormattedPrize string      `json:"formattedPrize"`
}

type CheckCombinationResponse struct {
	Success bool                 `json:"success"`
	Data    CheckCombinationData `json:"data"`
	Error   string               `json:"error,omitempty"`
}

type CheckCombinationData struct {
	HasResults          bool           `json:"hasResults"`
	CheckedNumbers      []int          `json:"checkedNumbers"`
	CheckedExtraNumbers []int          `json:"checkedExtraNumbers"`
	Game                LoteriaAPIGame `json:"game"`
	DrawDate            string         `json:"drawDate"`
	DrawId              string         `json:"drawId"`
	WinningCombination  []int          `json:"winningCombination"`
	WinningExtraNumbers []int          `json:"winningExtraNumbers"`
	MainNumbersMatched  int            `json:"mainNumbersMatched"`
	ExtraNumbersMatched int            `json:"extraNumbersMatched"`
	MatchedNumbers      []int          `json:"matchedNumbers"`
	MatchedExtraNumbers []int          `json:"matchedExtraNumbers"`
	IsWinner            bool           `json:"isWinner"`
}

func NewLoteriaAPIClient() *LoteriaAPIClient {
	apiKeys := []string{}
	if keys := os.Getenv("LOTERIA_API_KEY"); keys != "" {
		for _, k := range strings.Split(keys, ",") {
			if trimmed := strings.TrimSpace(k); trimmed != "" {
				apiKeys = append(apiKeys, trimmed)
			}
		}
	}
	return &LoteriaAPIClient{
		BaseURL: "https://api.loteriasapi.com/api/v1",
		APIKeys: apiKeys,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *LoteriaAPIClient) getCurrentKey() string {
	if len(c.APIKeys) == 0 {
		return ""
	}
	return c.APIKeys[c.currentKey]
}

func (c *LoteriaAPIClient) switchToNextKey() bool {
	if c.currentKey+1 < len(c.APIKeys) {
		c.currentKey++
		return true
	}
	return false
}

func (c *LoteriaAPIClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("X-API-Key", c.getCurrentKey())
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode == 429 || resp.StatusCode == 403 {
		if c.switchToNextKey() {
			req.Header.Set("X-API-Key", c.getCurrentKey())
			resp2, err2 := c.HTTPClient.Do(req)
			if err2 != nil {
				return nil, fmt.Errorf("error making request with fallback key: %w", err2)
			}
			defer resp2.Body.Close()
			body, err = io.ReadAll(resp2.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading fallback response: %w", err)
			}
			return body, nil
		}
	}

	return body, nil
}

func (c *LoteriaAPIClient) CheckCombination(gameSlug string, numbers string, extraNumbers string, drawId string) (*CheckCombinationResponse, error) {
	baseURL := fmt.Sprintf("%s/results/%s/check", c.BaseURL, gameSlug)

	params := url.Values{}
	params.Set("numbers", numbers)
	if extraNumbers != "" {
		params.Set("extraNumbers", extraNumbers)
	}
	if drawId != "" {
		params.Set("drawId", drawId)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result CheckCombinationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s", result.Error)
	}

	return &result, nil
}

type ResultsListItem struct {
	Game             LoteriaAPIGame    `json:"game"`
	DrawId           string            `json:"drawId"`
	DrawDate         string            `json:"drawDate"`
	DayOfWeek        string            `json:"dayOfWeek"`
	Status           string            `json:"status"`
	Combination      []int             `json:"combination"`
	ResultData       LoteriaAPIResult  `json:"resultData"`
	JackpotFormatted string            `json:"jackpotFormatted"`
	Prizes           []LoteriaAPIPrize `json:"prizes"`
}

type ResultsListMeta struct {
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

type ResultsListResponse struct {
	Success bool              `json:"success"`
	Data    []ResultsListItem `json:"data"`
	Meta    ResultsListMeta   `json:"meta"`
	Error   string            `json:"error,omitempty"`
}

func (c *LoteriaAPIClient) GetResults(gameSlug string, fromDate string, toDate string) (*ResultsListResponse, error) {
	baseURL := fmt.Sprintf("%s/results", c.BaseURL)

	params := url.Values{}
	params.Set("gameType", gameSlug)
	if fromDate != "" {
		params.Set("from", fromDate)
	}
	if toDate != "" {
		params.Set("to", toDate)
	}
	params.Set("sort", "asc")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result ResultsListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s", result.Error)
	}

	return &result, nil
}

var gameNamesBySlug = map[string]string{
	"euromillones": "EuroMillones",
	"primitiva":    "La Primitiva",
	"bonoloto":     "Bonoloto",
	"gordo":        "El Gordo de la Primitiva",
	"nacional":     "Lotería Nacional",
	"eurodreams":   "Eurodreams",
	"quiniela":     "La Quiniela",
	"quinigol":     "El Quinigol",
	"lototurf":     "Lototurf",
	"quintuple":    "Quíntuple Plus",
}

func MapGameSlugToName(slug string) string {
	if name, ok := gameNamesBySlug[slug]; ok {
		return name
	}
	return slug
}

func MapNameToGameSlug(name string) string {
	for slug, gameName := range gameNamesBySlug {
		if gameName == name {
			return slug
		}
	}
	return ""
}

func formatInts(values []int) string {
	if len(values) == 0 {
		return ""
	}

	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.Itoa(v)
	}

	return strings.Join(parts, ",")
}

func FormatCombination(numbers []int) string {
	return formatInts(numbers)
}

func FormatStars(stars []int) string {
	return formatInts(stars)
}
