package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ezhigval/weather-proxy/internal/model"
)

type OpenWeather struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

func NewOpenWeather(client *http.Client, apiKey, baseURL string) *OpenWeather {
	return &OpenWeather{client: client, apiKey: apiKey, baseURL: baseURL}
}

type owResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func (p *OpenWeather) Fetch(ctx context.Context, city string) (*model.Weather, error) {
	u, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("q", city)
	q.Set("appid", p.apiKey)
	q.Set("units", "metric")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openweather request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openweather status %d: %s", resp.StatusCode, string(body))
	}

	var data owResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	desc := "unknown"
	if len(data.Weather) > 0 {
		desc = data.Weather[0].Description
	}

	return &model.Weather{
		City:        data.Name,
		Temperature: data.Main.Temp,
		Description: desc,
		Humidity:    data.Main.Humidity,
		WindSpeed:   data.Wind.Speed,
	}, nil
}

type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Fetch(_ context.Context, city string) (*model.Weather, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, fmt.Errorf("city is empty")
	}

	seed := int([]rune(strings.ToLower(city))[0])
	return &model.Weather{
		City:        city,
		Temperature: float64(10 + seed%20),
		Description: "partly cloudy (mock)",
		Humidity:    40 + seed%40,
		WindSpeed:   float64(seed%10) + 1.5,
	}, nil
}
