# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] v0.6.0

## 23 Jan 2021
### Added

### Changed

### Fixed

### Removed


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


[Unreleased]: https://github.com/rovergulf/utils/compare/v0.5.0...main
[v0.5.0]: https://github.com/rovergulf/utils/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/rovergulf/utils/compare/v0.3.0...v0.4.0
[v0.2.0]: https://github.com/rovergulf/utils/compare/v0.1.7...v0.2.0
[v0.1.7]: https://github.com/rovergulf/utils/compare/v0.1.5...v0.1.7
[v0.1.5]: https://github.com/rovergulf/utils/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/rovergulf/utils/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/rovergulf/utils/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/rovergulf/utils/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/rovergulf/utils/tree/v0.1.1
