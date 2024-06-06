package repos

import (
	"testing"
)

func TestDownload(t *testing.T) {

	DownloadRepo("magicbutton", "magic-master")
}
