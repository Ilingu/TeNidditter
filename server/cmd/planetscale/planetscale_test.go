package ps_test

import (
	"os"
	utils_env "teniditter-server/cmd/global/utils/env"
	ps "teniditter-server/cmd/planetscale"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
	utils_env.LoadEnv()
}

func TestNewDB(t *testing.T) {
	db := ps.DBManager.NewDB()
	if db == nil {
		t.Fatal("cannot connect to db")
	}
}

func TestConnectDB(t *testing.T) {
	db := ps.DBManager.Connect()
	if db == nil {
		t.Fatal("connection not opened")
	}
}
func TestDisconnectDB(t *testing.T) {
	if ok := ps.DBManager.Disconnect(); !ok {
		t.Fatal("cannot disconnect db")
	}
}
