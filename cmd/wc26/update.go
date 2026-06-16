package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	updateCheck bool
	updateForce bool
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wc26 to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateCheck {
			return checkOnly()
		}
		return doUpdate()
	},
}

func latestRelease() (*githubRelease, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/Infran/wc26/releases/latest", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("parsing release: %w", err)
	}
	return &rel, nil
}

func parseVersion(tag string) string {
	return strings.TrimPrefix(tag, "v")
}

func checkOnly() error {
	rel, err := latestRelease()
	if err != nil {
		return err
	}

	latest := parseVersion(rel.TagName)
	current := parseVersion(Version)

	if latest == current {
		fmt.Fprintf(os.Stderr, "Already up to date (v%s)\n", current)
		return nil
	}
	if latest > current {
		fmt.Fprintf(os.Stderr, "Update available: v%s → v%s\n", current, latest)
	} else {
		fmt.Fprintf(os.Stderr, "Current version v%s is newer than latest v%s (dev build?)\n", current, latest)
	}
	return nil
}

func doUpdate() error {
	rel, err := latestRelease()
	if err != nil {
		return err
	}

	latest := parseVersion(rel.TagName)
	current := parseVersion(Version)

	if !updateForce && latest <= current {
		fmt.Fprintf(os.Stderr, "Already up to date (v%s)\n", current)
		return nil
	}

	if latest == current && !updateForce {
		fmt.Fprintf(os.Stderr, "Already up to date (v%s)\n", current)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Update available: v%s → v%s\n", current, latest)

	fmt.Fprintf(os.Stderr, "Show release notes? [y/N]: ")
	var showNotes string
	fmt.Scanln(&showNotes)
	if showNotes == "y" || showNotes == "Y" {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(rel.Body))
		fmt.Fprintln(os.Stderr)
	}

	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding current binary: %w", err)
	}
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return fmt.Errorf("resolving symlinks: %w", err)
	}
	dir := filepath.Dir(currentPath)

	osName := runtime.GOOS
	arch := runtime.GOARCH
	var binaryName string
	switch osName {
	case "windows":
		binaryName = fmt.Sprintf("wc26_windows_%s.exe", arch)
	default:
		binaryName = fmt.Sprintf("wc26_%s_%s", osName, arch)
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/Infran/wc26/releases/download/%s/%s",
		rel.TagName, binaryName,
	)

	fmt.Fprintf(os.Stderr, "Downloading %s ...\n", downloadURL)

	tmpFile, err := os.CreateTemp(dir, "wc26.download.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	dlResp, err := http.Get(downloadURL)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("downloading update: %w", err)
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != 200 {
		os.Remove(tmpPath)
		body, _ := io.ReadAll(io.LimitReader(dlResp.Body, 512))
		return fmt.Errorf("download failed (%d): %s", dlResp.StatusCode, strings.TrimSpace(string(body)))
	}

	written, err := io.Copy(tmpFile, dlResp.Body)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("writing download: %w", err)
	}
	if written == 0 {
		os.Remove(tmpPath)
		return fmt.Errorf("downloaded file is empty")
	}
	tmpFile.Close()

	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("setting permissions: %w", err)
	}

	backupPath := currentPath + ".old"
	if err := os.Rename(currentPath, backupPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("backing up current binary: %w", err)
	}

	if err := os.Rename(tmpPath, currentPath); err != nil {
		os.Rename(backupPath, currentPath)
		os.Remove(tmpPath)
		return fmt.Errorf("installing update: %w", err)
	}

	verify := exec.Command(currentPath, "--version")
	if out, err := verify.CombinedOutput(); err != nil {
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("new binary verification failed (%v): %s", err, string(out))
	}

	os.Remove(backupPath)

	fmt.Fprintf(os.Stderr, "Updated to %s\n", rel.TagName)
	return nil
}

func init() {
	updateCmd.Flags().BoolVarP(&updateCheck, "check", "c", false, "Only check for updates without upgrading")
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Re-download even if same version")
}
