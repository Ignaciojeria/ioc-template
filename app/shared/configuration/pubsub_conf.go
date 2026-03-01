package configuration

import (
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewPubSubConfiguration)

type PubSubConfiguration struct {
	GOOGLE_PROJECT_ID string `env:"GOOGLE_PROJECT_ID"`
}

func NewPubSubConfiguration() (PubSubConfiguration, error) {
	return Parse[PubSubConfiguration]()
}
