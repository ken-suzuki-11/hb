package libs

import (
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type Config struct {
	Common   CommonConfig   `toml:"common" validate:"required"`
	Http     HttpConfig     `toml:"http" validate:"required"`
	Function FunctionConfig `toml:"function" validate:"required"`
	Parallel ParallelConfig `toml:"parallel" validate:"required"`
}

type CommonConfig struct {
	ListSizeLimit int64 `toml:"list_size_limit" validate:"required"`
	Debug         bool  `toml:"debug" validate:"required"`
	Timeout       int   `toml:"timeout" validate:"required"`
}
type HttpConfig struct {
	DisableKeepalive bool `toml:"disable_keepalive" validate:"required"`
}
type FunctionConfig struct {
	Name string `toml:"name" validate:"required"`
}
type ParallelConfig struct {
	ProcessNum int `toml:"process_num" validate:"required"`
}

func NewConfig(configFile string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, errors.Wrap(err, "toml decode")
	}
	return &config, nil
}

func (c Config) Validation() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return errors.Wrap(err, "config validation")
	}
	return nil
}
