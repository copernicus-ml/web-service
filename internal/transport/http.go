package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(router *gin.Engine, stockHandler *StockHandler) {
	storeGroup := router.Group("/data", SessionClaimMiddleware())

	// /data/stock endpoints
	{
		stockGroup := storeGroup.Group("/stock")
		// get stock list
		stockGroup.GET("/list", func(c *gin.Context) {
			stockList := stockHandler.GetStockList()
			OK(c, stockList)
		})
		// get stock infomation
		stockGroup.GET("/info/:ticker", func(c *gin.Context) {
			ticker := c.Param("ticker")
			stockInfo, err := stockHandler.GetStockInfo(ticker)
			if err != nil {
				Fail(c, http.StatusNotFound, err.Error())
				return
			}
			OK(c, stockInfo)
		})
		// get stock prices data
		stockGroup.GET("/:ticker", func(c *gin.Context) {
			ticker := c.Param("ticker")
			if err := stockHandler.ValidateTicker(ticker); err != nil {
				Fail(c, http.StatusNotFound, err.Error())
				return
			}
			window := c.Query("window")
			if err := stockHandler.ValidateWindow(window); err != nil {
				Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			stockPrices, err := stockHandler.GetStockPrices(c.Request.Context(), ticker, window)
			if err != nil {
				Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			OK(c, stockPrices)
		})
	}

	// /data/features endpoints
	{
		featuresGroup := storeGroup.Group("/features")
		// get stock features data
		featuresGroup.GET("/:ticker", func(c *gin.Context) {
			ticker := c.Param("ticker")
			if err := stockHandler.ValidateTicker(ticker); err != nil {
				Fail(c, http.StatusNotFound, err.Error())
				return
			}
			featureGroup := c.Query("group")
			if err := stockHandler.ValidateFeatureGroup(featureGroup); err != nil {
				Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			window := c.Query("window")
			if err := stockHandler.ValidateWindow(window); err != nil {
				Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			stockFeatures, err := stockHandler.GetStockFeatures(
				c.Request.Context(), ticker, featureGroup, window)
			if err != nil {
				Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			OK(c, stockFeatures)
		})
	}

	// /data/market endpoints
	{
		marketGroup := storeGroup.Group("/market")
		// get market pct change analysis
		marketGroup.GET("/analysis", func(c *gin.Context) {
			scope := c.Query("scope")
			if err := stockHandler.ValidateMarket(scope); err != nil {
				Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			analysis, err := stockHandler.GetMarketAnalysis(c.Request.Context(), scope)
			if err != nil {
				Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			OK(c, analysis)
		})
	}
}

func SetupInternalRoutes(router *gin.Engine, stockHandler *StockHandler) {
	internalGroup := router.Group("/internal", InternalServiceMiddleware())
	internalGroup.GET("/clear-lru", func(c *gin.Context) {
		response := stockHandler.InvalidateCache()
		OK(c, response)
	})
}
