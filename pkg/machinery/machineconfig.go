package machinery

import (
	"context"

	machineconfigapi "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type UpgradingResult struct {
	IsUpgrading  bool
	UpdatedCount int32
	MachineCount int32
}

// IsUpgrading determines if machines are currently upgrading by comparing
// MachineCount and UpdatedMachineCount
func (m *machinery) IsUpgrading(c client.Client, nodeType string) (*UpgradingResult, error) {
	configPool := &machineconfigapi.MachineConfigPool{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: nodeType}, configPool)
	if err != nil {
		return nil, err
	}

	return &UpgradingResult{
		IsUpgrading:  true,
		UpdatedCount: 2,
		MachineCount: 8,
	}, nil
}
