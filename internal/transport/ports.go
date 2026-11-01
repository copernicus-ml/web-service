package transport

import (
	"context"

	"github.com/copernicus-ml/web-service/internal/repositories"
)

type StockStoreRepo interface {
	GetStockInfo(ticker string) (*repositories.StockInfo, error)
	GetStockList() map[string]repositories.StockInfo
	InvalidateCache()
	ReadFeatures(ctx context.Context, ticker string) ([]repositories.StockFeature, error)
	ReadIndustryPctChanges(ctx context.Context) ([]repositories.IndustryPctChange, error)
	ReadStock(ctx context.Context, ticker string) ([]repositories.StockPrice, error)
	ReadStockDayPctChanges(ctx context.Context) ([]repositories.StockPctChange, error)
	ValidateTicker(ticker string) bool
}
