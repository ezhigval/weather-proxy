package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ezhigval/go-toolkit/httputil"
	"github.com/ezhigval/weather-proxy/internal/model"
	"github.com/ezhigval/weather-proxy/internal/service"
)

type WeatherHandler struct {
	svc *service.WeatherService
	log *slog.Logger
}

func New(svc *service.WeatherService, log *slog.Logger) *WeatherHandler {
	return &WeatherHandler{svc: svc, log: log}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if strings.TrimSpace(city) == "" {
		httputil.WriteError(w, httputil.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "city query param is required", nil))
		return
	}

	weather, err := h.svc.GetWeather(r.Context(), city)
	if err != nil {
		httputil.WriteError(w, httputil.NewAppError(http.StatusBadGateway, "UPSTREAM_ERROR", "failed to fetch weather", err))
		return
	}

	if weather.Cached {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	httputil.WriteJSON(w, http.StatusOK, weather)
}

func (h *WeatherHandler) GetBatch(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("cities")
	cities := service.ParseCities(raw)
	if len(cities) == 0 {
		httputil.WriteError(w, httputil.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "cities query param is required", nil))
		return
	}

	if len(cities) > 20 {
		httputil.WriteError(w, httputil.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "max 20 cities per batch", nil))
		return
	}

	results := h.svc.GetBatch(r.Context(), cities)
	httputil.WriteJSON(w, http.StatusOK, model.BatchResponse{Results: results})
}
