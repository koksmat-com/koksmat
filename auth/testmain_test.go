package auth

import (
	"os"
	"testing"

	"github.com/koksmat-com/koksmat/config"
	//"github.com/koksmat-com/koksmat/config"
)

func TestMain(m *testing.M) {

	config.Setup("../.env")

	code := m.Run()

	os.Exit(code)
}
