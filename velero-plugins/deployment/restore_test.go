// unit tests for restore.go in deployment

package deployment

import (
	"testing"
	"github.com/konveyor/openshift-velero-plugin/velero-plugins/util/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
        appsv1API "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"encoding/json"
	"reflect"
        velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	//"fmt"
	//"k8s.io/apimachinery/pkg/runtime"
)

func TestRestorePluginAppliesTo(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}
	actual, err := restorePlugin.AppliesTo()
	require.NoError(t, err)
	assert.Equal(t, velero.ResourceSelector{IncludedResources: []string{"deployments.apps"}}, actual)
}

func TestRestorePluginExecute(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}

	tests := map[string]struct{
		deployment appsv1API.Deployment
		exp	   appsv1API.Deployment
	}{
		"1": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "backup-host/namespace-old/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
                                                },
                                        },
                                },
                        },
		},

		"2": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "backup-host/namespace-old/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
                                                },
					},
				},
			},
		},

		"3": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace/foo"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace/foo"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "restore-host/namespace/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "restore-host/namespace/foo"},
							},
                                                },
					},
				},
			},
		},

		"4": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/namespace-old/foo"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "restore-host/namespace-new/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "restore-host/namespace-new/foo"},
							},
                                                },
                                        },
                                },
                        },
		},

		"5": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host-2/namespace-old/foo"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host-2/namespace-old/foo"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "backup-host-2/namespace-old/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host-2/namespace-old/foo"},
							},
                                                },
                                        },
                                },
                        },
		},

		"6": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/openshift/foo@bar"},
							},
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "backup-host/openshift/foo@bar"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "restore-host/openshift/foo"},
                                                        },
							InitContainers: []apiv1.Container {
								apiv1.Container{Image: "restore-host/openshift/foo"},
							},
                                                },
                                        },
                                },
                        },
		},
	}

	for name, tc := range tests {
                t.Run(name, func(t *testing.T) {
			var out map[string]interface{}
			item := unstructured.Unstructured{}
			deploymentRec, _ := json.Marshal(tc.deployment) // Marshal it to JSON
			json.Unmarshal(deploymentRec, &out) // Unmarshal into the proper format
			item.SetUnstructuredContent(out) // Set unstructured object
			restore := velerov1.Restore{
				Spec: velerov1.RestoreSpec{
					NamespaceMapping: map[string]string{
						"namespace-old": "namespace-new",
					},
				},
			}
			input := &velero.RestoreItemActionExecuteInput{Item: &item, Restore: &restore}

			output, _ := restorePlugin.Execute(input)

			deploymentOut := appsv1API.Deployment{}
			itemMarshal, _ := json.Marshal(output.UpdatedItem)
			json.Unmarshal(itemMarshal, &deploymentOut)

			if !reflect.DeepEqual(deploymentOut, tc.exp) {
                                t.Fatalf("expected: %v, got: %v", tc.exp, deploymentOut)
                        }
		})
        }
}

