# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] v1.15.0

## 19 May 2022

### Added

### Changed

### Fixed

### Removed


## 18 May 2022

### Added

### Changed
- Updated dependencies

### Fixed

### Removed


## Released [v1.14.4]

## 24 Feb 2022

### Changed
- pgxs - fix debug query logging


## Released [v1.14.2]

## 8 Feb 2022

### Changed
- Update dependencies


## Released [v1.14.1]

## 6 Feb 2022

### Changed
- optimized root package functions
- `httplib` - use Cloudflare connecting ip header as request context value

### Fixed
- `ipaddr` - Cloudflare headers naming
- `timers` - removed unnecessary zap logger usage
- `utils` - sorting array comparison operators

### Removed
- `tlsconf` and `acksequence` packages


## Released [v1.14.0]

## 13 Jan 2022

### Added
- pgxs database stats methods (thanks )

### Changed
- Updated dependencies

### Fixed

### Removed
- GitHub and MongoDB packages


## Released [v1.11.0]

## 10 Nov 2021

### Changed
- go.mod updated up to go 1.17 version
- dependencies versions are up to last

### Fixed
- `httplib` cookies value extraction, using `http` package native method

### Notes
There are nats-server Dependabot security issue, which probably references to go.sum file  
Not sure how to handle this, so version may be increased at its minor (to fix that issue) later,  
if I not hit that by this tag

## Released [v1.10.1]

## 12 Sep 2021

### Changed
- Trace and set request context value of `X-Forwarded-For` header in `httplib`
- Update dependencies

## Released [v1.10.0]

## 7 Sep 2021

### Changed
- `httplib` now uses `http.Request.RemoteAddr` to identify client as context value, instead `ipaddr.GetRequestIPAddress`
  I'm pretty sure `ipaddr` may be deleted in near future, due poor design and lack of usage
- Update dependencies

### Fixed

### Removed

## Released [v1.9.0]

## 30 Aug 2021

### Changed
- Update dependencies


## Released [v1.8.0]

## 29 Apr 2021

### Added
- **response** package - handles json and yaml output to specified `io.Writer` interface


## Released [v1.7.0]

## 16 Apr 2021

### Added
- `mdb` package to support [MongoDB driver](https://github.com/mongodb/mongo-go-driver)


## Released [v1.6.0]

## 2 Apr 2021

### Added
- validation.go shortlinkRegex, fullShortlinkRegex and resourceName regex
- **zapx** package as [uber-go/zap](https://github.com/uber-go/zap) wrapper
- `ipaddr` cloudflare http header ip handling, 
- `httplib.Interceptor` ipaddr usage for tracing


## Released [v1.5.0]

## 1 Apr 2021

### Added
- ./github package `ghx` with Github API & OAuth2 client simple wrapper

## Released [v1.4.0]

## 30 Mar 2021

### Added
- ❕`pgxs` package now have `Migrate` method which represents [jackc/tern](https://github.com/jackc/tern) library
- Github Workflow configuration

### Changed
- ❗BREAKING: pgxs package config now uses `pgx.Config.MigrationsPath` for migrations
- pgxs - separate package logics by files, renamed functions for `pgxpool.Pool` and some preparations to handle `pgx.Conn`
- go.mod now uses Golang 1.16

### Fixed

### Removed
- `pgxs.Config.ActualSchema` parameter


## Released [v1.3.2]

## 19 Mar 2021

### Changed
- `tracing` package connection log level set to debug

## Released [v1.3.1]

## 11 Mar 2021

### Added
- `httplib.Interceptor` struct

### Fixed
- `natsmq` client id declaration


## Released [v1.3.0]

## 23 Feb 2021

### Added
- `useragent` package

### Changed

### Fixed

### Removed


## Released [v1.2.0]

## 18 Feb 2021

### Added
- `natsmq` – Message handlers moved to separate Nats/StreamingSub functions


## Released [v1.1.0]

## 18 Feb 2021

### Added
- `natsmq` – Message handler `context.Context` argument


## Released [v1.0.0]

## 17 Feb 2021

### Added
- `colors` – random color generating package


## Released [v0.11.0]

## 17 Feb 2021

### Fixed
- `natsmq` nuid lock
- `natsmq` lowercase struct internals


## Released [v0.10.8]

## 17 Feb 2021

### Changed
- `natsmq` log beautify


## Released [v0.10.0]

## 17 Feb 2021

### Added
- `natsmq` enhanced zap logging
- `natsmq` client dynamic nuid for ack's

### Changed
- go mod tidy
- update dependencies
- `natsmq.Config` includes `opentracing.Tracer` now to reduce call arguments and duplicate functions

### Removed
- `natsmq` removed unnecessary loges


## Released [v0.9.1]

## 16 Feb 2021

### Added
- `natsmq.MessageInfo` struct for `StanSubHandler` argument


## Released [v0.8.3]

## 16 Feb 2021

### Fixed
- `natsmq` nuid argument value


## Released [v0.8.2]

## 16 Feb 2021

### Fixed
- `natsmq` named loggin

### Removed
- `natsmq` subscription.Delivered() method call


## Released [v0.8.0]

## 16 Feb 2021

### Added
- `natsmq.NatsSubOpts` and `natsmq.StanSubOpts` specified for subscriptions only

### Changed
- `natsmq` methods only require nats config now


## Released [v0.7.0]

## 16 Feb 2021

### Added
- `natsmq.StanSubHandler` func interface supporting sequence and reply arguments
- `natsmq.NatsSUbHandler` func interface supporting reply argument


## Released [v0.6.2]

## 14 Feb 2021

### Fixed
- `natsmq` logging


## Released [v0.6.1]

## 14 Feb 2021

### Fixed
- `pgxs` tls connection setup



## Released [v0.6.0]

## 11 Feb 2021

### Changed
- `pgxs` package Api changed
- mq/nats - removed ack manager usage

### Removed
- `hub` package


## Released [v0.5.1]

## 22 Jan 2021

### Added
- strings unwrap quotes method


## Released [v0.5.0]

## 13 Jan 2021

### Changed
- httplib interceptor removed in prior gorilla/mux middlewares for better experience, I will think about how package is helpful 
- pgxs

### Removed
- rpcd package
- clog package


## Released [v0.2.0]

## 10 Jan 2021

### Added
- `httplib` interceptor JSON Response methods with debug zap logging


## Released [v0.1.7] & [v0.1.6]
## 8 Jan 2021

### Changed
- added zap.SugaredLogger usage to `pgxs` packages

## Released [v0.1.5]
## 8 Jan 2021

### Changed
- added zap.SugaredLogger usage to `mq` packages


## Released [v0.1.4]
## 7 Jan 2021

### Changed
- downgraded `grpc` version from 1.34 to 1.29.1 in case of backward compatibility with ectd client


## Released [v0.1.3]
## 7 Jan 2021

### Changed
- updated `httplib` model from `rovergulf/auth/httplib`


## Released [v0.1.2]
## 7 Jan 2021

### Changed
- added zap logger to strings hash method
- `httplib` interceptors to use zap logger


## Released [v0.1.1]
## 7 Jan 2021

### Added
- moved `pkg` from private repos to public `utils`


[Unreleased]: https://github.com/rovergulf/utils/compare/v1.14.4...main
[v1.14.4]: https://github.com/rovergulf/utils/compare/v1.14.2...v1.14.4
[v1.14.2]: https://github.com/rovergulf/utils/compare/v1.14.1...v1.14.2
[v1.14.1]: https://github.com/rovergulf/utils/compare/v1.14.0...v1.14.1
[v1.14.0]: https://github.com/rovergulf/utils/compare/v1.11.0...v1.14.0
[v1.11.0]: https://github.com/rovergulf/utils/compare/v1.10.1...v1.11.0
[v1.10.0]: https://github.com/rovergulf/utils/compare/v1.10.0...v1.10.1
[v1.10.0]: https://github.com/rovergulf/utils/compare/v1.9.0...v1.10.0
[v1.9.0]: https://github.com/rovergulf/utils/compare/v1.8.0...v1.9.0
[v1.8.0]: https://github.com/rovergulf/utils/compare/v1.7.0...v1.8.0
[v1.7.0]: https://github.com/rovergulf/utils/compare/v1.6.0...v1.7.0
[v1.6.0]: https://github.com/rovergulf/utils/compare/v1.5.0...v1.6.0
[v1.5.0]: https://github.com/rovergulf/utils/compare/v1.4.0...v1.5.0
[v1.4.0]: https://github.com/rovergulf/utils/compare/v1.3.2...v1.4.0
[v1.3.2]: https://github.com/rovergulf/utils/compare/v1.3.1...v1.3.2
[v1.3.1]: https://github.com/rovergulf/utils/compare/v1.3.0...v1.3.1
[v1.3.0]: https://github.com/rovergulf/utils/compare/v1.2.0...v1.3.0
[v1.2.0]: https://github.com/rovergulf/utils/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/rovergulf/utils/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/rovergulf/utils/compare/v1.0.0...v1.0.0
[v0.11.0]: https://github.com/rovergulf/utils/compare/v0.10.8...v0.11.0
[v0.10.8]: https://github.com/rovergulf/utils/compare/v0.10.0...v0.10.8
[v0.10.0]: https://github.com/rovergulf/utils/compare/v0.9.1...v0.10.0
[v0.9.1]: https://github.com/rovergulf/utils/compare/v0.8.3...v0.9.1
[v0.8.3]: https://github.com/rovergulf/utils/compare/v0.8.2...v0.8.3
[v0.8.2]: https://github.com/rovergulf/utils/compare/v0.8.0...v0.8.2
[v0.8.0]: https://github.com/rovergulf/utils/compare/v0.7.0...v0.8.0
[v0.7.0]: https://github.com/rovergulf/utils/compare/v0.6.2...v0.7.0
[v0.6.2]: https://github.com/rovergulf/utils/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/rovergulf/utils/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/rovergulf/utils/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/rovergulf/utils/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/rovergulf/utils/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/rovergulf/utils/compare/v0.3.0...v0.4.0
[v0.2.0]: https://github.com/rovergulf/utils/compare/v0.1.7...v0.2.0
[v0.1.7]: https://github.com/rovergulf/utils/compare/v0.1.5...v0.1.7
[v0.1.5]: https://github.com/rovergulf/utils/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/rovergulf/utils/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/rovergulf/utils/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/rovergulf/utils/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/rovergulf/utils/tree/v0.1.1
