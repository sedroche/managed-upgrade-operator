package scaler

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -destination=mocks/scaler.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/scaler Scaler
type Scaler interface {
	EnsureScaleUpNodes(client.Client, logr.Logger) (bool, error)
	EnsureScaleDownNodes(client.Client, logr.Logger) (bool, error)
}

func NewScaler(cfg ScaleConfig) Scaler {
	return &machineSetScaler{cfg}
}

type scaleTimeOutError struct {
	message string
}

func (stoErr *scaleTimeOutError) Error() string {
	return stoErr.message
}

func IsScaleTimeOutError(err error) bool {
	_, ok := err.(*scaleTimeOutError)
	return ok
}

func NewScaleTimeOutError(msg string) *scaleTimeOutError {
	return &scaleTimeOutError{message: msg}
}
