package eventbus

import (
	"testing"
)

func TestNewTemplateConsumer(t *testing.T) {
	// A real unit test will substitute eventbus.Subscriber with a Mock.
	// In the template we mock nil so IoC resolution tests dont panic.
	// c, err := NewTemplateConsumer(nil) ...
}
