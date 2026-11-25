package transport

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/copernicus-ml/web-service/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getStockStoreHandler() (*mockStockStoreRepo, *StockHandler) {
	mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockStore := new(mockStockStoreRepo)
	return mockStore, &StockHandler{
		logger:     mockLogger,
		stockStore: mockStore,
	}
}

// test validate ticker
func TestValidateTicker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		inputTicker string
		expectError bool
	}{
		{
			name:        "success",
			inputTicker: "AAPL",
			expectError: false,
		},
		{
			name:        "invalid ticker",
			inputTicker: "invalid",
			expectError: true,
		},
	}
	mockStore, storeHandler := getStockStoreHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			input := strings.ToUpper(tt.inputTicker)
			if !tt.expectError {
				mockStore.On("ValidateSymbol", input).Return(true).Once()
				err := storeHandler.ValidateSymbol(input)
				assert.Nil(t, err)
			} else {
				mockStore.On("ValidateSymbol", input).Return(false).Once()
				err := storeHandler.ValidateSymbol(strings.ToUpper(input))
				assert.Error(t, err)
				assert.Contains(t, err.Error(), strings.ToUpper(input))
			}
			mockStore.AssertExpectations(t)
		})
	}
}

// test validate market scope
func TestValidateMarket(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		inputScope  string
		expectError bool
	}{
		{
			name:        "success",
			inputScope:  StockPct,
			expectError: false,
		},
		{
			name:        "success with uppercase scope",
			inputScope:  "INDUSTRYPCT",
			expectError: false,
		},
		{
			name:        "invalid market scope",
			inputScope:  "invalid",
			expectError: true,
		},
	}
	_, storeHandler := getStockStoreHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			err := storeHandler.ValidateMarket(tt.inputScope)
			if !tt.expectError {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), strings.ToLower(tt.inputScope))
			}
		})
	}
}

// test validate data window
func TestValidateWindow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		inputWindow string
		expectError bool
	}{
		{
			name:        "success",
			inputWindow: "1year",
			expectError: false,
		},
		{
			name:        "invalid window",
			inputWindow: "invalid",
			expectError: true,
		},
	}
	_, storeHandler := getStockStoreHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			err := storeHandler.ValidateWindow(tt.inputWindow)
			if !tt.expectError {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.inputWindow)
			}
		})
	}
}

// test validate feature group
func TestValidateFeatureGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		inputFeature string
		expectError  bool
	}{
		{
			name:         "success",
			inputFeature: PriceMomentum,
			expectError:  false,
		},
		{
			name:         "success with uppercase feature group",
			inputFeature: "MACD",
			expectError:  false,
		},
		{
			name:         "invalid feature group",
			inputFeature: "invalid",
			expectError:  true,
		},
	}
	_, storeHandler := getStockStoreHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			err := storeHandler.ValidateFeatureGroup(tt.inputFeature)
			if !tt.expectError {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), strings.ToLower(tt.inputFeature))
			}
		})
	}
}

// test invalidate cache response
func TestInvalidateCache(t *testing.T) {
	t.Parallel()
	mockStore, storeHandler := getStockStoreHandler()
	mockStore.On("InvalidateCache").Return().Once()

	response := storeHandler.InvalidateCache()

	assert.Equal(t, "LRU Cache cleared.", response.Message)
	mockStore.AssertExpectations(t)
}

// test get stocks' day percent changes
func TestGetStocksPctChanges(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		pctChanges repositories.StockPctChangeMap
		storeError error
		want       *StocksPctChangesResponse
	}{
		{
			name: "success with an unknown symbol",
			pctChanges: repositories.StockPctChangeMap{
				"AAPL": {PctChange1D: 1.25, PctChange5D: 3.5},
				"MSFT": {PctChange1D: 2.5},
			},
			want: &StocksPctChangesResponse{"AAPL": 1.25},
		},
		{
			name:       "failed to read stock percent changes",
			storeError: errors.New("stock percent changes unavailable"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore, storeHandler := getStockStoreHandler()
			mockStore.On(
				"ReadStockPctChangesMap", mock.Anything,
			).Return(tt.pctChanges, tt.storeError).Once()

			response, err := storeHandler.GetStocksPctChanges(t.Context(), []string{"AAPL", "UNKNOWN"})
			mockStore.AssertExpectations(t)
			if tt.storeError != nil {
				assert.ErrorIs(t, err, tt.storeError)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, response)
		})
	}
}

// test get market analysis
func TestGetMarketAnalysis(t *testing.T) {
	t.Parallel()
	mockStore, storeHandler := getStockStoreHandler()

	tests := []struct {
		name        string
		scope       string
		expectError bool
	}{
		{
			name:        "success stock day pct",
			scope:       StockPct,
			expectError: false,
		},
		{
			name:        "success industry pct",
			scope:       IndustryPct,
			expectError: false,
		},
		{
			name:        "failed stock day pct",
			scope:       StockPct,
			expectError: true,
		},
		{
			name:        "failed industry pct",
			scope:       IndustryPct,
			expectError: true,
		},
		{
			name:        "invalid scope",
			scope:       "invalid",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			if !tt.expectError {
				if tt.scope == StockPct {
					mockStore.On(
						"ReadStockDayPctChanges", mock.Anything,
					).Return(mockStockPctChanges, nil).Once()
					mockStore.On("GetStockList").Return(mockStockInfos).Once()
				}
				if tt.scope == IndustryPct {
					mockStore.On(
						"ReadIndustryPctChanges", mock.Anything,
					).Return(mockIndustryPctChanges, nil).Once()
				}
			} else {
				if tt.scope == StockPct {
					mockStore.On(
						"ReadStockDayPctChanges", mock.Anything,
					).Return(nil, errors.New("stock")).Once()
				}
				if tt.scope == IndustryPct {
					mockStore.On(
						"ReadIndustryPctChanges", mock.Anything,
					).Return(nil, errors.New("industry")).Once()
				}
			}
			response, err := storeHandler.GetMarketAnalysis(t.Context(), tt.scope)
			mockStore.AssertExpectations(t)
			if tt.expectError {
				assert.Error(t, err)
				return
			} else if tt.scope == StockPct {
				output := response.([]StockPctChangeResponse)
				firstData := output[0]
				firstReference := mockStockPctChanges[0]
				assert.Equal(t, len(output), len(mockStockPctChanges))
				assert.Equal(t, firstReference.Symbol, firstData.Symbol)
				assert.Equal(t, firstReference.Date, firstData.Date)
				assert.Equal(t, firstReference.PctChange1D, firstData.PctChange1D)
				assert.Equal(t, firstReference.PctChange20D, firstData.PctChange20D)
			} else if tt.scope == IndustryPct {
				output := response.([]IndustryPctChangeResponse)
				firstData := output[0]
				firstReference := mockIndustryPctChanges[0]
				assert.Equal(t, len(output), len(mockStockPctChanges))
				assert.Equal(t, firstReference.SubIndustry, firstData.SubIndustry)
				assert.Equal(t, firstReference.PctChange1D, firstData.PctChange1D)
				assert.Equal(t, firstReference.PctChange20D, firstData.PctChange20D)
			}
		})
	}
}

func TestGetStockList(t *testing.T) {
	t.Parallel()
	mockStore, storeHandler := getStockStoreHandler()
	mockStore.On("GetStockList").Return(mockStockInfos).Once()
	outputList := storeHandler.GetStockList()
	mockStore.AssertExpectations(t)
	assert.Equal(t, len(outputList), len(mockStockInfos))
	for i, stockName := range sortedStockNames {
		assert.Equal(t, stockName, outputList[i][0])
		assert.Equal(t, mockStockInfos[stockName].CompanyName, outputList[i][1])
	}
}
