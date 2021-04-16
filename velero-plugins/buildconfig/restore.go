package buildconfig

import (
	"context"
	"encoding/json"

	"github.com/konveyor/openshift-velero-plugin/velero-plugins/build"
	"github.com/konveyor/openshift-velero-plugin/velero-plugins/clients"
	"github.com/konveyor/openshift-velero-plugin/velero-plugins/common"
	buildv1API "github.com/openshift/api/build/v1"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to buildconfigs
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"buildconfigs"},
	}, nil
}

// Execute action for the restore plugin for the buildconfig resource
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[buildconfig-restore] Entering buildconfig restore plugin")

	buildconfig := buildv1API.BuildConfig{}
	itemMarshal, _ := json.Marshal(input.Item)
	json.Unmarshal(itemMarshal, &buildconfig)

	buildconfigUnmodified := buildv1API.BuildConfig{}
	itemMarshal, _ = json.Marshal(input.ItemFromBackup)
	json.Unmarshal(itemMarshal, &buildconfigUnmodified)

	buildconfig, err := p.updateSecretsAndDockerRefs(buildconfig, buildconfigUnmodified.Namespace, input.Restore.Spec.NamespaceMapping)
	if err != nil {
		p.Log.Error("[buildconfig-restore] error modifying buildconfig: ", err)
		return nil, err
	}

	var out map[string]interface{}
	objrec, _ := json.Marshal(buildconfig)
	json.Unmarshal(objrec, &out)

	return velero.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: out}), nil
}

func (p *RestorePlugin) updateSecretsAndDockerRefs(buildconfig buildv1API.BuildConfig, srcNamespace string, namespaceMapping map[string]string) (buildv1API.BuildConfig, error) {
	client, err := clients.CoreClient()
	if err != nil {
		return buildconfig, err
	}

	destNamespace := srcNamespace
	if namespaceMapping[destNamespace] != "" {
		destNamespace = namespaceMapping[destNamespace]
	}
	secretList, err := client.Secrets(destNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return buildconfig, err
	}

	registry := buildconfig.Annotations[common.RestoreRegistryHostname]
	backupRegistry := buildconfig.Annotations[common.BackupRegistryHostname]

	newCommonSpec, err := build.UpdateCommonSpec(buildconfig.Spec.CommonSpec, registry, backupRegistry, secretList, p.Log, namespaceMapping)
	if err != nil {
		return buildconfig, err
	}
	buildconfig.Spec.CommonSpec = newCommonSpec
	return buildconfig, nil
}
