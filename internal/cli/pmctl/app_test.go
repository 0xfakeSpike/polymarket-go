package pmctl

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunTools(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := App{Stdout: &stdout, Stderr: &stderr}

	if code := app.Run([]string{"tools"}); code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"client_call"`) {
		t.Fatalf("expected tool list, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"get_markets_by_annualized_return"`) {
		t.Fatalf("expected get_markets_by_annualized_return tool, got %s", stdout.String())
	}
}

func TestRunToolMethods(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := App{Stdout: &stdout, Stderr: &stderr}

	code := app.Run([]string{"tool", "-params", `{"long":false}`, "methods"})
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"GetOK"`) {
		t.Fatalf("expected method list, got %s", stdout.String())
	}
}

func TestRunRejectsOldBusinessShortcut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := App{Stdout: &stdout, Stderr: &stderr}

	if code := app.Run([]string{"orderbook"}); code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "pmctl tool") {
		t.Fatalf("expected generic tool usage, got %s", stderr.String())
	}
}
