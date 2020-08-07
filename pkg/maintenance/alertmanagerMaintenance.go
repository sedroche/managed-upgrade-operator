package maintenance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-multierror"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/managed-upgrade-operator/config"
	amSilence "github.com/prometheus/alertmanager/api/v2/client/silence"
	amv2Models "github.com/prometheus/alertmanager/api/v2/models"
)

var (
	alertManagerNamespace          = "openshift-monitoring"
	alertManagerRouteName          = "alertmanager-main"
	alertManagerServiceAccountName = "prometheus-k8s"
	alertManagerBasePath           = "/api/v2/"
)

type alertManagerMaintenanceBuilder struct{}

func (ammb *alertManagerMaintenanceBuilder) NewClient(client client.Client, mcfg MaintenanceConfig) (Maintenance, error) {
	transport, err := getTransport(client)
	if err != nil {
		return nil, err
	}

	transport.DefaultAuthentication, err = getAuthentication(client)
	if err != nil {
		return nil, err
	}

	return &alertManagerMaintenance{
		client: &alertManagerSilenceClient{
			transport: transport,
		},
		cfg: mcfg,
	}, nil
}

type alertManagerMaintenance struct {
	//	client alertManagerSilenceClient
	client AlertManagerSilencer
	cfg    MaintenanceConfig
}

func getTransport(c client.Client) (*httptransport.Runtime, error) {
	amRoute := &routev1.Route{}
	err := c.Get(
		context.TODO(),
		types.NamespacedName{Namespace: alertManagerNamespace, Name: alertManagerRouteName},
		amRoute,
	)
	if err != nil {
		return nil, err
	}

	return httptransport.New(
		amRoute.Spec.Host,
		alertManagerBasePath,
		[]string{"https"},
	), nil
}

func getAuthentication(c client.Client) (runtime.ClientAuthInfoWriter, error) {
	sl := &corev1.SecretList{}
	err := c.List(
		context.TODO(),
		sl,
		&client.ListOptions{Namespace: alertManagerNamespace},
	)
	if err != nil {
		return nil, err
	}

	var token string
	for _, s := range sl.Items {
		if strings.Contains(s.Name, alertManagerServiceAccountName+"-token") {
			token = string(s.Data["token"])
		}
	}

	return httptransport.BearerToken(token), nil
}

// Start a control plane maintenance in Alertmanager for version
func (amm *alertManagerMaintenance) StartControlPlane(version string) error {
	defaultComment := fmt.Sprintf("Silence for OSD control plane upgrade to version %s", version)
	criticalAlertComment := fmt.Sprintf("Silence for critical alerts during OSD control plane upgrade to version %s", version)
	mList, err := amm.client.List([]string{})
	if err != nil {
		return err
	}
	defaultExists := silenceExists(mList, defaultComment)
	criticalExists := silenceExists(mList, criticalAlertComment)

	if defaultExists && criticalExists {
		return nil
	}

	endsAt := time.Now().Add(amm.cfg.GetControlPlaneDuration())
	now := strfmt.DateTime(time.Now().UTC())
	end := strfmt.DateTime(endsAt.UTC())
	if !defaultExists {
		err = amm.client.Create(createDefaultMatchers(), now, end, config.OperatorName, defaultComment)
		if err != nil {
			return err
		}
	}

	if !criticalExists {
		ic := amm.cfg.IgnoredAlerts.ControlPlaneCriticals
		if len(ic) > 0 {
			icRegex := "(" + strings.Join(ic, "|") + ")"
			matchers := []*amv2Models.Matcher{createMatcher("alertname", icRegex, true)}
			err = amm.client.Create(matchers, now, end, config.OperatorName, criticalAlertComment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Start a worker node maintenance in Alertmanager for version
func (amm *alertManagerMaintenance) StartWorker(numWorkers int, version string) error {
	comment := fmt.Sprintf("Silence for OSD worker node upgrade to version %s", version)
	mList, err := amm.client.List([]string{})
	if err != nil {
		return err
	}
	exists := silenceExists(mList, comment)
	if exists {
		return nil
	}

	workerMaintenanceExpectedDuration := time.Duration(numWorkers) * amm.cfg.GetWorkerNodeDuration()
	endsAt := time.Now().Add(workerMaintenanceExpectedDuration)
	now := strfmt.DateTime(time.Now().UTC())
	end := strfmt.DateTime(endsAt.UTC())
	err = amm.client.Create(createDefaultMatchers(), now, end, config.OperatorName, comment)
	if err != nil {
		return err
	}

	return nil
}

// End all active maintenances created by managed-upgrade-operator in Alertmanager
func (amm *alertManagerMaintenance) End() error {
	silences, err := amm.client.List([]string{})
	if err != nil {
		return err
	}

	var deleteErrors *multierror.Error
	for _, s := range silences.Payload {
		if *s.CreatedBy == config.OperatorName && *s.Status.State == amv2Models.SilenceStatusStateActive {
			err := amm.client.Delete(*s.ID)
			if err != nil {
				deleteErrors = multierror.Append(deleteErrors, err)
			}
		}
	}
	return deleteErrors.ErrorOrNil()
}

func createMatcher(alertMatchKey string, alertValue string, isRegex bool) *amv2Models.Matcher {
	return &amv2Models.Matcher{
		Name:    &alertMatchKey,
		IsRegex: &isRegex,
		Value:   &alertValue,
	}
}

func createDefaultMatchers() []*amv2Models.Matcher {
	// Upgrades can impact some availability which may trigger info/warning alerts. ignore those.
	nonCriticalAlertMatcher := createMatcher("severity", "(warning|info|none)", true)

	inNamespaceAlertMatcher := createMatcher("namespace", "(^openshift.*|^kube.*|^redhat.*|^default$)", true)
	return amv2Models.Matchers{nonCriticalAlertMatcher, inNamespaceAlertMatcher}
}

// Check if a silence with comment exists
func silenceExists(mList *amSilence.GetSilencesOK, comment string) bool {
	for _, m := range mList.Payload {
		if *m.Comment == comment {
			return true
		}
	}

	return false
}

func (amm *alertManagerMaintenance) IsActive() (bool, error) {
	silences, err := amm.client.List([]string{})
	if err != nil {
		return false, err
	}

	for _, s := range silences.Payload {
		if *s.CreatedBy == config.OperatorName && *s.Status.State == amv2Models.AlertStatusStateActive {
			return true, nil
		}
	}
	return false, nil
}
