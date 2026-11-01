package transport

import (
	"context"

	"github.com/copernicus-ml/web-service/internal/repositories"
	"github.com/stretchr/testify/mock"
)

// Mock StockStore Repository
type mockStockStoreRepo struct {
	mock.Mock
}

func (m *mockStockStoreRepo) GetStockInfo(ticker string) (*repositories.StockInfo, error) {
	args := m.Called(ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.StockInfo), args.Error(1)
}

func (m *mockStockStoreRepo) GetStockList() map[string]repositories.StockInfo {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]repositories.StockInfo)
}

func (m *mockStockStoreRepo) InvalidateCache() {
	m.Called()
}

func (m *mockStockStoreRepo) ReadFeatures(
	ctx context.Context, ticker string) ([]repositories.StockFeature, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.StockFeature), args.Error(1)
}

func (m *mockStockStoreRepo) ReadIndustryPctChanges(
	ctx context.Context) ([]repositories.IndustryPctChange, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.IndustryPctChange), args.Error(1)
}

func (m *mockStockStoreRepo) ReadStock(
	ctx context.Context, ticker string) ([]repositories.StockPrice, error) {
	args := m.Called(ctx, ticker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.StockPrice), args.Error(1)
}

func (m *mockStockStoreRepo) ReadStockDayPctChanges(
	ctx context.Context) ([]repositories.StockPctChange, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.StockPctChange), args.Error(1)
}

func (m *mockStockStoreRepo) ValidateTicker(ticker string) bool {
	args := m.Called(ticker)
	return args.Bool(0)
}

var mockStockPctChanges = []repositories.StockPctChange{
	{
		TSCode:       "AAPL",
		TradeDate:    "20260102",
		SubIndustry:  "Technology Hardware, Storage & Peripherals",
		PctChange1D:  1.25,
		PctChange5D:  3.50,
		PctChange10D: 5.75,
		PctChange20D: 8.25,
		PctChange60D: 12.50,
	},
	{
		TSCode:       "MSFT",
		TradeDate:    "20260102",
		SubIndustry:  "Systems Software",
		PctChange1D:  -0.75,
		PctChange5D:  1.80,
		PctChange10D: 4.20,
		PctChange20D: 6.10,
		PctChange60D: 10.30,
	},
}

var mockIndustryPctChanges = []repositories.IndustryPctChange{
	{
		SubIndustry:  "Technology Hardware, Storage & Peripherals",
		PctChange1D:  0.90,
		PctChange5D:  2.65,
		PctChange10D: 4.95,
		PctChange20D: 7.35,
		PctChange60D: 11.40,
	},
	{
		SubIndustry:  "Systems Software",
		PctChange1D:  -0.30,
		PctChange5D:  1.15,
		PctChange10D: 3.85,
		PctChange20D: 5.80,
		PctChange60D: 9.90,
	},
}

var mockStockInfos = map[string]repositories.StockInfo{
	"MSFT": {
		TSCode:       "MSFT",
		CompanyName:  "Microsoft Corporation",
		Sector:       "Information Technology",
		SubIndustry:  "Systems Software",
		Headquarters: "Redmond, Washington",
		DateAdded:    "19940601",
		CIK:          "0000789019",
		Founded:      "1975",
	},
	"AAPL": {
		TSCode:       "AAPL",
		CompanyName:  "Apple Inc.",
		Sector:       "Information Technology",
		SubIndustry:  "Technology Hardware, Storage & Peripherals",
		Headquarters: "Cupertino, California",
		DateAdded:    "19821231",
		CIK:          "0000320193",
		Founded:      "1976",
	},
	"JPM": {
		TSCode:       "JPM",
		CompanyName:  "JPMorgan Chase & Co.",
		Sector:       "Financials",
		SubIndustry:  "Diversified Banks",
		Headquarters: "New York, New York",
		DateAdded:    "19910506",
		CIK:          "0000019617",
		Founded:      "1799",
	},
	"JNJ": {
		TSCode:       "JNJ",
		CompanyName:  "Johnson & Johnson",
		Sector:       "Health Care",
		SubIndustry:  "Pharmaceuticals",
		Headquarters: "New Brunswick, New Jersey",
		DateAdded:    "19730331",
		CIK:          "0000200406",
		Founded:      "1886",
	},
	"XOM": {
		TSCode:       "XOM",
		CompanyName:  "Exxon Mobil Corporation",
		Sector:       "Energy",
		SubIndustry:  "Integrated Oil & Gas",
		Headquarters: "Spring, Texas",
		DateAdded:    "19570304",
		CIK:          "0000034088",
		Founded:      "1999",
	},
	"AMZN": {
		TSCode:       "AMZN",
		CompanyName:  "Amazon.com, Inc.",
		Sector:       "Consumer Discretionary",
		SubIndustry:  "Broadline Retail",
		Headquarters: "Seattle, Washington",
		DateAdded:    "20051118",
		CIK:          "0001018724",
		Founded:      "1994",
	},
}
var sortedStockNames = []string{"AAPL", "AMZN", "JNJ", "JPM", "MSFT", "XOM"}
