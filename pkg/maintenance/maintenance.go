package maintenance

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -destination=mocks/maintenance.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/maintenance Maintenance
type Maintenance interface {
	StartControlPlane(version string) error
	StartWorker(numWorkers int, version string) error
	End() error
	IsActive() (bool, error)
}

//go:generate mockgen -destination=mocks/maintenanceBuilder.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/maintenance MaintenanceBuilder
type MaintenanceBuilder interface {
	NewClient(client client.Client, mcfg MaintenanceConfig) (Maintenance, error)
}

func NewBuilder() MaintenanceBuilder {
	return &alertManagerMaintenanceBuilder{}
}
