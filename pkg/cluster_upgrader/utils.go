package cluster_upgrader

import (
	"k8s.io/apimachinery/pkg/types"
	configv1 "github.com/openshift/api/config/v1"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetClusterVersion(c client.Client) (*configv1.ClusterVersion, error){
	cv := &configv1.ClusterVersion{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: "version"}, cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

