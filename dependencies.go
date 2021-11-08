package igopher

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	"sort"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/v27/github"
	log "github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"google.golang.org/api/option"
)

const (
	// desiredChromeBuild is the known build of Chromium to download from the
	// chromium-browser-snapshots/Linux_x64 bucket.
	//
	// See https://omahaproxy.appspot.com for a list of current releases.
	//
	// Update this periodically.
	desiredChromeBuild = "920003"

	// desiredFirefoxVersion is the known version of Firefox to download.
	//
	// Update this periodically.
	desiredFirefoxVersion             = "68.0.1"
	desiredSeleniumVersion            = "3.141.59"
	desiredProxyLoginAutomatorVersion = "1.0.0"
	desiredHTMLUnitDriver             = "2.54.0"
	desiredGeckodriver                = "0.30.0"
	desiredSauceLabs                  = "4.6.3"

	windowsOs = "windows"
	macOs     = "darwin"
	linuxOs   = "linux"
)

type file struct {
	URL        string   `json:"url"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Hash       string   `json:"hash"`
	HashType   string   `json:"hash_method"` // default is sha256
	Rename     []string `json:"rename,omitempty"`
	Os         string   `json:"os,omitempty"`
	compressed bool
	browser    bool
}

type downloadsTracking map[string]*downloadStatus

type downloadStatus struct {
	TotalSize      int64
	DownloadedSize int64
	Progress       float64
	Speed          int64
	Started        bool
	Completed      bool
	Failed         bool
}

var (
	files             []file
	downloadDirectory = filepath.FromSlash("./lib/")
	manifestPath      = filepath.FromSlash("./lib/manifest.json")
)

func init() {
	switch runtime.GOOS {
	case windowsOs:
		files = []file{
			{
				URL:        fmt.Sprintf("https://saucelabs.com/downloads/sc-%s-win32.zip", desiredSauceLabs),
				Name:       "sauce-connect.zip",
				Path:       downloadDirectory + "sauce-connect.zip",
				Rename:     []string{downloadDirectory + fmt.Sprintf("sc-%s-win32", desiredSauceLabs), downloadDirectory + "sauce-connect"},
				Os:         runtime.GOOS,
				compressed: true,
			},
		}

	case macOs:
		files = []file{
			{
				URL:        fmt.Sprintf("https://saucelabs.com/downloads/sc-%s-osx.zip", desiredSauceLabs),
				Name:       "sauce-connect.zip",
				Path:       downloadDirectory + "sauce-connect.zip",
				Rename:     []string{downloadDirectory + fmt.Sprintf("sc-%s-osx", desiredSauceLabs), downloadDirectory + "sauce-connect"},
				Os:         runtime.GOOS,
				compressed: true,
			},
		}

	case linuxOs:
		files = []file{
			{
				URL:        fmt.Sprintf("https://saucelabs.com/downloads/sc-%s-linux.tar.gz", desiredSauceLabs),
				Name:       "sauce-connect.tar.gz",
				Path:       downloadDirectory + "sauce-connect.tar.gz",
				Rename:     []string{downloadDirectory + fmt.Sprintf("sc-%s-linux", desiredSauceLabs), downloadDirectory + "sauce-connect"},
				Os:         runtime.GOOS,
				compressed: true,
			},
		}

	default:
		log.Fatal("Unsupported OS")
	}
}

// addGithubRelease adds a file to the list of files to download from the
// release of the specified Github repository that matches the asset
// name and tag. The file will be downloaded to localFileName.
func addGithubRelease(ctx context.Context, owner, repo, assetName, tag, localFileName string, comp bool) error {
	client := github.NewClient(nil)
	rel, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
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
			Name:       localFileName,
			Path:       downloadDirectory + localFileName,
			URL:        u,
			compressed: comp,
		})
		return nil
	}

	return fmt.Errorf("Release for %s not found at http://github.com/%s/%s/releases", assetName, owner, repo)
}

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
			Name:       localFileName,
			Path:       downloadDirectory + localFileName,
			URL:        u,
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
		downloadDriverPath         string
		targetDriverPath           string
	)

	const chromedriverzip = "chromedriver.zip"
	if runtime.GOOS == windowsOs {
		prefixOS = "Win_x64"
		lastChangeFile = "Win_x64/LAST_CHANGE"
		chromeFilename = "chrome-win.zip"
		chromeDriverFilename = "chromedriver_win32.zip"
		chromeDriverTargetFilename = chromedriverzip
		downloadDriverPath = filepath.FromSlash("chromedriver_win32/chromedriver.exe")
		targetDriverPath = "chromedriver.exe"
	} else if runtime.GOOS == macOs {
		prefixOS = "Mac"
		lastChangeFile = "Mac/LAST_CHANGE"
		chromeFilename = "chrome-mac.zip"
		chromeDriverFilename = "chromedriver_mac64.zip"
		chromeDriverTargetFilename = chromedriverzip
		downloadDriverPath = filepath.FromSlash("chromedriver_mac64/chromedriver")
		targetDriverPath = "chromedriver"
	} else {
		prefixOS = "Linux_x64"
		lastChangeFile = "Linux_x64/LAST_CHANGE"
		chromeFilename = "chrome-linux.zip"
		chromeDriverFilename = "chromedriver_linux64.zip"
		chromeDriverTargetFilename = chromedriverzip
		downloadDriverPath = filepath.FromSlash("chromedriver_linux64/chromedriver")
		targetDriverPath = "chromedriver"
	}

	gcsPath := fmt.Sprintf("gs://%s/", storageBktName)
	client, err := storage.NewClient(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return fmt.Errorf("cannot create a storage client for downloading the chrome browser: %v", err)
	}
	bkt := client.Bucket(storageBktName)
	if latestChromeBuild == "" {
		var r *storage.Reader
		r, err = bkt.Object(lastChangeFile).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("cannot create a reader for %s%s file: %v", gcsPath, lastChangeFile, err)
		}
		defer r.Close()
		// Read the last change file content for the latest build directory name
		var data []byte
		data, err = ioutil.ReadAll(r)
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
		Name:     chromeFilename,
		Path:     downloadDirectory + chromeFilename,
		browser:  true,
		URL:      cpAttrs.MediaLink,
		Hash:     hex.EncodeToString(cpAttrs.MD5),
		HashType: "md5",
	})
	latestChromeDriverPackage := path.Join(prefixOS, latestChromeBuild, chromeDriverFilename)
	cpAttrs, err = bkt.Object(latestChromeDriverPackage).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("cannot get the chrome driver package %s%s attrs: %v", gcsPath, latestChromeDriverPackage, err)
	}
	files = append(files, file{
		Name:     chromeDriverTargetFilename,
		Path:     downloadDirectory + chromeDriverTargetFilename,
		URL:      cpAttrs.MediaLink,
		Rename:   []string{downloadDirectory + downloadDriverPath, downloadDirectory + targetDriverPath},
		Hash:     hex.EncodeToString(cpAttrs.MD5),
		HashType: "md5",
	})
	return nil
}

// addFirefox adds the appropriate Firefox files to the list.
//
// If `desiredVersion` is empty, the the latest version will be used.
// Otherwise, the specific version will be used.
func addFirefox(desiredVersion string) {
	if runtime.GOOS == windowsOs {
		if desiredVersion == "" {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				URL:        "https://download.mozilla.org/?product=firefox-nightly-latest-ssl&lang=en-US",
				Name:       "firefox-nightly.exe",
				Path:       downloadDirectory + "firefox-nightly.exe",
				compressed: false,
				browser:    true,
			})
		} else {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				URL: "https://download-installer.cdn.mozilla.net/pub/firefox/releases/" +
					url.PathEscape(desiredVersion) + "/en-US/firefox-" +
					url.PathEscape(desiredVersion) + ".exe",
				Name:       "firefox.exe",
				Path:       downloadDirectory + "firefox.exe",
				compressed: false,
				browser:    true,
			})
		}
	} else {
		if desiredVersion == "" {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				URL:        "https://download.mozilla.org/?product=firefox-nightly-latest-ssl&os=linux64&lang=en-US",
				Name:       "firefox-nightly.tar.bz2",
				Path:       downloadDirectory + "firefox-nightly.tar.bz2",
				compressed: true,
				browser:    true,
			})
		} else {
			files = append(files, file{
				// This is a recent nightly. Update this path periodically.
				URL: "https://download-installer.cdn.mozilla.net/pub/firefox/releases/" +
					url.PathEscape(desiredVersion) + "/linux-x86_64/en-US/firefox-" +
					url.PathEscape(desiredVersion) + ".tar.bz2",
				Name:       "firefox.tar.bz2",
				Path:       downloadDirectory + "firefox.tar.bz2",
				compressed: true,
				browser:    true,
			})
		}
	}
}

func CheckDependencies() {
	for i := 0; i < len(files); i++ {
		log.Debugf("Checking %s with i=%d\n", files[i].Name, i)
		if f, finded := findInManifest(files[i].Name); finded {
			log.Debugf("%s is already in the manifest\n", files[i].Name)
			if fileSameHash(f) || f.Hash == "" {
				files = append(files[:i], files[i+1:]...)
				i--
			} else {
				log.Errorf("%s: hash is different from the one in the manifest, reinstallation planned\n", files[i].Name)
			}
		}
	}
}

// DownloadDependencies automate selenium dependencies downloading
// (ChromeDriver binary, the Firefox binary, the Selenium WebDriver JARs, and the Sauce Connect proxy binary)
func DownloadDependencies(downloadBrowsers, downloadLatest, forceDl bool) {
	log.Info("Downloading and installing dependencies...")
	ctx := context.Background()
	if len(files) <= 4 || files == nil {
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
	}

	if err := addGithubRelease(ctx, "SeleniumHQ", "selenium", "selenium-server-.*.jar",
		fmt.Sprintf("selenium-%s", desiredSeleniumVersion), "selenium-server.jar", false); err != nil {
		log.Errorf("Unable to find the requested Selenium Server: %s", err)
	}
	if err := addGithubRelease(ctx, "SeleniumHQ", "htmlunit-driver", "htmlunit-driver-.*-jar-with-dependencies.jar",
		desiredHTMLUnitDriver, "htmlunit-driver.jar", false); err != nil {
		log.Errorf("Unable to find the requested HTMLUnit Driver: %s", err)
	}

	if runtime.GOOS == windowsOs {
		if err := addGithubRelease(ctx, "mozilla", "geckodriver", "geckodriver-.*win64.zip",
			fmt.Sprintf("v%s", desiredGeckodriver), "geckodriver.zip", true); err != nil {
			log.Errorf("Unable to find the requested Geckodriver: %s", err)
		}
		if err := addGithubRelease(ctx, "hbollon", "proxy-login-automator",
			fmt.Sprintf("v%s", desiredProxyLoginAutomatorVersion), "proxy-login-automator-.*win64.exe",
			"proxy-login-automator.exe", false); err != nil {
			log.Errorf("Unable to find the requested proxy-login-automator: %s", err)
		}
	} else if runtime.GOOS == macOs {
		if err := addGithubRelease(ctx, "mozilla", "geckodriver",
			"geckodriver-.*macos.tar.gz", fmt.Sprintf("v%s", desiredGeckodriver), "geckodriver.tar.gz", true); err != nil {
			log.Errorf("Unable to find the requested Geckodriver: %s", err)
		}
		if err := addGithubRelease(ctx, "hbollon", "proxy-login-automator",
			fmt.Sprintf("v%s", desiredProxyLoginAutomatorVersion), "proxy-login-automator-.*macos", "proxy-login-automator", false); err != nil {
			log.Errorf("Unable to find the requested proxy-login-automator: %s", err)
		}
	} else {
		if err := addGithubRelease(ctx, "mozilla", "geckodriver", "geckodriver-.*linux64.tar.gz",
			fmt.Sprintf("v%s", desiredGeckodriver), "geckodriver.tar.gz", true); err != nil {
			log.Errorf("Unable to find the requested Geckodriver: %s", err)
		}
		if err := addGithubRelease(ctx, "hbollon", "proxy-login-automator", "proxy-login-automator-.*linux64",
			fmt.Sprintf("v%s", desiredProxyLoginAutomatorVersion), "proxy-login-automator", false); err != nil {
			log.Errorf("Unable to find the requested proxy-login-automator: %s", err)
		}
	}

	if CheckDependencies(); len(files) == 0 {
		msg := MessageOut{
			Status: SUCCESS,
			Msg:    "downloads done",
		}
		SendMessageToElectron(msg)
		log.Info("All dependencies are already installed, skipping.")
		return
	}

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)

	var wg sync.WaitGroup
	var filesToDl []file
	dlTracking := make(downloadsTracking)
	for _, file := range files {
		if file.Os == "" || file.Os == runtime.GOOS {
			wg.Add(1)
			dlTracking[file.Name] = &downloadStatus{}
			filesToDl = append(filesToDl, file)
			bar := p.Add(0,
				mpb.NewBarFiller("[=>-|"),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.OnComplete(decor.Name(file.Name+": ", decor.WCSyncSpaceR), file.Name+": done!"),
					decor.OnComplete(decor.CountersKibiByte("% .2f / % .2f", decor.WCSyncWidth), ""),
				),
				mpb.AppendDecorators(
					decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 90, decor.WCSyncWidth), ""),
					decor.OnComplete(decor.Name(" ] "), ""),
					decor.OnComplete(decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60, decor.WCSyncWidth), ""),
				),
			)
			file := file
			go func() {
				time.Sleep(2 * time.Second)
				if err := handleFile(bar, dlTracking, file, downloadBrowsers, forceDl); err != nil {
					log.Fatalf("Error handling %s: %s", file.Name, err)
				}
				wg.Done()
			}()
		}
	}
	if IsElectronRunning() {
		done := make(chan bool)
		go followUpDownloads(dlTracking, filesToDl, done)
		defer func(done chan bool) {
			done <- true
			close(done)
		}(done)
	}
	p.Wait()
	wg.Wait()
}

func handleFile(bar *mpb.Bar, dlTracking downloadsTracking, file file, downloadBrowsers, forceDl bool) error {
	if file.browser && !downloadBrowsers {
		log.Infof("Skipping %q because --download_browser is not set.", file.Name)
		bar.Abort(true)
		return nil
	}
	if _, err := os.Stat(file.Path); err == nil && !forceDl {
		log.Debugf("Skipping file %q which has already been downloaded.", file.Name)
		bar.Abort(true)
	} else {
		if err := downloadFile(bar, dlTracking, file); err != nil {
			bar.Abort(true)
			return err
		}
	}

	if err := extractFile(file); err != nil {
		return err
	}
	if rename := file.Rename; len(rename) == 2 {
		log.Debugf("Renaming %q to %q", rename[0], rename[1])
		os.RemoveAll(rename[1]) // Ignore error.
		if err := os.Rename(rename[0], rename[1]); err != nil {
			log.Warnf("Error renaming %q to %q: %v", rename[0], rename[1], err)
		}
	}
	if err := dumpDependencyToManifest(file); err != nil {
		log.Error(err)
	}
	return nil
}

func extractFile(file file) error {
	switch path.Ext(file.Name) {
	case ".zip":
		log.Debugf("Unzipping %q", file.Path)
		if runtime.GOOS == windowsOs {
			if err := exec.Command("tar", "-xf", file.Path, "-C", downloadDirectory).Run(); err != nil {
				return fmt.Errorf("Error unzipping %q: %v", file.Path, err)
			}
		} else {
			if err := exec.Command("unzip", "-o", file.Path, "-d", downloadDirectory).Run(); err != nil {
				return fmt.Errorf("Error unzipping %q: %v", file.Path, err)
			}
		}
	case ".gz":
		log.Debugf("Unzipping %q", file.Path)
		if err := exec.Command("tar", "-xzf", file.Path, "-C", downloadDirectory).Run(); err != nil {
			return fmt.Errorf("Error unzipping %q: %v", file.Path, err)
		}
	case ".bz2":
		log.Debugf("Unzipping %q", file.Path)
		if err := exec.Command("tar", "-xjf", file.Path, "-C", downloadDirectory).Run(); err != nil {
			return fmt.Errorf("Error unzipping %q: %v", file.Path, err)
		}
	}

	return nil
}

func findInManifest(name string) (file, bool) {
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return file{}, false
	}

	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Error("Failed to read existing manifest file.")
		return file{}, false
	}

	data := []file{}
	json.Unmarshal(manifest, &data)
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name <= data[j].Name
	})
	idx := sort.Search(len(data), func(i int) bool {
		return data[i].Name >= name
	})
	if idx < len(data) && data[idx].Name == name {
		return data[idx], true
	}
	return file{}, false
}

func dumpDependencyToManifest(f file) error {
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		_, err := os.Create(manifestPath)
		if err != nil {
			return err
		}
	}

	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	data := []file{}
	json.Unmarshal(manifest, &data)
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name <= data[j].Name
	})
	idx := sort.Search(len(data), func(i int) bool {
		return data[i].Name >= f.Name
	})
	if idx < len(data) && data[idx].Name == f.Name {
		data[idx] = f
	} else {
		data = append(data, f)
	}

	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(manifestPath, dataBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(bar *mpb.Bar, dlTracking downloadsTracking, file file) (err error) {
	f, err := os.Create(file.Path)
	if err != nil {
		return fmt.Errorf("error creating %q: %v", file.Path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing %q: %v", file.Path, err)
		}
	}()

	resp, err := http.Get(file.URL)
	if err != nil {
		return fmt.Errorf("%s: error downloading %q: %v", file.Name, file.URL, err)
	}
	defer resp.Body.Close()

	bar.SetTotal(resp.ContentLength, false)
	if track, ok := dlTracking[file.Name]; ok {
		track.TotalSize = resp.ContentLength
	}

	if file.Hash != "" {
		var h hash.Hash
		switch strings.ToLower(file.HashType) {
		case "md5":
			h = md5.New()
		case "sha1":
			h = sha1.New()
		default:
			h = sha256.New()
		}
		dlTracking[file.Name].Started = true
		if _, err := io.Copy(io.MultiWriter(f, h), bar.ProxyReader(resp.Body)); err != nil {
			return fmt.Errorf("%s: error downloading %q: %v", file.Name, file.URL, err)
		}
		if h := hex.EncodeToString(h.Sum(nil)); h != file.Hash {
			return fmt.Errorf("%s: got %s hash %q, want %q", file.Name, file.HashType, h, file.Hash)
		}
	} else {
		dlTracking[file.Name].Started = true
		if _, err := io.Copy(f, bar.ProxyReader(resp.Body)); err != nil {
			return fmt.Errorf("%s: error downloading %q: %v", file.Name, file.URL, err)
		}
	}
	return nil
}

// Will regulrary check downloads status and calculate progress pencentages
// to send it to Electron GUI if it's running
func followUpDownloads(dlTracking downloadsTracking, srcFiles []file, done chan bool) {
	time.Sleep(2 * time.Second)
	for {
		select {
		case <-done:
			msg := MessageOut{
				Status: SUCCESS,
				Msg:    "downloads done",
			}
			SendMessageToElectron(msg)
			log.Infof("Downloads finished")
			return
		default:
			for _, srcFile := range srcFiles {
				if track, ok := dlTracking[srcFile.Name]; ok {
					if track.Started && !track.Completed && !track.Failed {
						file, err := os.Open(srcFile.Path)
						if err != nil {
							log.Error(err)
						}

						fi, err := file.Stat()
						if err != nil {
							log.Error(err)
						}

						size := fi.Size()
						if size == 0 {
							size = 1
						}

						track.DownloadedSize = size
						track.Progress = float64(size) / float64(track.TotalSize) * 100
						if track.Progress == 100 {
							track.Completed = true
						}
					}
				} else {
					log.Errorf("%s: download tracking not found", srcFile.Name)
				}
			}

			msg := MessageOut{
				Status:  INFO,
				Msg:     "downloads tracking",
				Payload: dlTracking,
			}
			SendMessageToElectron(msg)
		}
		time.Sleep(2 * time.Second)
	}
}

func fileSameHash(file file) bool {
	if file.Hash == "" {
		return false
	}
	if _, err := os.Stat(file.Path); err != nil {
		return false
	}

	var h hash.Hash
	switch strings.ToLower(file.HashType) {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	default:
		h = sha256.New()
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return false
	}

	sum := hex.EncodeToString(h.Sum(nil))
	if sum != file.Hash {
		log.Warningf("File %q: got hash %q, expect hash %q", file.Path, sum, file.Hash)
		return false
	}
	return true
}
