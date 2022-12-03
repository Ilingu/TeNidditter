package utils_env_test

import (
	"os"
	"teniditter-server/cmd/global/utils"
	utils_env "teniditter-server/cmd/global/utils/env"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
}

func TestLoadEnv(t *testing.T) {
	if !utils.IsEmptyString(os.Getenv("APP_MODE")) || !utils.IsEmptyString(os.Getenv("PORT")) {
		t.Fatal("env already loaded, cannot test LoadEnv()")
	}
	utils_env.LoadEnv()
	if utils.IsEmptyString(os.Getenv("APP_MODE")) || utils.IsEmptyString(os.Getenv("PORT")) {
		t.Fatal("LoadEnv failed to load env")
	}
}
