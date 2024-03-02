package kitchen

import (
	"testing"
)

func TestEnvFileChech(t *testing.T) {

	s, err := envFileCheck()
	if err != nil {
		t.Error(err)
	}
	t.Log(s)

}

func TestCreateDefaultEnvFile(t *testing.T) {

	_, err := createDefaultEnvFile(".env")
	if err != nil {
		t.Error(err)
	}

}

func TestCreateNoneDefaultEnvFile(t *testing.T) {

	_, err := createDefaultEnvFile(".xenv")
	if err != nil {
		t.Error(err)
	}

}

func TestBoot(t *testing.T) {
	err := Boot(BootOptions{Verbose: true})
	if err != nil {
		t.Error(err)
	}

}
