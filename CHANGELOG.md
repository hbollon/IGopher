# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2021-02-21
### Added
- Useful repository files including:
  - CONTRIBUTING.md
  - Issues & PR templates
  - Changelog file

### Changed
- Add new config & linters to golangci workflow

### Fixed
- Chrome/ChromeDriver dependencies incompatibility issues with MacOS
- Terminal cleaning issue with MacOS
- Variable shadowing issues
- goconst & lll linters related issues

## [0.1.2] - 2021-02-06
### Changed
- Moved TUI to internal sub-package
- Refactor TUI Update/View logic

### Fixed
- Reduce cyclomatic complexities of some functions
- Golint issues

## [0.1.1] - 2021-01-31
### Changed
- Update README with better installation instructions

### Fixed
- Issue with scrapper config model

## [0.1.0] - 2021-01-31
IGopher come in this first public pre-release with cross-platform (Linux/Windows at the moment) compatibility and a user-friendly terminal user interface!
At this point, the bot will first retrieve a user list from the followers of the source users that you have entered. It will then send a message according to the templates that you put to them.
In addition, you can activate certain modules such as:

- A heuristic and daily quota limiter
- A scheduler
- The use of a blacklist to avoid duplicates interactions

[Unreleased]: https://github.com/hbollon/igopher/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/hbollon/igopher/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/hbollon/igopher/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/hbollon/igopher/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/hbollon/igopher/releases/tag/v0.1.0
