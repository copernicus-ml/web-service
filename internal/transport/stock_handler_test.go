package transport

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

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
				mockStore.On("ValidateTicker", input).Return(true).Once()
				err := storeHandler.ValidateTicker(input)
				assert.Nil(t, err)
			} else {
				mockStore.On("ValidateTicker", input).Return(false).Once()
				err := storeHandler.ValidateTicker(strings.ToUpper(input))
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
				assert.Equal(t, firstReference.TSCode, firstData.TSCode)
				assert.Equal(t, firstReference.TradeDate, firstData.TradeDate)
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
