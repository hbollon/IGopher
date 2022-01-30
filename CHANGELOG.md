# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0](https://www.github.com/hbollon/IGopher/compare/v0.3.1...v0.4.0) (2022-01-30)


### Features

* add dependencies hash checking during dependency check at startup ([664589e](https://www.github.com/hbollon/IGopher/commit/664589efc7aa7c72affaf1bde056eff916669712))
* add fixed versions for dependencies ([1a1e20f](https://www.github.com/hbollon/IGopher/commit/1a1e20fefd8798b21c8fb60350b9ba52d68a82de))
* **dependency-manager:** set fixed versions for all files and add hashes ([6eea031](https://www.github.com/hbollon/IGopher/commit/6eea031c4b97c5d7d7c1c90a9b0a48bd72fc8e1f))
* manifest file generation and checking ([5d0902f](https://www.github.com/hbollon/IGopher/commit/5d0902f359609df7a11573d197d85a471b41a7e6))


### Bug Fixes

* browser stuck after connexion on /#reactivated url ([b47e1b8](https://www.github.com/hbollon/IGopher/commit/b47e1b8d029cb6b419dac8816a4e497a2e259090))
* dependency search in manifest ([1a1e20f](https://www.github.com/hbollon/IGopher/commit/1a1e20fefd8798b21c8fb60350b9ba52d68a82de))
* **gui:** now handle 'bot crash' msg from Go ([592766a](https://www.github.com/hbollon/IGopher/commit/592766ae725157178f49c30da1249d6c5f7aaf4a))
* handle browser closing event and avoid crash ([70da812](https://www.github.com/hbollon/IGopher/commit/70da812a817f7333ed262173cc82a5ffaafe7926))
* **logs:** removed config dump from logs ([2c437dd](https://www.github.com/hbollon/IGopher/commit/2c437dd340971e834c6773de234df58c66835724))
* update xpath for scrapper and user search input ([408962a](https://www.github.com/hbollon/IGopher/commit/408962a4b5f02a235259c6bd10b3ccef152b10b8))
* update xpath for users scrap and search ([2be277b](https://www.github.com/hbollon/IGopher/commit/2be277bbd20628abc1a6821794c5ed1180b04676))
* **vue:** remove useless console log ([c776b8d](https://www.github.com/hbollon/IGopher/commit/c776b8d0871e0ddd3a358e625bb7ee24db55fe4e))

### [0.3.1](https://www.github.com/hbollon/IGopher/compare/v0.3.0...v0.3.1) (2021-10-26)


### Bug Fixes

* crash during followers list focus ([d020662](https://www.github.com/hbollon/IGopher/commit/d0206621e7548fb86e3950745461bed9797180fc))
* **engine:** add --incognito flag to chrome & update user-agent ([a20c0eb](https://www.github.com/hbollon/IGopher/commit/a20c0eba0d4b9d81f363cd647f03b5b8e012145c))
* **frontend:** launch/hotreload callbacks swiped ([596bc9a](https://www.github.com/hbollon/IGopher/commit/596bc9af65cd15564025cc2be61e0cdb71a17867))
* **frontend:** modal-backdrop not destroying on modal close ([4c5e1c8](https://www.github.com/hbollon/IGopher/commit/4c5e1c840b3f75f4474f45c6d45ceff01ffcc50d))
* now check for the second cookies prompt after login (not always present) ([d020662](https://www.github.com/hbollon/IGopher/commit/d0206621e7548fb86e3950745461bed9797180fc))
* **perf:** slow down goroutines loop to significantly reduce cpu usage ([f0d8739](https://www.github.com/hbollon/IGopher/commit/f0d87398cfe76fbc153be94fa513782832d857f0))

## [0.3.0] - 2021-03-06
### Features
* native proxy support configurable from both GUI/TUI
* new `--background-task` flag for the TUI version to execute IGopher as background task! TUI is also capable to detect running tasks so to stop a background one relaunch TUI without the flag :)
* rework all GUI's frontend to Vue 3 and Bootstrap 5
* **bundler:** update to be compatible with new Vue binaries ([924e15f](https://www.github.com/hbollon/IGopher/commit/924e15f28db14a364bc8e13713836418696ca12e))
* **frontend:** now update radio state with current igopher config ([359431f](https://www.github.com/hbollon/IGopher/commit/359431fe532f91e8c91206f463a9a58cb4aaa3db))
* **frontend:** update scrapper's src users tag input with current config ([ec6aad3](https://www.github.com/hbollon/IGopher/commit/ec6aad31a4453c5ef24a3c60afd3f434076a4f5b))
* **gui:** add download tracking interface on bot launch ([61c4531](https://www.github.com/hbollon/IGopher/commit/61c45312a3d9a579ff6d091143ea39f448c6215e))
* **gui:** replace src users text field by tags input ([3c9d224](https://www.github.com/hbollon/IGopher/commit/3c9d2244f0ce7819379b0ec62c3d19835d1c8915))
* **scripts:** update bundle.sh to do npm operations ([ecfa656](https://www.github.com/hbollon/IGopher/commit/ecfa656c050e10779cbdfc7bb31d442fbfc9568d))
* **vuejs:** add logs view ([fae88b1](https://www.github.com/hbollon/IGopher/commit/fae88b1afec71e0fadf41600235ac0d6c341ea30))
* **vuejs:** add mixin to handle title property on views ([2afcf8e](https://www.github.com/hbollon/IGopher/commit/2afcf8ec88f2a7da39414f7b32317167142c13e0))
* **vuejs:** add settings view and controller with router config ([ac3550b](https://www.github.com/hbollon/IGopher/commit/ac3550b88e4d31acc90e6d9cd55f7983bf56a1cc))
* **vuejs:** convert DmAutomation component script to typescript ([50bd2ba](https://www.github.com/hbollon/IGopher/commit/50bd2ba3eb4f5cceff5e1240d052e01fb61b3c23))bad struct field access on msg listening

* **vuejs:** remove old compatibility scripts for JQuery and obsoletes files ([33e6eb1](https://www.github.com/hbollon/IGopher/commit/33e6eb18b51850a823579f01fc8a03e54e2877f4))
* **vuejs:** replace old izitoast by sweetalert2 ([8323661](https://www.github.com/hbollon/IGopher/commit/8323661602ea3b71956a3c8564c65d2d52584d2a))
* **vuejs:** view for 404 errors and router configuration ([6f42449](https://www.github.com/hbollon/IGopher/commit/6f42449fa80277e7a760b45ffb513bd488d4a375))

### Bug Fixes

* **engine:** cleanup routine execution on electron window closing/crashing
* **astor:** bad struct field access on msg listening ([61c4531](https://www.github.com/hbollon/IGopher/commit/61c45312a3d9a579ff6d091143ea39f448c6215e))
* **chrome:** DevToolsActivePort file doesn't exist error on linux ([#8](https://www.github.com/hbollon/IGopher/issues/8)) ([#10](https://www.github.com/hbollon/IGopher/issues/10)) ([32e66e9](https://www.github.com/hbollon/IGopher/commit/32e66e954277730f29560d327e1e22b5d7bbc9a8))
* **vuejs:** hook execution on route change caused by astilectron listener ([822554e](https://www.github.com/hbollon/IGopher/commit/822554e29686cb3e745372ea0ebf80c20de79393))

**The CHANGELOG and releases will now be automated by __release-please__ workflow.**

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

[0.3.0]: https://github.com/hbollon/igopher/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/hbollon/igopher/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/hbollon/igopher/compare/v0.1.3...v0.2.0
[0.1.3]: https://github.com/hbollon/igopher/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/hbollon/igopher/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/hbollon/igopher/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/hbollon/igopher/releases/tag/v0.1.0
