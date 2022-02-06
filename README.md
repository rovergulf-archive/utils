![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rovergulf/utils)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rovergulf/utils)
[![Golang CI](https://github.com/rovergulf/utils/actions/workflows/main.yml/badge.svg)](https://github.com/rovergulf/utils/actions/workflows/main.yml)

# 🚜 utils
Rovergulf Engineers common utils repository

### 🦍 packages
- colors - package to generate random colors in hsv or rgb
- httplib - http utility library
  - Interceptor wraps [gorilla/mux](https://github.com/gorilla/mux) Router
- mq
  - [Sarama/shopify]([jackc/pgx](https://github.com/Sarama/shopify)) Kafka consumer wrapper
  - [nats/nats.go](https://github.com/nats-io/nats.go) and [nats/stan.go](https://github.com/nats-io/stan.go) wrapper
- pgxs - [jackc/pgx](https://github.com/jackc/pgx) and [jackc/tern](https://github.com/jackc/tern) wrapper
- tracing - [Jaeger Tracing client](https://github.com/uber/jaeger-client-go) wrapper
- useragent - discover http.Request User Agent header value
- zapx - [uber-go/zap](https://github.com/uber-go/zap) initializer
