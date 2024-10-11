package command

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// InstallRawImage downloads the raw image URL and writes it to the given device.
func InstallRawImage(m remotecommand.Machine, rawImageURL, device string, cb ssh.HostKeyCallback) error {
	parsedURL, err := url.Parse(rawImageURL)
	if err != nil {
		return err
	}

	switch ext := filepath.Ext(parsedURL.Path); ext {
	case ".xz":
		return installRawImageXZ(m, rawImageURL, device, cb)
	case ".zst":
		return installImageZstandard(m, rawImageURL, device, cb)
	default:
		return fmt.Errorf("InstallRawImage canot handle a %s file", ext)
	}
}

// installRawImageXZ downloads the raw.xz image URL and writes it to the given device.
func installRawImageXZ(m remotecommand.Machine, rawImageURL, device string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
    wget %s -O talos-metal.raw.xz \
    && cat talos-metal.raw.xz \
    |  xz -d \
    |  dd of=%s bs=4M \
    && sync`, rawImageURL, device)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command InstallImage failed: %w", err)
	}
	return nil
}

// installImageZstandard downloads the raw.zst image URL and writes it to the given device.
func installImageZstandard(m remotecommand.Machine, rawImageURL, device string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
    wget %s -O talos-metal.raw.zst \
    && cat talos-metal.raw.zst \
    |  zstd -d \
    |  dd of=%s bs=4M \
    && sync`, rawImageURL, device)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command InstallImage failed: %w", err)
	}
	return nil
}
