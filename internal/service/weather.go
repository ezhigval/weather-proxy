package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ezhigval/weather-proxy/internal/cache"
	"github.com/ezhigval/weather-proxy/internal/model"
	"github.com/ezhigval/weather-proxy/internal/pool"
	"github.com/ezhigval/weather-proxy/internal/provider"
)

type WeatherService struct {
	provider provider.Provider
	cache    *cache.WeatherCache
	pool     *pool.WorkerPool[string, *model.Weather]
}

func New(p provider.Provider, c *cache.WeatherCache, workers, queueSize int) *WeatherService {
	s := &WeatherService{provider: p, cache: c}
	s.pool = pool.New(workers, queueSize, s.fetchOne)
	return s
}

func (s *WeatherService) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*model.Weather, error) {
	city = normalizeCity(city)
	if city == "" {
		return nil, fmt.Errorf("city is required")
	}

	cached, err := s.cache.Get(ctx, city)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	weather, err := s.provider.Fetch(ctx, city)
	if err != nil {
		return nil, err
	}

	weather.City = city
	if err := s.cache.Set(ctx, weather); err != nil {
		return weather, fmt.Errorf("cache set: %w", err)
	}

	return weather, nil
}

func (s *WeatherService) GetBatch(ctx context.Context, cities []string) []model.BatchResult {
	results := make([]model.BatchResult, len(cities))
	var wg sync.WaitGroup

	for i, raw := range cities {
		wg.Add(1)
		go func(idx int, city string) {
			defer wg.Done()
			city = normalizeCity(city)
			if city == "" {
				results[idx] = model.BatchResult{City: raw, Error: "empty city"}
				return
			}

			weather, err := s.pool.Submit(ctx, city)
			if err != nil {
				results[idx] = model.BatchResult{City: city, Error: err.Error()}
				return
			}

			results[idx] = model.BatchResult{City: city, Weather: weather}
		}(i, raw)
	}

	wg.Wait()
	return results
}

func (s *WeatherService) fetchOne(ctx context.Context, city string) (*model.Weather, error) {
	cached, err := s.cache.Get(ctx, city)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	weather, err := s.provider.Fetch(ctx, city)
	if err != nil {
		return nil, err
	}

	weather.City = city
	if err := s.cache.Set(ctx, weather); err != nil {
		return weather, fmt.Errorf("cache set: %w", err)
	}

	return weather, nil
}

func normalizeCity(city string) string {
	return strings.TrimSpace(city)
}

func ParseCities(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if c := normalizeCity(p); c != "" {
			out = append(out, c)
		}
	}
	return out
}
