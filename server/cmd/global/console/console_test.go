package console_test

import (
	"io"
	"log"
	"os"
	"teniditter-server/cmd/global/console"
	"testing"
)

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

func TestLogMsg(t *testing.T) {
	res := hijackStdout(func() {
		console.LogMsg("test", console.NEUTRAL)
	})
	log.Println(res)
}
