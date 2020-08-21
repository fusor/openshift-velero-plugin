// unit tests for restore.go in cronjob

package cronjob

import (
	"testing"
	"github.com/konveyor/openshift-velero-plugin/velero-plugins/util/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
        //appsv1API "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"encoding/json"
	"reflect"
        velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	//"fmt"
	//"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/batch/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
)

func TestRestorePluginAppliesTo(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}
	actual, err := restorePlugin.AppliesTo()
	require.NoError(t, err)
	assert.Equal(t, velero.ResourceSelector{IncludedResources: []string{"cronjobs"}}, actual)
}

func TestRestorePluginExecute(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}

	tests := map[string]struct{
		cronJob v1beta1.CronJob
		exp	   v1beta1.CronJob
	}{
		"1": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
                        },
		},

		"2": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
			},
		},

		"3": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
				},
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
			},
		},

		"4": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
                        },
		},

		"5": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
                        },
		},

		"6": {
			cronJob: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
				Spec: v1beta1.CronJobSpec {
					JobTemplate: v1beta1.JobTemplateSpec {
						Spec: batchv1.JobSpec {
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
				},
			},
			exp: v1beta1.CronJob {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "backup-host",
						"openshift.io/restore-registry-hostname": "restore-host",
					},
				},
                                Spec: v1beta1.CronJobSpec {
                                        JobTemplate: v1beta1.JobTemplateSpec {
                                                Spec: batchv1.JobSpec {
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
                        },
		},

	}

	for name, tc := range tests {
                t.Run(name, func(t *testing.T) {
			var out map[string]interface{}
			item := unstructured.Unstructured{}
			cronJobRec, _ := json.Marshal(tc.cronJob) // Marshal it to JSON
			json.Unmarshal(cronJobRec, &out) // Unmarshal into the proper format
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

			cronJobOut := v1beta1.CronJob{}
			itemMarshal, _ := json.Marshal(output.UpdatedItem)
			json.Unmarshal(itemMarshal, &cronJobOut)

			if !reflect.DeepEqual(cronJobOut, tc.exp) {
                                t.Fatalf("expected: %v, got: %v", tc.exp, cronJobOut)
                        }
		})
        }
}

