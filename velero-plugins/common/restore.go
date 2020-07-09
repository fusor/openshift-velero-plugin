package common

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
)

// RestorePlugin is a restore item action plugin for Heptio Ark.
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to everything.
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{}, nil
}

// Execute sets a custom annotation on the item being restored.
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[common-restore] Entering common restore plugin")

	metadata, annotations, err := getMetadataAndAnnotations(input.Item)
	if err != nil {
		return nil, err
	}
	name := metadata.GetName()
	p.Log.Infof("[common-restore] common restore plugin for %s", name)

	major, minor, err := GetServerVersion()
	if err != nil {
		return nil, err
	}

	annotations[RestoreServerVersion] = fmt.Sprintf("%v.%v", major, minor)
	registryHostname, err := GetRegistryInfo(major, minor, p.Log)
	if err != nil {
		return nil, err
	}
	annotations[RestoreRegistryHostname] = registryHostname

	if input.Restore.Labels[MigrationApplicationLabelKey] != MigrationApplicationLabelValue {
		// if the current workflow is not CAM(i.e B/R) then get the backup registry route and set the same on annotation to use in plugins.
		backupLocation, err := getBackupStorageLocationForBackup(input.Restore.Spec.BackupName, input.Restore.Namespace)
		if err != nil {
			return nil, err
		}
		tempRegistry, err := getOADPRegistryRoute(input.Restore.Namespace, backupLocation, RegistryConfigMap)
		if err != nil {
			return nil, err
		}
		annotations[MigrationRegistry] = tempRegistry
	} else {
		// if the current workflow is CAM then get migration registry from backup object and set the same on annotation to use in plugins.
		annotations[MigrationRegistry] = input.Restore.Annotations[MigrationRegistry]
	}
	metadata.SetAnnotations(annotations)

	return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
}
