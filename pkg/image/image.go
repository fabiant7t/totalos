package image

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type githubRelease struct {
	Name       string        `json:"name"`
	TagName    string        `json:"tag_name"`
	Draft      bool          `json:"draft"`
	Prerelease bool          `json:"prerelease"`
	Assets     []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// LatestImageURL takes the machine hardware name (architecture, the
// result of `uname -m`), reads the GitHub releases API and returns
// the latest non-draft and non-prerelease metal ISO URL.
func LatestImageURL(ctx context.Context, machineHardwareName string, client *http.Client) (string, error) {
	if client == nil {
		client = &http.Client{}
	}

	url := "https://api.github.com/repos/siderolabs/talos/releases"

	var wantName string
	switch arch := machineHardwareName; arch {
	case "x86_64":
		wantName = "metal-amd64.raw.zst"
	case "aarch64", "arm64":
		wantName = "metal-arm64.raw.zst"
	default:
		return "", fmt.Errorf("Unknown machine hardware name (architecture: %s)", arch)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var rr []githubRelease
	if err := json.Unmarshal(b, &rr); err != nil {
		return "", err
	}
	for _, r := range rr {
		if r.Draft || r.Prerelease {
			continue
		}
		for _, a := range r.Assets {
			if a.Name == wantName {
				return a.BrowserDownloadURL, nil
			}
		}
	}
	return "", errors.New("Cannot parse latest ISO")
}
