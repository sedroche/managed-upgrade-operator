package structs

import (
	"fmt"
	api "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

type testUpgradeConfigBuilder struct {
	uc api.UpgradeConfig
}

func (t *testUpgradeConfigBuilder) GetUpgradeConfig() *api.UpgradeConfig {
	return &t.uc
}

func NewUpgradeConfigBuilder() *testUpgradeConfigBuilder {
	return &testUpgradeConfigBuilder{
		uc: api.UpgradeConfig{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "fakeUpgradeConfig",
				Namespace: "fakeNamespace",
			},
			Spec: api.UpgradeConfigSpec{
				Type: api.OSD,
				Desired: api.Update{
					Version: "fakeVersion",
					Channel: "fakeChannel",
				},
				UpgradeAt:            time.Now().Format(time.RFC3339),
				PDBForceDrainTimeout: 60,
			},
		},
	}
}

func (t *testUpgradeConfigBuilder) WithNamespacedName(namespacedName types.NamespacedName) *testUpgradeConfigBuilder {
	t.uc.ObjectMeta.Name = namespacedName.Name
	t.uc.ObjectMeta.Namespace = namespacedName.Namespace
	return t
}

func (t *testUpgradeConfigBuilder) WithPhase(phase api.UpgradePhase) *testUpgradeConfigBuilder {
	t.uc.Status.History = []api.UpgradeHistory{
		api.UpgradeHistory{
			Version: t.uc.Spec.Desired.Version,
			Phase:   phase,
		},
	}
	return t
}

type UpgradeConfigMatcher struct {
	ActualUpgradeConfig api.UpgradeConfig
	FailReason          string
}

func NewUpgradeConfigMatcher() *UpgradeConfigMatcher {
	return &UpgradeConfigMatcher{}
}

func (m *UpgradeConfigMatcher) Matches(x interface{}) bool {
	ref, isCorrectType := x.(*api.UpgradeConfig)
	if !isCorrectType {
		m.FailReason = fmt.Sprintf("Unexpected type passed: want '%T', got '%T'", api.UpgradeConfig{}, x)
		return false
	}
	m.ActualUpgradeConfig = *ref.DeepCopy()
	return true
}

func (m *UpgradeConfigMatcher) String() string {
	return "Fail reason: " + m.FailReason
}
