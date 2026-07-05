# weather-proxy

HTTP weather proxy with Redis caching, worker pool for batch requests, and OpenWeatherMap integration.

Part of [Go Backend Portfolio](https://github.com/ezhigval).

## Features

- `GET /weather?city=Moscow` — single city lookup
- `GET /weather/batch?cities=Moscow,London,Tokyo` — parallel batch via worker pool
- Redis TTL cache
- OpenWeatherMap provider (or mock mode without API key)
- Graceful shutdown, structured logging

## Quick Start

```bash
make docker-up

curl -s "http://localhost:8081/weather?city=Moscow" | jq
curl -s "http://localhost:8081/weather/batch?cities=Moscow,London,Berlin" | jq
```

Without API key the service runs in **mock mode** — handy for local demo and CI.

## OpenWeatherMap

```bash
export OPENWEATHER_API_KEY=your_key
export REDIS_ADDR=localhost:6380
make run
```

## API

| Method | Path | Description |
|---|---|---|
| GET | `/weather?city=` | Weather for one city |
| GET | `/weather/batch?cities=` | Batch up to 20 cities |
| GET | `/health` | Health check |

## Config

| Env | Default | Description |
|---|---|---|
| `PORT` | `8081` | HTTP port |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `CACHE_TTL` | `10m` | Cache TTL |
| `OPENWEATHER_API_KEY` | — | API key (mock if empty) |
| `BATCH_WORKERS` | `4` | Worker pool size |

## Architecture

```
HTTP → handler → service → cache (redis)
                         ↘ provider (openweather/mock)
                         ↘ worker pool (batch)
```

## Stack

Go · chi · Redis · OpenWeatherMap · go-toolkit

## License

MIT
