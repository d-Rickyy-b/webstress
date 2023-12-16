# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2023-12-16

### Added
- Make websockets recover from certain errors
- Implement ratelimiting for websockets

### Fixed
- Fixed display issue with Msg/s/socket


## [1.1.0] - 2023-04-02

### Security
- Fix vulnerability in [golang.org/x/text/language](https://github.com/d-Rickyy-b/webstress/security/dependabot/1)
- Fix Out-of-bounds Read vulnerability in [golang.org/x/text/language](https://github.com/d-Rickyy-b/webstress/security/dependabot/2)
- Fix vunlearbility in [golang.org/x/sys/unix](https://github.com/d-Rickyy-b/webstress/security/dependabot/3)


## [1.0.0] - 2022-08-11

First release


[unreleased]: https://github.com/d-Rickyy-b/webstress/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/d-Rickyy-b/webstress/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/d-Rickyy-b/webstress/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/d-Rickyy-b/webstress/releases/tag/v1.0.0
