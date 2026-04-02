package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const checkInterval = 24 * time.Hour

// VersionCheck holds the cached result of a version check
type VersionCheck struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// VersionChecker performs async version checks with caching
type VersionChecker struct {
	currentVersion string
	cacheDir       string
	result         *VersionCheck
	done           chan struct{}
}

// NewVersionChecker creates a checker that runs in the background
func NewVersionChecker(currentVersion, homeDir string) *VersionChecker {
	return &VersionChecker{
		currentVersion: currentVersion,
		cacheDir:       filepath.Join(homeDir, ".magebox"),
		done:           make(chan struct{}),
	}
}

func (vc *VersionChecker) cachePath() string {
	return filepath.Join(vc.cacheDir, "version-check.json")
}

// Start begins the version check in the background
func (vc *VersionChecker) Start() {
	go func() {
		defer close(vc.done)

		// Try to load from cache first
		if cached, err := vc.loadCache(); err == nil {
			if time.Since(cached.CheckedAt) < checkInterval {
				vc.result = cached
				return
			}
		}

		// Cache is stale or missing — fetch from GitHub with a short timeout
		client := &http.Client{Timeout: 3 * time.Second}
		url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", GitHubAPIURL, GitHubOwner, GitHubRepo)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "MageBox-VersionCheck")

		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		var release GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return
		}

		check := &VersionCheck{
			LatestVersion: release.TagName,
			CheckedAt:     time.Now(),
		}
		vc.result = check
		_ = vc.saveCache(check)
	}()
}

// Result waits for the check to complete and returns a message if an update is available.
// Returns empty string if no update or if the check failed.
func (vc *VersionChecker) Result() string {
	<-vc.done

	if vc.result == nil {
		return ""
	}

	u := &Updater{currentVersion: vc.currentVersion}
	if !u.isNewerVersion(vc.result.LatestVersion) {
		return ""
	}

	return fmt.Sprintf("A new version of MageBox is available: %s (current: %s). Run 'magebox self-update' to upgrade.",
		vc.result.LatestVersion, vc.currentVersion)
}

func (vc *VersionChecker) loadCache() (*VersionCheck, error) {
	data, err := os.ReadFile(vc.cachePath())
	if err != nil {
		return nil, err
	}
	var check VersionCheck
	if err := json.Unmarshal(data, &check); err != nil {
		return nil, err
	}
	return &check, nil
}

func (vc *VersionChecker) saveCache(check *VersionCheck) error {
	_ = os.MkdirAll(vc.cacheDir, 0755)
	data, err := json.Marshal(check)
	if err != nil {
		return err
	}
	return os.WriteFile(vc.cachePath(), data, 0644)
}
