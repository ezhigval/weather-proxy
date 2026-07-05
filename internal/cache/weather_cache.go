package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ezhigval/weather-proxy/internal/model"
	"github.com/redis/go-redis/v9"
)

type WeatherCache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *WeatherCache {
	return &WeatherCache{client: client, ttl: ttl}
}

func (c *WeatherCache) key(city string) string {
	return fmt.Sprintf("weather:%s", city)
}

func (c *WeatherCache) Get(ctx context.Context, city string) (*model.Weather, error) {
	data, err := c.client.Get(ctx, c.key(city)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var w model.Weather
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("unmarshal weather: %w", err)
	}

	w.Cached = true
	return &w, nil
}

func (c *WeatherCache) Set(ctx context.Context, weather *model.Weather) error {
	copy := *weather
	copy.Cached = false

	data, err := json.Marshal(copy)
	if err != nil {
		return fmt.Errorf("marshal weather: %w", err)
	}

	if err := c.client.Set(ctx, c.key(copy.City), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
