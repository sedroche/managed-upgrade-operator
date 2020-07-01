package validation

import (
	configv1 "github.com/openshift/api/config/v1"
	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Validation interface {
	IsValidUpgradeConfig(uc upgradev1alpha1.UpgradeConfig, cv configv1.ClusterVersion) (bool, error)
}

type ValidationBuilder interface {
	NewClient(client client.Client) (Validation, error)
}

func NewBuilder() ValidationBuilder {
	return &validationBuilder{}
}

type validationBuilder struct{}

func (vb *validationBuilder) NewClient(client client.Client) (Validation, error) {
	return nil, nil
}
