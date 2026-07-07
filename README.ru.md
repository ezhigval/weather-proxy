# weather-proxy

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
[![CI](https://github.com/ezhigval/weather-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/ezhigval/weather-proxy/actions/workflows/ci.yml)
![License](https://img.shields.io/badge/license-MIT-blue)
![Tier](https://img.shields.io/badge/tier-junior-1d76db)

**[English](README.md)** · Русский

Кэширует ответы погодного API. Batch-эндпоинт работает через worker pool.

## Запуск

```bash
make docker-up
curl "localhost:8081/weather?city=Moscow"
curl "localhost:8081/weather/batch?cities=Moscow,Berlin,Tokyo"
```

Без `OPENWEATHER_API_KEY` отдаёт mock-данные — для разработки нормально.

## Почему порт 8081?

Чтобы не пересечься с url-shortener на 8080 при локальном запуске обоих сервисов.

ADR: [docs/adr/001-mock-provider.ru.md](docs/adr/001-mock-provider.ru.md)
