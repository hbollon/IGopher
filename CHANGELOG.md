# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2021-03-06
### Added
- [GUI] Information notification on bot stop/hot-reload
### Changed
- Improve selenium closing routine
- [Gui] Notification triggering on bot crash and bot running state reset
- Set icons, single instance and version options to astilectron in gui development package
- Replace custom scripts min by full ones for easier contributions
- Allow bot exit before ig connection and scrapping process

### Fixed
- Duplicate CloseSelenium call on bot stop
- Scrapper issue if src user doesn't exist
- Scrapper issue if src user is private
- Scrapper issue if src user hasn't enough followers than requested
- Abort blocking mpb progress bar on user fetching error
- Clean go.mod
## [0.2.0] - 2021-03-04
### Added
- Electron GUI with: 
  - DM Automation config screen with launch/stop/hot-reload actions
  - Global settings view
  - Logs explorer
- Logrus dual output on stdout with curom formatter and log file with json formatter
- Bundler github workflow
### Changed
- Parallelization of bot execution on several goroutines (once for engine and once for communication with main goroutine) with context/channels
- IGopher architecture refactor

### Fixed
- Fix project environment location issue #3
- Linters related issues
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

[Unreleased]: https://github.com/hbollon/igopher/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/hbollon/igopher/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/hbollon/igopher/compare/v0.1.3...v0.2.0
[0.1.3]: https://github.com/hbollon/igopher/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/hbollon/igopher/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/hbollon/igopher/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/hbollon/igopher/releases/tag/v0.1.0
