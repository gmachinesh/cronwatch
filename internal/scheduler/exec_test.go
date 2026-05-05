package scheduler

import (
	"strings"
	"testing"
	"time"
)

func TestSplitCommand(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"echo hello", []string{"echo", "hello"}},
		{"echo \"hello world\"", []string{"echo", "hello world"}},
		{"/usr/bin/env bash -c script.sh", []string{"/usr/bin/env", "bash", "-c", "script.sh"}},
		{"", nil},
	}

	for _, tc := range cases {
		got := splitCommand(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("splitCommand(%q): got %v, want %v", tc.input, got, tc.expected)
			continue
		}
		for i := range got {
			if got[i] != tc.expected[i] {
				t.Errorf("splitCommand(%q)[%d]: got %q, want %q", tc.input, i, got[i], tc.expected[i])
			}
		}
	}
}

func TestExecCommand_Success(t *testing.T) {
	if err := execCommand("echo cronwatch"); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestExecCommand_Failure(t *testing.T) {
	err := execCommand("false")
	if err == nil {
		t.Fatal("expected error from 'false', got nil")
	}
}

func TestExecCommand_Timeout(t *testing.T) {
	err := execCommandWithTimeout("sleep 10", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestExecCommand_NoTimeout(t *testing.T) {
	// Ensure a fast command completes successfully within a generous timeout.
	err := execCommandWithTimeout("echo cronwatch", 2*time.Second)
	if err != nil {
		t.Fatalf("expected success within timeout, got: %v", err)
	}
}

func TestExecError_Message(t *testing.T) {
	err := execCommand("ls /nonexistent_cronwatch_path_xyz")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exit status") {
		t.Errorf("unexpected error message: %v", err)
	}
}
