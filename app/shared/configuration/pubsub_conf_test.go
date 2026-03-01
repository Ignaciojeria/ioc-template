package configuration

import (
	"os"
	"testing"
)

func TestNewPubSubConfiguration(t *testing.T) {
	os.Setenv("GOOGLE_PROJECT_ID", "test-project-99")

	defer func() {
		os.Unsetenv("GOOGLE_PROJECT_ID")
	}()

	conf, err := NewPubSubConfiguration()
	if err != nil {
		t.Fatalf("unexpected error parsing pubsub config: %v", err)
	}

	if conf.GOOGLE_PROJECT_ID != "test-project-99" {
		t.Errorf("expected test-project-99 got %s", conf.GOOGLE_PROJECT_ID)
	}
}
