package osd_cluster_upgrader

import (
	"fmt"
	"time"

	"github.com/openshift/managed-upgrade-operator/pkg/maintenance"
)

type osdUpgradeConfig struct {
	Maintenance maintenance.MaintenanceConfig `yaml:"maintenance"`
	Scale       scaleConfig                   `yaml:"scale"`
	NodeDrain   nodeDrain                     `yaml:"nodeDrain"`
}

type scaleConfig struct {
	TimeOut int `yaml:"timeOut"`
}

type nodeDrain struct {
	TimeOut int `yaml:"timeOut"`
}

func (cfg *osdUpgradeConfig) IsValid() error {
	if err := cfg.Maintenance.IsValid(); err != nil {
		return err
	}
	if cfg.Scale.TimeOut <= 0 {
		return fmt.Errorf("Config scale timeOut is invalid")
	}
	if cfg.NodeDrain.TimeOut <= 0 {
		return fmt.Errorf("Config nodeDrain timeOut is invalid")
	}

	return nil
}

func (cfg *osdUpgradeConfig) GetScaleDuration() time.Duration {
	return time.Duration(cfg.Scale.TimeOut) * time.Minute
}

func (cfg *osdUpgradeConfig) GetNodeDrainDuration() time.Duration {
	return time.Duration(cfg.NodeDrain.TimeOut) * time.Minute
}
