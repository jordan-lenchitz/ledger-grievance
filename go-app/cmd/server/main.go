package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jordan-lenchitz/ledger-grievance/go-app/docs"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/config"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/handler"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/middleware"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/repository"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/service"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/time/rate"
)

// @title Ledger-Grievance API
// @version 1.0
// @description API for the incident grievance ledger
// @host localhost:8000
// @BasePath /
func main() {
	cfg := config.Load()
	ctx := context.Background()

	var handlerOpts slog.HandlerOptions
	if cfg.LogLevel == "debug" {
		handlerOpts.Level = slog.LevelDebug
	} else {
		handlerOpts.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOpts))
	slog.SetDefault(logger)

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := connectDB(dsn, logger)
	if err != nil {
		logger.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	repo := repository.NewMySQLIncidentRepository(db)
	pkgsiteSvc := service.NewPkgsiteService("")
	svc := service.NewIncidentService(repo, pkgsiteSvc)
	h := handler.NewIncidentHandler(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.SLAGuarantee(cfg.SLAGuaranteeMode))
	r.Use(middleware.QuineNarrator())
	r.Use(otelgin.Middleware("ledger-grievance"))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.ErrorHandler())
	if cfg.StrictCompliance {
		r.Use(middleware.StrictAccept())
		r.Use(middleware.ContentGatekeeper())
		r.Use(middleware.SubSecondPrecondition(func(ctx context.Context, id uint64) (time.Time, error) {
			inc, err := svc.GetIncident(ctx, id)
			if err != nil {
				return time.Time{}, err
			}
			return inc.RecordedAt, nil
		}))
	}
	r.Use(middleware.CompassionateRateLimiter(rate.Limit(5), 10))

	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
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
