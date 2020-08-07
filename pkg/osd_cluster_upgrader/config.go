package osd_cluster_upgrader

import (
	"fmt"
	"time"

	"github.com/openshift/managed-upgrade-operator/pkg/maintenance"
	"github.com/openshift/managed-upgrade-operator/pkg/scaler"
)

type osdUpgradeConfig struct {
	Maintenance maintenance.MaintenanceConfig `yaml:"maintenance"`
	Scale       scaler.ScaleConfig            `yaml:"scale"`
	NodeDrain   nodeDrain                     `yaml:"nodeDrain"`
	HealthCheck healthCheck                   `yaml:"healthCheck"`
}

type healthCheck struct {
	IgnoredCriticals []string `yaml:"ignoredCriticals"`
}

type nodeDrain struct {
	TimeOut int `yaml:"timeOut"`
}

func (cfg *osdUpgradeConfig) IsValid() error {
	if err := cfg.Maintenance.IsValid(); err != nil {
		return err
	}
	if err := cfg.Scale.IsValid(); err != nil {
		return err
	}
	if cfg.NodeDrain.TimeOut <= 0 {
		return fmt.Errorf("Config nodeDrain timeOut is invalid")
	}

	return nil
}

func (cfg *osdUpgradeConfig) GetNodeDrainDuration() time.Duration {
	return time.Duration(cfg.NodeDrain.TimeOut) * time.Minute
}
