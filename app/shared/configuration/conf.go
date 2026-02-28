package configuration

import (
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewConf)

type Conf struct {
	PORT         string `env:"PORT" envDefault:"8080"`
	PROJECT_NAME string `env:"PROJECT_NAME"`
	VERSION      string `env:"VERSION"`
}

// NewConf loads the configuration and provides it.
// It is returned by value because it's lightweight and immutable.
func NewConf() (Conf, error) {
	return Parse[Conf]()
}
