package app

import (
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/chigopher/pathlib"
)

// test the fileVersion
func TestFileVersion(t *testing.T) {
	ver := "1.2.3"
	// Set a test version
	fileVersion(ver)
	// Get the version
	got := fileVersion()

	if got != ver {
		t.Errorf("fileVersion() = %q, want %q", got, ver)
	}

	LatestVersion.Fetch()

	time.Sleep(1 * time.Second)
	got, _ = LatestVersion.Get(1)
	if got != ver {
		t.Errorf("LatestVersion.Get(1) = %q, want %q", got, ver)
	}

	// Test the expiration
	// set the filetime to 2 hours ago
	p := pathlib.NewPath(os.TempDir()).Join(".labctl.version")
	err := os.Chtimes(p.String(), time.Now().Add(-2*time.Hour), time.Now().Add(-2*time.Hour))
	assert.NoError(t, err)

	got = fileVersion()
	if got != "" {
		t.Errorf("fileVersion() = %q, want %q", got, "")
	}
}
