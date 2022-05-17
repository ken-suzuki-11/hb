package libs

import (
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type Config struct {
	Common     CommonConfig     `toml:"common" validate:"required"`
	Parallel ParallelConfig `toml:"parallel" validate:"required"`
}
type CommonConfig struct {
	Function      string `toml:"function" validate:"required"`
	ListSizeLimit int64  `toml:"list_size_limit" validate:"required"`
	Debug         bool   `toml:"debug" validate:"required"`
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
