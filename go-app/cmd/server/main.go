package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/config"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/handler"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/middleware"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/repository"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/service"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/time/rate"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/jordan-lenchitz/ledger-grievance/go-app/docs"
)

// @title Ledger-Grievance API
// @version 1.0
// @description API for the incident grievance ledger
// @host localhost:8000
// @BasePath /
func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Setup Logger
	var handlerOpts slog.HandlerOptions
	if cfg.LogLevel == "debug" {
		handlerOpts.Level = slog.LevelDebug
	} else {
		handlerOpts.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOpts))
	slog.SetDefault(logger)

	// Initialize OTEL
	shutdown, err := telemetry.SetupOTEL(ctx)
	if err != nil {
		logger.Error("failed to setup otel", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Error("failed to shutdown otel", "error", err)
		}
	}()

	// Database Connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := connectDB(dsn, logger)
	if err != nil {
		logger.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize Layers
	repo := repository.NewMySQLIncidentRepository(db)
	pkgsiteSvc := service.NewPkgsiteService("")
	svc := service.NewIncidentService(repo, pkgsiteSvc)
	h := handler.NewIncidentHandler(svc)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("ledger-grievance"))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CompassionateRateLimiter(rate.Limit(5), 10))

	// Routes
	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h.RegisterRoutes(r)

	// Server Configuration
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exiting")
}

func connectDB(dsn string, logger *slog.Logger) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		logger.Warn("failed to connect to database, retrying...", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
