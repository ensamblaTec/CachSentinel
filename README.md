# CachSentinel 🛡️

### High-Performance Predictive Cache Proxy

CachSentinel es un proxy inverso de caché de grado industrial desarrollado en Go. Está diseñado para sistemas donde la latencia P99 debe ser plana y la disponibilidad es innegociable (MES, WMS, Fintech).

## Características Principales

- **Arquitectura Hexagonal:** Desacoplamiento total entre la lógica de negocio y la infraestructura (Memoria/Redis/HTTP).
- **Predictive Background Refresh:** El sistema aprende del tráfico. Si una llave es "caliente" (popular) y está próxima a expirar, se refresca automáticamente en segundo plano antes de que el usuario lo solicite.
- **Zero-Copy Optimization:** Almacenamiento directo de `[]byte`. Elimina el overhead de reflexión y serialización JSON en el path crítico.
- **Singleflight Consolidation:** Evita el *Cache Stampede* agrupando peticiones concurrentes idénticas en una sola llamada al upstream.
- **Stale-While-Revalidate (SWR):** Si el servidor de origen falla, CachSentinel sirve datos expirados durante un periodo de gracia, garantizando que el sistema nunca se detenga.

## Arquitectura de Capas

```text
internal/
├── core/
│   ├── domain/     # Entidades puras (CacheEntry, Config)
│   ├── ports/      # Interfaces (Contratos del repositorio y fetcher)
│   └── service/    # El "Cerebro": Lógica de predicción, SWR y concurrencia
└── infrastructure/
    ├── adapter/    # Implementaciones (MemoryStore, HTTPFetcher)
    └── api/        # Capa de transporte (Proxy HTTP)
```

