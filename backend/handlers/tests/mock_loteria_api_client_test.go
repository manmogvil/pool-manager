package tests

import (
	"lottery-pool-manager/services"
)

type mockLoteriaAPI struct {
	checkCombinationFn  func(gameSlug string, numbers string, extraNumbers string, drawId string) (*services.CheckCombinationResponse, error)
	getResultsFn        func(gameSlug string, fromDate string, toDate string) (*services.ResultsListResponse, error)
}

func (m *mockLoteriaAPI) CheckCombination(gameSlug string, numbers string, extraNumbers string, drawId string) (*services.CheckCombinationResponse, error) {
	return m.checkCombinationFn(gameSlug, numbers, extraNumbers, drawId)
}
func (m *mockLoteriaAPI) GetResults(gameSlug string, fromDate string, toDate string) (*services.ResultsListResponse, error) {
	return m.getResultsFn(gameSlug, fromDate, toDate)
}
