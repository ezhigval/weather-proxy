# weather-proxy

Caches weather API responses. Batch endpoint uses a worker pool.

## Run

```bash
make docker-up
curl "localhost:8081/weather?city=Moscow"
curl "localhost:8081/weather/batch?cities=Moscow,Berlin,Tokyo"
```

Without `OPENWEATHER_API_KEY` you get mock data — fine for dev.

## Why port 8081?

Avoids collision with url-shortener on 8080 when running both locally.

ADR: [docs/adr/001-mock-provider.md](docs/adr/001-mock-provider.md)
