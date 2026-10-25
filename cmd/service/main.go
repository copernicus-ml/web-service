package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/copernicus-ml/web-service/internal/repositories"
	"github.com/copernicus-ml/web-service/internal/transport"
)

const (
	portNumber        = 8800
	goEnv             = "GO_ENV"
	clerkKeyEnv       = "CLERK_SECRET_KEY"
	allowedOriginsEnv = "ALLOWED_ORIGINS"
)

const (
	runtimeDev  = "dev"
	runtimeProd = "prod"
)

func main() {
	// init logger and context
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()
	// get runtime env
	runtimeEnv := os.Getenv(goEnv)
	if runtimeEnv != runtimeDev && runtimeEnv != runtimeProd {
		logger.Warn("invalid runtime env", "env", runtimeEnv)
		runtimeEnv = runtimeDev
	}

	// init clerk auth client
	clerkSecretKey := os.Getenv(clerkKeyEnv)
	if len(clerkSecretKey) == 0 {
		logger.Error("clerk secret key should not be empty")
		os.Exit(1)
	}
	clerk.SetKey(clerkSecretKey)

	// get allowed orgins
	allowedOrigins := os.Getenv(allowedOriginsEnv)
	if len(allowedOrigins) == 0 {
		logger.Error("allowed origins should not be empty")
		os.Exit(1)
	}
	origins := strings.Split(strings.TrimSpace(allowedOrigins), ";")

	// create a Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// init repositories
	stockStore, err := repositories.NewStockStore(ctx, runtimeEnv, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// init handlers
	stockHandler := transport.NewStockStoreHandler(stockStore, logger)
	transport.SetupStoreRoutes(router, stockHandler)

	protectedRouter := clerkhttp.WithHeaderAuthorization(
		clerkhttp.AuthorizedPartyMatches(origins...),
		clerkhttp.CustomClaimsConstructor(func(context.Context) any {
			return &transport.SessionCustomClaims{}
		}),
	)(router)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", portNumber),
		Handler:      protectedRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start application %w", err)
		os.Exit(1)
	}
}
