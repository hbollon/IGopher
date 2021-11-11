<h1 align="center">IGopher : (WIP) Golang smart bot for Instagram DM automation</h1>
<p align="center">
    <img alt="IGopher logo" height="250" src="doc/IGopher.png">
</p>
<p align="center">
  <a href="https://github.com/hbollon/IGopher/actions" target="_blank">
    <img alt="Build CI" src="https://github.com/hbollon/igopher/workflows/build/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/hbollon/igopher" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hbollon/igopher" />
  </a>
  <a href="https://github.com/hbollon/igopher/blob/master/LICENSE.md" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/hbollon/go-instadm" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/hbollon/go-instadm" alt="PkgGoDev">
  </a>
</p>

<p align="center">⚡ Powerful, customizable and easy to use Instagram dm bot. With TUI and Eletron.js GUI! Using Selenium webdriver and Yaml configuration files.</p>

<p align="center"><strong>This project is under active development, there may be bugs or missing features. If you have any problem or would like to see a feature implemented, please, open an issue. This is essential so that we can continue to improve IGopher! </strong></p>


---

> Disclaimer: This is a research project. I am in no way responsible for the use you made of this tool. In addition, I am not responsible for any sanctions and/or limitations imposed on your account after using this bot.

---

## Table of Contents

- [Presentation](#presentation)
  - [Graphical User Interface](#graphical-user-interface)
  - [Terminal User Interface](#terminal-user-interface)
- [Features](#features)
- [Getting Started](#getting-started)
  - [From release](#from-release)
  - [From sources](#from-sources)
  - [Flags](#flags)
- [Known Issues](#known-issues)
- [Contributing](#-contributing)
- [Author](#author)
- [License](#-license)

## Presentation

IGopher is a new Instagram automation tool that aims to simplify the deployment of such tools and make their use more pleasant thanks to a TUI (Terminal User Interface) as well as a GUI (Graphical User Interface) powered with Electron.js!

### Graphical User Interface

<p align="center">
  <img src="doc/gifs/demo_gui.gif">
  <small>A beautiful, cross-platform and easy to use interface! Build with Electron.js and <a href="https://github.com/asticode/go-astilectron">go-astilectron</a>.</small>
</p>

Come with **Hot Reload** functionality to apply configuration changes without restart !
Bot stopping and hot reloading are actions safe by waiting bot idle to execute.

### Terminal User Interface

<p align="center">
  <img src="doc/gifs/demo.gif">
  <small>Automatic user fetching and message sending!</small>
</p>

Thanks to the TUI you can easily use this tool on a not very powerful machine, in ssh, on a Vps or even on an operating system without graphical interface!
The bot configuration is very easy thanks to the different configuration menus in the TUI. Parameters are managed and saved in Yaml files easy to edit manually!
All dependencies are downloaded and managed automatically.

<p align="center">
  <img src="doc/gifs/demo_tui.gif">
  <small>Easily configurable and easy to use thanks to his TUI !</small>
</p>

### Requirements
- [Java 8 or 11](https://java.com/fr/download/) (incompatible with newer versions yet)
- For Windows:
  - [Optionnal] [Windows Terminal](https://www.microsoft.com/fr-fr/p/windows-terminal/9n0dx20hk701?activetab=pivot:overviewtab) -> in order to have a best TUI experience

## Features
- Selenium webdriver engine :stars:
- Automatic dependencies downloading and installation :stars:
- Automated IG connection & message sending :stars:
- Users scrapping from ig user followers :stars:
- Scheduler :stars:
- Quotas & user blacklist modules :stars:
- Human writing simulation :stars:
- Fully and easily customizable through Yaml files or with TUI :stars:
- TUI (Terminal User Interface) :stars:
- GUI (Graphical User Interface) powered with Electron.js :stars:
  - Hot Reload functionality to apply configuration changes without restart !
  - Stop and Hot Reload are actions safe by waiting bot idle to execute !
- Many more to come ! 🥳

**Check this [Project](https://github.com/hbollon/igopher/projects/1) to see all planned features for this tool! Feel free to suggest additional features to implement! 🥳**

## Getting Started

### From release

#### GUI version:

1. Download and install [Java 8 or 11](https://java.com/fr/download/) (needed for Selenium webdriver) and add them to your path (on Windows)
2. Download [lastest release](https://github.com/hbollon/igopher/releases/latest) GUI executable for your operating system
3. Move the executable to a dedicated folder (it will create folders/files)
4. Launch it
- For the moment, on MacOS, you must move the .app to your Applications folder and execute the binary file located inside the .app one. It will be improved soon!
5. Configure the bot with your Instagram credentials and your desired scrapping and autodm settings.
6. You're ready! Just hit the "Launch" option on the dm automation page 🚀 
IGopher will download all needed dependencies automatically, don't panic if it seems stuck. I will implement a download monitoring view soon :smile:

#### TUI version:

1. Download and install [Java 8 or 11](https://java.com/fr/download/) (needed for Selenium webdriver) and add them to your path (on Windows)
2. Download [lastest release](https://github.com/hbollon/igopher/releases/latest) TUI executable for your operating system
3. Move the executable to a dedicated folder (it will create folders/files)
4. Launch it:
- On Windows, open a **Windows Terminal** in the folder (or powershell/cmd but the experience quality can be lower) and execute it: ```./tui.exe``` or just drag and drop tui.exe in your command prompt
- On Linux or MacOS, open you favorite shell in the folder, allow it to be executed with ```chmod +x ./tui``` and launch it: ```./tui```
5. Configure the bot with your Instagram credentials and set your desired scrapping and autodm settings. To do that, you can use the TUI settings screen or directly edit the config.yaml file.
6. You're ready! Just hit the "Launch" option in the TUI main menu 🚀

### From sources

#### GUI version:

##### With bundles

1. Download and install [Java 8 or 11](https://java.com/fr/download/) (needed for Selenium webdriver) and add them to your path (on Windows)
2. Install [Go](https://golang.org/doc/install) on your system
3. Download [lastest release](https://github.com/hbollon/igopher/releases/latest) source archive or clone the master branch
4. Launch ```bundle.sh``` script from the project root directory
5. Once done, you can find all generated executables in ```cmd/igopher/gui-bundle/output``` for all operating systems!

##### Without bundles

1. Download and install [Java 8 or 11](https://java.com/fr/download/) (needed for Selenium webdriver) and add them to your path (on Windows)
2. Install [Go](https://golang.org/doc/install) on your system
3. Download [lastest release](https://github.com/hbollon/igopher/releases/latest) source archive or clone the master branch
4. Launch it with this command: ```go run ./cmd/igopher/gui```

#### TUI version:

1. Download and install [Java 8 or 11](https://java.com/fr/download/) (needed for Selenium webdriver) and add them to your path (on Windows)
2. Install [Go](https://golang.org/doc/install) on your system
3. Download [lastest release](https://github.com/hbollon/igopher/releases/latest) source archive or clone the master branch
4. Launch it with this command: ```go run ./cmd/igopher/tui```
5. Configure the bot with your Instagram credentials and set your desired scrapping and autodm settings. To do that, you can use the TUI settings screen or directly edit the config.yaml file.
6. You're ready! Just hit the "Launch" option in the TUI main menu 🚀

### Flags

IGopher have a flags system for debuging or to enable system feature.
You can activate them by adding them after the executable call, for exemple to activate headless mode:
```./tui --headless```

There is the list of all available flags:
```
--debug
      Display debug and selenium output
--force-download
      Force redownload of all dependencies even if exists
--headless
      Run WebDriver with frame buffer
--ignore-dependencies
      Skip dependencies management
--loglevel string
      Log level threshold (default "info")
--port int
      Specify custom communication port (default 8080)
```

You can recover this list by adding **--help** flag.

## Known Issues

#### [GUI] Microsoft Smart Screen block IGopher.exe execution

At the moment Microsoft Smart Screen block IGopher.exe from launching. To avoid that, you must whitelist IGopher.
I'm currently investigating on this issue, I submitted my exe to Microsoft so we will see.

#### [GUI] Running the .app on MacOs does nothing

At the moment, you must move the .app to your Applications folder and run the binary file located in it.
It can also block the execution since the app isn't signed yet. You can avoid it by launching it from terminal or by right clicking on it and open it.

#### Javascript error just after bot launch

This issue ofter happen with an incompatible Java version installed. 
Indeed, IGopher isn't compatible with versions of the JRE greater than 11 yet due to the use of Selenium 3.

Working Java versions tested:
- Windows: [Java 8](https://java.com/fr/download/)
- Linux (Manjaro): **jre11-openjdk** -> `sudo pacman -S jre11-openjdk`


**If you find other problems, please open an issue. This is essential so that we can continue to improve IGopher! :smile:**

## 🤝 Contributing

Contributions are greatly appreciated!

1. Fork the project
2. Create your feature branch (```git checkout -b feature/AmazingFeature```)
3. Commit your changes (```git commit -m 'Add some amazing stuff'```)
4. Push to the branch (```git push origin feature/AmazingFeature```)
5. Create a new Pull Request

Issues and feature requests are welcome!
Feel free to check [issues page](https://github.com/hbollon/igopher/issues).

## Author

👤 **Hugo Bollon**

* Github: [@hbollon](https://github.com/hbollon)
* LinkedIn: [@Hugo Bollon](https://www.linkedin.com/in/hugobollon/)
* Portfolio: [hugobollon.me](https://www.hugobollon.me)

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

This project is under [MIT](https://github.com/hbollon/igopher/blob/master/LICENSE.md) license.
