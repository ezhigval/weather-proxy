package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezhigval/go-toolkit/httputil"
	"github.com/ezhigval/go-toolkit/logger"
	tkmw "github.com/ezhigval/go-toolkit/middleware"
	tkredis "github.com/ezhigval/go-toolkit/redis"
	"github.com/ezhigval/weather-proxy/internal/cache"
	"github.com/ezhigval/weather-proxy/internal/config"
	"github.com/ezhigval/weather-proxy/internal/handler"
	"github.com/ezhigval/weather-proxy/internal/provider"
	"github.com/ezhigval/weather-proxy/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(logger.Config{Level: cfg.LogLevel, Format: cfg.LogFormat})

	ctx := context.Background()

	rdb := tkredis.NewClient(tkredis.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer func() { _ = tkredis.Close(rdb) }()

	if err := tkredis.Ping(ctx, rdb); err != nil {
		log.Error("redis connection failed", "error", err)
		os.Exit(1)
	}

	var p provider.Provider
	if cfg.OpenWeatherAPIKey != "" {
		p = provider.NewOpenWeather(&http.Client{Timeout: cfg.HTTPTimeout}, cfg.OpenWeatherAPIKey, cfg.OpenWeatherURL)
		log.Info("using openweather provider")
	} else {
		p = provider.NewMock()
		log.Warn("OPENWEATHER_API_KEY not set, using mock provider")
	}

	weatherCache := cache.New(rdb, cfg.CacheTTL)
	svc := service.New(p, weatherCache, cfg.BatchWorkers, cfg.BatchQueueSize)
	defer svc.Close()

	h := handler.New(svc, log)

	r := chi.NewRouter()
	r.Use(tkmw.RequestID)
	r.Use(tkmw.RealIP)
	r.Use(tkmw.Recoverer(log))
	r.Use(tkmw.AccessLog(log))
	r.Use(chimw.Timeout(30 * time.Second))

	r.Get("/health", httputil.HealthHandler(map[string]func() error{
		"redis": func() error { return tkredis.Ping(ctx, rdb) },
	}))
	r.Get("/weather", h.GetWeather)
	r.Get("/weather/batch", h.GetBatch)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
