// Binary init downloads the necessary files to perform an integration test
// between this WebDriver client and multiple versions of Selenium and
// browsers.

package igopher

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/v27/github"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

const (
	// desiredChromeBuild is the known build of Chromium to download from the
	// chromium-browser-snapshots/Linux_x64 bucket.
	//
	// See https://omahaproxy.appspot.com for a list of current releases.
	//
	// Update this periodically.
	desiredChromeBuild = "664981" // This corresponds to version 76.0.3809.0

	// desiredFirefoxVersion is the known version of Firefox to download.
	//
	// Update this periodically.
	desiredFirefoxVersion = "68.0.1"
)

type file struct {
	url        string
	name       string
	path       string
	hash       string
	hashType   string // default is sha256
	rename     []string
	os         string
	compressed bool
	browser    bool
}

var (
	files = []file{
		{
			url:        "https://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar",
			name:       "selenium-server.jar",
			path:       downloadDirectory + "selenium-server.jar",
			compressed: false,
		},
		{
			url:        "https://saucelabs.com/downloads/sc-4.6.3-linux.tar.gz",
			name:       "sauce-connect.tar.gz",
			path:       downloadDirectory + "sauce-connect.tar.gz",
			rename:     []string{downloadDirectory + "sc-4.6.3-linux", downloadDirectory + "sauce-connect"},
			os:         "linux",
			compressed: true,
		},
		{
			url:        "https://saucelabs.com/downloads/sc-4.6.3-win32.zip",
			name:       "sauce-connect.zip",
			path:       downloadDirectory + "sauce-connect.zip",
			rename:     []string{downloadDirectory + "sc-4.6.3-win32", downloadDirectory + "sauce-connect"},
			os:         "windows",
			compressed: true,
		},
	}

	downloadDirectory = filepath.FromSlash("./lib/")
)

// addLatestGithubRelease adds a file to the list of files to download from the
// latest release of the specified Github repository that matches the asset
// name. The file will be downloaded to localFileName.
func addLatestGithubRelease(ctx context.Context, owner, repo, assetName, localFileName string, comp bool) error {
	client := github.NewClient(nil)

	rel, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return err
	}
	assetNameRE, err := regexp.Compile(assetName)
	if err != nil {
		return fmt.Errorf("invalid asset name regular expression %q: %s", assetName, err)
	}
	for _, a := range rel.Assets {
		if !assetNameRE.MatchString(a.GetName()) {
			continue
		}
		u := a.GetBrowserDownloadURL()
		if u == "" {
			return fmt.Errorf("%s does not have a download URL", a.GetName())
		}
		files = append(files, file{
			name:       localFileName,
			path:       downloadDirectory + localFileName,
			url:        u,
			compressed: comp,
		})
		return nil
	}

	return fmt.Errorf("Release for %s not found at http://github.com/%s/%s/releases", assetName, owner, repo)
}

// addChrome adds the appropriate chromium files to the list.
//
// If `latestChromeBuild` is empty, then the latest build will be used.
// Otherwise, that specific build will be used.
func addChrome(ctx context.Context, latestChromeBuild string) error {
	// Bucket URL: https://console.cloud.google.com/storage/browser/chromium-browser-continuous/?pli=1
	const storageBktName = "chromium-browser-snapshots"
	var (
		lastChangeFile             string
		prefixOS                   string
		chromeFilename             string
		chromeDriverFilename       string
		chromeDriverTargetFilename string // For backward compatibility
	)

	if runtime.GOOS == "windows" {
		prefixOS = "Win_x64"
		lastChangeFile = "Win_x64/LAST_CHANGE"
		chromeFilename = "chrome-win.zip"
		chromeDriverFilename = "chromedriver_win32.zip"
		chromeDriverTargetFilename = "chromedriver.zip"
	} else {
		prefixOS = "Linux_x64"
		lastChangeFile = "Linux_x64/LAST_CHANGE"
		chromeFilename = "chrome-linux.zip"
		chromeDriverFilename = "chromedriver_linux64.zip"
		chromeDriverTargetFilename = "chromedriver.zip"
	}

	gcsPath := fmt.Sprintf("gs://%s/", storageBktName)
	client, err := storage.NewClient(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return fmt.Errorf("cannot create a storage client for downloading the chrome browser: %v", err)
	}
	bkt := client.Bucket(storageBktName)
	if latestChromeBuild == "" {
		r, err := bkt.Object(lastChangeFile).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("cannot create a reader for %s%s file: %v", gcsPath, lastChangeFile, err)
		}
		defer r.Close()
		// Read the last change file content for the latest build directory name
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return fmt.Errorf("cannot read from %s%s file: %v", gcsPath, lastChangeFile, err)
		}
		latestChromeBuild = string(data)
	}
	latestChromePackage := path.Join(prefixOS, latestChromeBuild, chromeFilename)
	cpAttrs, err := bkt.Object(latestChromePackage).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("cannot get the chrome package %s%s attrs: %v", gcsPath, latestChromePackage, err)
	}
	files = append(files, file{
		name:    chromeFilename,
		path:    downloadDirectory + chromeFilename,
		browser: true,
		url:     cpAttrs.MediaLink,
	})
	latestChromeDriverPackage := path.Join(prefixOS, latestChromeBuild, chromeDriverFilename)
	cpAttrs, err = bkt.Object(latestChromeDriverPackage).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("cannot get the chrome driver package %s%s attrs: %v", gcsPath, latestChromeDriverPackage, err)
	}
	files = append(files, file{
		name:   chromeDriverTargetFilename,
		path:   downloadDirectory + chromeDriverTargetFilename,
		url:    cpAttrs.MediaLink,
		rename: []string{downloadDirectory + "chromedriver_linux64/chromedriver", downloadDirectory + "chromedriver"},
	})
	return nil
}

// addFirefox adds the appropriate Firefox files to the list.
//
// If `desiredVersion` is empty, the the latest version will be used.
// Otherwise, the specific version will be used.
func addFirefox(desiredVersion string) {
	if runtime.GOOS == "windows" {
		if desiredVersion == "" {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				url:        "https://download.mozilla.org/?product=firefox-nightly-latest-ssl&lang=en-US",
				name:       "firefox-nightly.exe",
				path:       downloadDirectory + "firefox-nightly.exe",
				compressed: false,
				browser:    true,
			})
		} else {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				url:        "https://download-installer.cdn.mozilla.net/pub/firefox/releases/" + url.PathEscape(desiredVersion) + "/en-US/firefox-" + url.PathEscape(desiredVersion) + ".exe",
				name:       "firefox.exe",
				path:       downloadDirectory + "firefox.exe",
				compressed: false,
				browser:    true,
			})
		}
	} else {
		if desiredVersion == "" {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				url:        "https://download.mozilla.org/?product=firefox-nightly-latest-ssl&os=linux64&lang=en-US",
				name:       "firefox-nightly.tar.bz2",
				path:       downloadDirectory + "firefox-nightly.tar.bz2",
				compressed: true,
				browser:    true,
			})
		} else {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				url:        "https://download-installer.cdn.mozilla.net/pub/firefox/releases/" + url.PathEscape(desiredVersion) + "/linux-x86_64/en-US/firefox-" + url.PathEscape(desiredVersion) + ".tar.bz2",
				name:       "firefox.tar.bz2",
				path:       downloadDirectory + "firefox.tar.bz2",
				compressed: true,
				browser:    true,
			})
		}
	}
}

// DownloadDependencies automate selenium dependencies downloading
// (ChromeDriver binary, the Firefox binary, the Selenium WebDriver JARs, and the Sauce Connect proxy binary)
func DownloadDependencies(downloadBrowsers, downloadLatest, forceDl bool) {
	log.Info("Downloading and installing dependencies...")
	ctx := context.Background()
	if downloadBrowsers {
		chromeBuild := desiredChromeBuild
		firefoxVersion := desiredFirefoxVersion
		if downloadLatest {
			chromeBuild = ""
			firefoxVersion = ""
		}

		if err := addChrome(ctx, chromeBuild); err != nil {
			log.Errorf("Unable to download Google Chrome browser: %v", err)
		}
		addFirefox(firefoxVersion)
	}

	if err := addLatestGithubRelease(ctx, "SeleniumHQ", "htmlunit-driver", "htmlunit-driver-.*-jar-with-dependencies.jar", "htmlunit-driver.jar", false); err != nil {
		log.Errorf("Unable to find the latest HTMLUnit Driver: %s", err)
	}

	if runtime.GOOS == "windows" {
		if err := addLatestGithubRelease(ctx, "mozilla", "geckodriver", "geckodriver-.*win64.zip", "geckodriver.zip", true); err != nil {
			log.Errorf("Unable to find the latest Geckodriver: %s", err)
		}
	} else {
		if err := addLatestGithubRelease(ctx, "mozilla", "geckodriver", "geckodriver-.*linux64.tar.gz", "geckodriver.tar.gz", true); err != nil {
			log.Errorf("Unable to find the latest Geckodriver: %s", err)
		}
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if file.os == "" || file.os == runtime.GOOS {
			wg.Add(1)
			file := file
			go func() {
				if err := handleFile(file, downloadBrowsers, forceDl); err != nil {
					log.Fatalf("Error handling %s: %s", file.name, err)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func handleFile(file file, downloadBrowsers, forceDl bool) error {
	if file.browser && !downloadBrowsers {
		log.Infof("Skipping %q because --download_browser is not set.", file.name)
		return nil
	}
	if _, err := os.Stat(file.path); err == nil && !forceDl {
		log.Debugf("Skipping file %q which has already been downloaded.", file.name)
	} else {
		log.Debugf("Downloading %q from %q", file.name, file.url)
		if err := downloadFile(file); err != nil {
			return err
		}
	}

	switch path.Ext(file.name) {
	case ".zip":
		log.Debugf("Unzipping %q", file.path)
		if runtime.GOOS == "windows" {
			if err := exec.Command("tar", "-xf", file.path, "-C", downloadDirectory).Run(); err != nil {
				return fmt.Errorf("Error unzipping %q: %v", file.path, err)
			}
		} else {
			if err := exec.Command("unzip", "-o", file.path, "-d", downloadDirectory).Run(); err != nil {
				return fmt.Errorf("Error unzipping %q: %v", file.path, err)
			}
		}
	case ".gz":
		log.Debugf("Unzipping %q", file.path)
		if err := exec.Command("tar", "-xzf", file.path, "-C", downloadDirectory).Run(); err != nil {
			return fmt.Errorf("Error unzipping %q: %v", file.path, err)
		}
	case ".bz2":
		log.Debugf("Unzipping %q", file.path)
		if err := exec.Command("tar", "-xjf", file.path, "-C", downloadDirectory).Run(); err != nil {
			return fmt.Errorf("Error unzipping %q: %v", file.path, err)
		}
	}
	if rename := file.rename; len(rename) == 2 {
		log.Debugf("Renaming %q to %q", rename[0], rename[1])
		os.RemoveAll(rename[1]) // Ignore error.
		if err := os.Rename(rename[0], rename[1]); err != nil {
			log.Warnf("Error renaming %q to %q: %v", rename[0], rename[1], err)
		}
	}
	return nil
}

func downloadFile(file file) (err error) {
	f, err := os.Create(file.path)
	if err != nil {
		return fmt.Errorf("error creating %q: %v", file.path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing %q: %v", file.path, err)
		}
	}()

	resp, err := http.Get(file.url)
	if err != nil {
		return fmt.Errorf("%s: error downloading %q: %v", file.name, file.url, err)
	}
	defer resp.Body.Close()
	if file.hash != "" {
		var h hash.Hash
		switch strings.ToLower(file.hashType) {
		case "md5":
			h = md5.New()
		case "sha1":
			h = sha1.New()
		default:
			h = sha256.New()
		}
		if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
			return fmt.Errorf("%s: error downloading %q: %v", file.name, file.url, err)
		}
		if h := hex.EncodeToString(h.Sum(nil)); h != file.hash {
			return fmt.Errorf("%s: got %s hash %q, want %q", file.name, file.hashType, h, file.hash)
		}
	} else {
		if _, err := io.Copy(f, resp.Body); err != nil {
			return fmt.Errorf("%s: error downloading %q: %v", file.name, file.url, err)
		}
	}
	return nil
}

func fileSameHash(file file) bool {
	if _, err := os.Stat(file.path); err != nil {
		return false
	}
	var h hash.Hash
	switch strings.ToLower(file.hashType) {
	case "md5":
		h = md5.New()
	default:
		h = sha256.New()
	}
	f, err := os.Open(file.path)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return false
	}

	sum := hex.EncodeToString(h.Sum(nil))
	if sum != file.hash {
		log.Warningf("File %q: got hash %q, expect hash %q", file.path, sum, file.hash)
		return false
	}
	return true
}
