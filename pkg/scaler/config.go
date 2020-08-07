package scaler

import (
	"fmt"
	"time"
)

type ScaleConfig struct {
	TimeOut int `yaml:"timeOut"`
}

func (cfg *ScaleConfig) IsValid() error {
	if cfg.TimeOut <= 0 {
		return fmt.Errorf("Config scale timeOut is invalid")
	}

	return nil
}

func (cfg *ScaleConfig) GetScaleDuration() time.Duration {
	return time.Duration(cfg.TimeOut) * time.Minute
}
