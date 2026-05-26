package store

import "testing"

func TestProcessedStore(t *testing.T) {
	s := NewProcessedStore()

	if s.Exists("msg-1") {
		t.Fatal("message should not exist before mark")
	}

	s.MarkDone("msg-1")

	if !s.Exists("msg-1") {
		t.Fatal("message should exist after mark")
	}
}
