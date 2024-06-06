package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestShipGet(t *testing.T) {

	shipGetCmd(nil, []string{"test"})

}

func TestInstall(t *testing.T) {

	err := Install()
	if err != nil {
		t.Error(err)
	}

}

func TestBuildUrl(t *testing.T) {

	url := BuildUrlFromPackage("magicbutton/magic-mix:0.0.3.3")
	assert.Exactly(t, url, "https://github.com/magicbutton/magic-mix/archive/refs/tags/v0.0.3.3.zip")
}

func TestSetupConnectorsLocal(t *testing.T) {

	err := SyncConnectorsWithMaster("./testroot")
	if err != nil {
		t.Error(err)
	}

}

func TestSetupConnectorsRoot(t *testing.T) {
	root := viper.GetString("KITCHENROOT")

	err := SyncConnectorsWithMaster(root)
	if err != nil {
		t.Error(err)
	}

}
