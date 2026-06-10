package cloudrunproxy

import (
	"fmt"
	"testing"
)

// TestMustEnv は set 済み env の素通しと、未設定時に fatalf 経路へ入ることを
// 検証する。fatalf は package variable なので process exit なしで差し替える。
func TestMustEnv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		t.Setenv("GO_CLOUDRUN_PROXY_TEST_ENV", "value-1")
		if got := MustEnv("GO_CLOUDRUN_PROXY_TEST_ENV"); got != "value-1" {
			t.Fatalf("MustEnv = %q, want %q", got, "value-1")
		}
	})

	t.Run("missing fail-closed", func(t *testing.T) {
		orig := fatalf
		defer func() { fatalf = orig }()

		var msg string
		fatalf = func(format string, v ...any) {
			msg = fmt.Sprintf(format, v...)
		}

		got := MustEnv("GO_CLOUDRUN_PROXY_TEST_ENV_MISSING")
		if got != "" {
			t.Fatalf("MustEnv = %q, want empty", got)
		}
		want := "env GO_CLOUDRUN_PROXY_TEST_ENV_MISSING is required"
		if msg != want {
			t.Fatalf("fatalf message = %q, want %q", msg, want)
		}
	})
}
