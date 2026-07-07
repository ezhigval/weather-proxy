# ADR-001: Mock provider when API key missing

**English** · [Русский](001-mock-provider.ru.md)

**Status:** Accepted  
**Date:** 2026-07-06

## Context

OpenWeatherMap requires API key. Portfolio demos and CI should run without secrets.

## Decision

If `OPENWEATHER_API_KEY` is empty, log a warning and use `Mock` provider with deterministic output.

## Consequences

- Easy local demo
- Tests never hit external API
- Production deploy must set the key explicitly
