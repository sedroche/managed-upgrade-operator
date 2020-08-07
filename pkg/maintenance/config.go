package maintenance

import (
	"fmt"
	"time"
)

type MaintenanceConfig struct {
	ControlPlaneTime int           `yaml:"controlPlaneTime"`
	WorkerNodeTime   int           `yaml:"workerNodeTime"`
	IgnoredAlerts    ignoredAlerts `yaml:"ignoredAlerts"`
}

type ignoredAlerts struct {
	// Generally upgrades should not fire critical alerts but there are some critical alerts that will fire.
	// e.g. 'etcdMembersDown' happens as the masters drain/reboot and a master is offline but this is expected and will resolve.
	// This is a list of critical alerts that can be ignored while upgrading of controlplane occurs
	ControlPlaneCriticals []string `yaml:"controlPlaneCriticals"`
}

func (cfg *MaintenanceConfig) IsValid() error {
	if cfg.ControlPlaneTime <= 0 {
		return fmt.Errorf("Config maintenace controlPlaneTime out is invalid")
	}
	if cfg.WorkerNodeTime <= 0 {
		return fmt.Errorf("Config maintenace workerNodeTime is invalid")
	}

	return nil
}

func (cfg *MaintenanceConfig) GetControlPlaneDuration() time.Duration {
	return time.Duration(cfg.ControlPlaneTime) * time.Minute
}

func (cfg *MaintenanceConfig) GetWorkerNodeDuration() time.Duration {
	return time.Duration(cfg.WorkerNodeTime) * time.Minute
}
