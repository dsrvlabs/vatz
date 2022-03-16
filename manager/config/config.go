package config

type config struct {
}

func (c config) Parse() error {
	return nil
}

type Config interface {
	Parse() error
}

func NewConfig() Config {
	return &config{}
}
