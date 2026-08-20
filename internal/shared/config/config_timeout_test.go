package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadRejectsNegativeRequestTimeout(t *testing.T) {
	t.Setenv("SIGNALFORGE_REQUEST_TIMEOUT", "250ms")
	cfg, err := Load("")
	if err != nil || cfg.RuntimeTimeouts().Request != 250*time.Millisecond {
		t.Fatalf("request timeout override did not reach runtime: cfg=%+v err=%v", cfg, err)
	}
	t.Setenv("SIGNALFORGE_REQUEST_TIMEOUT", "-2s")
	_, e := Load("")
	if e == nil || !strings.Contains(e.Error(), "duration must be positive") {
		t.Fatal(e)
	}
}
func TestLoadRejectsNegativeReadTimeout(t *testing.T) {
	t.Setenv("SIGNALFORGE_READ_TIMEOUT", "350ms")
	cfg, err := Load("")
	if err != nil || cfg.RuntimeTimeouts().Read != 350*time.Millisecond {
		t.Fatalf("read timeout override did not reach runtime: cfg=%+v err=%v", cfg, err)
	}
	t.Setenv("SIGNALFORGE_READ_TIMEOUT", "-2s")
	_, e := Load("")
	if e == nil || !strings.Contains(e.Error(), "duration must be positive") {
		t.Fatal(e)
	}
}
func TestLoadRejectsNegativeWriteTimeout(t *testing.T) {
	t.Setenv("SIGNALFORGE_WRITE_TIMEOUT", "450ms")
	cfg, err := Load("")
	if err != nil || cfg.RuntimeTimeouts().Write != 450*time.Millisecond {
		t.Fatalf("write timeout override did not reach runtime: cfg=%+v err=%v", cfg, err)
	}
	t.Setenv("SIGNALFORGE_WRITE_TIMEOUT", "-2s")
	_, e := Load("")
	if e == nil || !strings.Contains(e.Error(), "duration must be positive") {
		t.Fatal(e)
	}
}
