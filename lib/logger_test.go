package lib

import "testing"

func TestNewLogger(t *testing.T) {
	l := NewLogger()
	if l.TestMode {
		t.Fatal("unexpected test mode")
	}
	if len(l.TestMessages) != 0 {
		t.Fatal("unexpected number of test messages")
	}
}

func TestLogger_Log(t *testing.T) {
	l := NewLogger()
	l.TestMode = true
	l.Log("a")
	l.Log("b")
	if l.TestMessages[0] != "a" {
		t.Fatal("expected a")
	}
	if l.TestMessages[1] != "b" {
		t.Fatal("expected b")
	}
}

func TestLogger_Fatal(t *testing.T) {
	l := NewLogger()
	l.TestMode = true
	l.Fatal("a")
	l.Fatal("b")
	if l.TestMessages[0] != "a" {
		t.Fatal("expected a")
	}
	if l.TestMessages[1] != "b" {
		t.Fatal("expected b")
	}
}

func TestLogger_Reset(t *testing.T) {
	l := NewLogger()
	l.TestMode = true
	l.Log("a")
	l.Log("b")
	if len(l.TestMessages) != 2 {
		t.Fatal("expecting 2")
	}
	l.Reset()
	if len(l.TestMessages) != 0 {
		t.Fatal("expecting 0")
	}
}
