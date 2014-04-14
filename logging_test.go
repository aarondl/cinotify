package cinotify

import (
	"bytes"
	"log"
	"strings"
	. "testing"
)

func TestLogger(t *T) {
	// Ensure no panics
	DoLog("")
	DoLogf("")

	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)

	DoLog("test")
	if s := strings.TrimSpace(buf.String()); s != "test" {
		t.Error("Expected test but got: ", s)
	}

	buf.Reset()
	DoLogf("te%dt", 5)
	if s := strings.TrimSpace(buf.String()); s != "te5t" {
		t.Error("Expected te5t but got: ", s)
	}
}
