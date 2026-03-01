package postgres

import (
	"testing"
)

func TestNewTemplateRepository(t *testing.T) {
	// A proper integration test will involve setting up a test DB or using interface mocks.
	// For standard unit testing, ensuring the constructor avoids panics or nil values is paramount.
	repo, err := NewTemplateRepository(nil)
	if err != nil {
		t.Fatalf("expected no error during repository creation, got %v", err)
	}
	if repo == nil {
		t.Fatal("expected repository instance, got nil")
	}
}
