package main

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
	os.Setenv("PORT", "3000")
}

// High order function
func hijackStdout(toExec func()) string {
	// save std's
	stdOut := os.Stdout
	stdErr := os.Stderr

	// hijack std's
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	log.SetOutput(w)

	// Exec Funcs
	toExec()

	// Close hijack
	_ = w.Close()

	// get results
	result, _ := io.ReadAll(r)
	output := string(result)

	// reset the std's
	os.Stdout = stdOut
	os.Stderr = stdErr
	log.SetOutput(os.Stderr)

	return output
}

func TestServer(t *testing.T) {
	res := hijackStdout(func() {
		main()
	})

	wantedLog := []string{"Cors Middleware Up and Running", "AuthHandler Registered", "TedinitterUserHandler Registered", "TedditHandler Registered", "NitterHandler Registered", "http server started on"}
	for _, log := range wantedLog {
		if !strings.Contains(res, log) {
			t.Error("server start failed, missing", log)
		}
	}
}
