package provider

import (
	"context"

	"github.com/ezhigval/weather-proxy/internal/model"
)

type Provider interface {
	Fetch(ctx context.Context, city string) (*model.Weather, error)
}
