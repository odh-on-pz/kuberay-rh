package model

import (
	"fmt"
	"reflect"
	"testing"

	util "github.com/ray-project/kuberay/apiserver/pkg/util"
	api "github.com/ray-project/kuberay/proto/go_client"
	rayv1api "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	enableIngress                    = true
	deploymentReplicas       int32   = 1
	headNodeReplicas         int32   = 1
	workerReplicas           int32   = 5
	unhealthySecondThreshold int32   = 900
	floatNumber              float64 = 1
	secondsValue             int32   = 100
)

var headSpecTest = rayv1api.HeadGroupSpec{
	ServiceType:   "ClusterIP",
	EnableIngress: &enableIngress,
	Replicas:      &headNodeReplicas,
	RayStartParams: map[string]string{
		"dashboard-host":      "0.0.0.0",
		"metrics-export-port": "8080",
		"num-cpus":            "0",
	},
	Template: v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"openshift.io/scc":    "restricted",
				"ray.io/ft-enabled":   "false",
				"ray.io/health-state": "",
				"custom":              "value",
			},
			Labels: map[string]string{
				"app.kubernetes.io/created-by": "kuberay-operator",
				"app.kubernetes.io/name":       "kuberay",
				"ray.io/cluster":               "boris-cluster",
				"ray.io/cluster-dashboard":     "boris-cluster-dashboard",
				"ray.io/group":                 "headgroup",
				"ray.io/identifier":            "boris-cluster-head",
				"ray.io/is-ray-node":           "yes",
				"ray.io/node-type":             "head",
				"test":                         "value",
			},
			Name:      "boris-cluster-head-f7zx2",
			Namespace: "max",
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "account",
			ImagePullSecrets: []v1.LocalObjectReference{
				{Name: "foo"},
			},
			Tolerations: []v1.Toleration{
				{
					Key:      "blah1",
					Operator: "Exists",
					Effect:   "NoExecute",
				},
			},
			Containers: []v1.Container{
				{
					Name:  "ray-head",
					Image: "blublinsky1/ray310:2.5.0",
					Env: []v1.EnvVar{
						{
							Name:  "AWS_KEY",
							Value: "123",
						},
						{
							Name: "REDIS_PASSWORD",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis-password-secret",
									},
									Key: "password",
								},
							},
						},
						{
							Name: "CONFIGMAP",
							ValueFrom: &v1.EnvVarSource{
								ConfigMapKeyRef: &v1.ConfigMapKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "special-config",
									},
									Key: "special.how",
								},
							},
						},
						{
							Name: "ResourceFieldRef",
							ValueFrom: &v1.EnvVarSource{
								ResourceFieldRef: &v1.ResourceFieldSelector{
									ContainerName: "my-container",
									Resource:      "resource",
								},
							},
						},
						{
							Name: "FieldRef",
							ValueFrom: &v1.EnvVarSource{
								FieldRef: &v1.ObjectFieldSelector{
									FieldPath: "path",
								},
							},
						},
					},
				},
			},
		},
	},
}

var configMapWithoutTolerations = v1.ConfigMap{
	Data: map[string]string{
		"cpu":             "4",
		"gpu":             "0",
		"gpu_accelerator": "",
		"memory":          "8",
		"name":            "head-node-template",
		"namespace":       "max",
	},
}

var configMapWithTolerations = v1.ConfigMap{
	Data: map[string]string{
		"cpu":             "4",
		"gpu":             "0",
		"gpu_accelerator": "",
		"memory":          "8",
		"name":            "head-node-template",
		"namespace":       "max",
		"tolerations":     "[{\"key\":\"blah1\",\"operator\":\"Exists\",\"effect\":\"NoExecute\"}]",
	},
}

var workerSpecTest = rayv1api.WorkerGroupSpec{
	GroupName:   "",
	Replicas:    &workerReplicas,
	MinReplicas: &workerReplicas,
	MaxReplicas: &workerReplicas,
	RayStartParams: map[string]string{
		"node-ip-address": "$MY_POD_IP",
	},
	Template: v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"cni.projectcalico.org/containerID": "cce862a899455385e98e3453ba9ef5a376e85ad45c3e95b18e04e001204af728",
				"cni.projectcalico.org/podIP":       "172.17.60.2/32",
				"cni.projectcalico.org/podIPs":      "172.17.60.2/32",
				"openshift.io/scc":                  "restricted",
				"ray.io/ft-enabled":                 "false",
				"ray.io/health-state":               "",
				"custom":                            "value",
			},
			Labels: map[string]string{
				"app.kubernetes.io/created-by": "kuberay-operator",
				"app.kubernetes.io/name":       "kuberay",
				"ray.io/cluster":               "boris-cluster",
				"ray.io/cluster-dashboard":     "boris-cluster-dashboard",
				"ray.io/group":                 "8-CPUs",
				"ray.io/identifier":            "boris-cluster-worker",
				"ray.io/is-ray-node":           "yes",
				"ray.io/node-type":             "worker",
				"test":                         "value",
			},
			Name:      "boris-cluster-worker-8-cpus-4dp9v",
			Namespace: "max",
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "account",
			ImagePullSecrets: []v1.LocalObjectReference{
				{Name: "foo"},
			},
			Tolerations: []v1.Toleration{
				{
					Key:      "blah1",
					Operator: "Exists",
					Effect:   "NoExecute",
				},
			},
			Containers: []v1.Container{
				{
					Name:  "ray-worker",
					Image: "blublinsky1/ray310:2.5.0",
					Env: []v1.EnvVar{
						{
							Name:  "AWS_KEY",
							Value: "123",
						},
						{
							Name:  "AWS_SECRET",
							Value: "1234",
						},
						{
							Name:  "RAY_DISABLE_DOCKER_CPU_WARNING",
							Value: "1",
						},
						{
							Name:  "TYPE",
							Value: "worker",
						},
						{
							Name:  "RAY_IP",
							Value: "boris-cluster-head-svc",
						},
						{
							Name:  "RAY_USAGE_STATS_KUBERAY_IN_USE",
							Value: "1",
						},
					},
				},
			},
		},
	},
}

var ClusterSpecTest = rayv1api.RayCluster{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "raycluster-sample",
		Namespace: "default",
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
		},
	},
	Spec: rayv1api.RayClusterSpec{
		HeadGroupSpec: headSpecTest,
		WorkerGroupSpecs: []rayv1api.WorkerGroupSpec{
			workerSpecTest,
		},
	},
}

var JobNewClusterTest = rayv1api.RayJob{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: "test",
		Labels: map[string]string{
			"ray.io/user": "user",
		},
	},
	Spec: rayv1api.RayJobSpec{
		Entrypoint: "python /home/ray/samples/sample_code.py",
		Metadata: map[string]string{
			"job_submission_id": "123",
		},
		RuntimeEnvYAML:          "mytest yaml",
		TTLSecondsAfterFinished: &secondsValue,
		RayClusterSpec:          &ClusterSpecTest.Spec,
	},
}

var JobExistingClusterTest = rayv1api.RayJob{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: "test",
		Labels: map[string]string{
			"ray.io/user": "user",
		},
	},
	Spec: rayv1api.RayJobSpec{
		Entrypoint:              "python /home/ray/samples/sample_code.py",
		RuntimeEnvYAML:          "mytest yaml",
		TTLSecondsAfterFinished: &secondsValue,
		ClusterSelector: map[string]string{
			util.RayClusterUserLabelKey: "test",
		},
	},
}

var JobExistingClusterSubmitterTest = rayv1api.RayJob{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: "test",
		Labels: map[string]string{
			"ray.io/user": "user",
		},
	},
	Spec: rayv1api.RayJobSpec{
		Entrypoint:              "python /home/ray/samples/sample_code.py",
		RuntimeEnvYAML:          "mytest yaml",
		TTLSecondsAfterFinished: &secondsValue,
		ClusterSelector: map[string]string{
			util.RayClusterUserLabelKey: "test",
		},
		SubmitterPodTemplate: &v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "test-submitter",
						Image: "image",
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("2"),
								v1.ResourceMemory: resource.MustParse("1Gi"),
							},
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("500m"),
								v1.ResourceMemory: resource.MustParse("200Mi"),
							},
						},
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		},
	},
}

var ServiceV1Test = rayv1api.RayService{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: "test",
		Labels: map[string]string{
			"ray.io/user": "user",
		},
	},
	Spec: rayv1api.RayServiceSpec{
		ServeDeploymentGraphSpec: rayv1api.ServeDeploymentGraphSpec{
			ImportPath: "fruit.deployment_graph",
			RuntimeEnv: "working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"",
			ServeConfigSpecs: []rayv1api.ServeConfigSpec{
				{
					Name:        "MangoStand",
					NumReplicas: &deploymentReplicas,
					UserConfig:  "price: 3",
					RayActorOptions: rayv1api.RayActorOptionSpec{
						NumCpus: &floatNumber,
					},
				},
				{
					Name:        "OrangeStand",
					NumReplicas: &deploymentReplicas,
				},
				{
					Name:        "PearStand",
					NumReplicas: &deploymentReplicas,
					UserConfig:  "price: 1",
					RayActorOptions: rayv1api.RayActorOptionSpec{
						NumCpus: &floatNumber,
					},
				},
			},
		},
		RayClusterSpec:                  ClusterSpecTest.Spec,
		ServiceUnhealthySecondThreshold: &unhealthySecondThreshold,
	},
}

var ServiceV2Test = rayv1api.RayService{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: "test",
		Labels: map[string]string{
			"ray.io/user": "user",
		},
	},
	Spec: rayv1api.RayServiceSpec{
		ServeConfigV2:                      "Some yaml value",
		RayClusterSpec:                     ClusterSpecTest.Spec,
		DeploymentUnhealthySecondThreshold: &unhealthySecondThreshold,
	},
}

var expectedAnnotations = map[string]string{
	"custom": "value",
}

var expectedLabels = map[string]string{
	"test": "value",
}

var expectedHeadEnv = &api.EnvironmentVariables{
	Values: map[string]string{
		"AWS_KEY": "123",
	},
	ValuesFrom: map[string]*api.EnvValueFrom{
		"REDIS_PASSWORD": {
			Source: api.EnvValueFrom_SECRET,
			Name:   "redis-password-secret",
			Key:    "password",
		},
		"CONFIGMAP": {
			Source: api.EnvValueFrom_CONFIGMAP,
			Name:   "special-config",
			Key:    "special.how",
		},
		"ResourceFieldRef": {
			Source: api.EnvValueFrom_RESOURCEFIELD,
			Name:   "my-container",
			Key:    "resource",
		},
		"FieldRef": {
			Source: api.EnvValueFrom_FIELD,
			Key:    "path",
		},
	},
}

var expectedEnv = &api.EnvironmentVariables{
	Values: map[string]string{
		"AWS_KEY":    "123",
		"AWS_SECRET": "1234",
	},
	ValuesFrom: map[string]*api.EnvValueFrom{},
}

var expectedTolerations = api.PodToleration{
	Key:      "blah1",
	Operator: "Exists",
	Effect:   "NoExecute",
}

func TestPopulateHeadNodeSpec(t *testing.T) {
	groupSpec := PopulateHeadNodeSpec(headSpecTest)

	if groupSpec.ServiceAccount != "account" {
		t.Errorf("failed to convert service account")
	}
	if groupSpec.EnableIngress != *headSpecTest.EnableIngress {
		t.Errorf("failed to convert enableIngress")
	}
	if groupSpec.ImagePullSecret != "foo" {
		t.Errorf("failed to convert image pull secret")
	}
	if !reflect.DeepEqual(groupSpec.Annotations, expectedAnnotations) {
		t.Errorf("failed to convert annotations, got %v, expected %v", groupSpec.Annotations, expectedAnnotations)
	}
	if !reflect.DeepEqual(groupSpec.Labels, expectedLabels) {
		t.Errorf("failed to convert labels, got %v, expected %v", groupSpec.Labels, expectedLabels)
	}
	if !reflect.DeepEqual(groupSpec.Environment, expectedHeadEnv) {
		t.Errorf("failed to convert environment, got %v, expected %v", groupSpec.Environment, expectedHeadEnv)
	}
}

func TestPopulateWorkerNodeSpec(t *testing.T) {
	groupSpec := PopulateWorkerNodeSpec([]rayv1api.WorkerGroupSpec{workerSpecTest})[0]

	if groupSpec.ServiceAccount != "account" {
		t.Errorf("failed to convert service account")
	}
	if groupSpec.ImagePullSecret != "foo" {
		t.Errorf("failed to convert image pull secret")
	}
	if !reflect.DeepEqual(groupSpec.Annotations, expectedAnnotations) {
		t.Errorf("failed to convert annotations, got %v, expected %v", groupSpec.Annotations, expectedAnnotations)
	}
	if !reflect.DeepEqual(groupSpec.Labels, expectedLabels) {
		t.Errorf("failed to convert labels, got %v, expected %v", groupSpec.Labels, expectedLabels)
	}
	if !reflect.DeepEqual(groupSpec.Environment, expectedEnv) {
		t.Errorf("failed to convert environment, got %v, expected %v", groupSpec.Environment, expectedEnv)
	}
}

func TestPopulateRayClusterSpec(t *testing.T) {
	cluster := FromCrdToApiCluster(&ClusterSpecTest, []v1.Event{})
	if len(cluster.Annotations) != 1 {
		t.Errorf("failed to convert cluster's annotations")
	}
}

func TestPopulateTemplate(t *testing.T) {
	template := FromKubeToAPIComputeTemplate(&configMapWithoutTolerations)
	if len(template.Tolerations) != 0 {
		t.Errorf("failed to convert config map, expected no tolerations, got %d", len(template.Tolerations))
	}

	template = FromKubeToAPIComputeTemplate(&configMapWithTolerations)
	if len(template.Tolerations) != 1 {
		t.Errorf("failed to convert config map, expected 1 toleration, got %d", len(template.Tolerations))
	}
	if template.Tolerations[0].Key != expectedTolerations.Key ||
		template.Tolerations[0].Operator != expectedTolerations.Operator ||
		template.Tolerations[0].Effect != expectedTolerations.Effect {
		t.Errorf("failed to convert config map, got %v, expected %v", tolerationToString(template.Tolerations[0]),
			tolerationToString(&expectedTolerations))
	}
}

func tolerationToString(toleration *api.PodToleration) string {
	return "Key: " + toleration.Key + " Operator: " + string(toleration.Operator) + " Effect: " + string(toleration.Effect)
}

func TestPopulateJob(t *testing.T) {
	job := FromCrdToApiJob(&JobNewClusterTest)
	fmt.Printf("jobWithCluster = %#v\n", job)
	assert.Equal(t, "test", job.Name)
	assert.Equal(t, "test", job.Namespace)
	assert.Equal(t, "user", job.User)
	assert.Greater(t, len(job.RuntimeEnv), 1)
	assert.Nil(t, job.ClusterSelector)
	assert.NotNil(t, job.ClusterSpec)

	job = FromCrdToApiJob(&JobExistingClusterTest)
	fmt.Printf("jobReferenceCluster = %#v\n", job)
	assert.Equal(t, "test", job.Name)
	assert.Equal(t, "test", job.Namespace)
	assert.Equal(t, "user", job.User)
	assert.Greater(t, len(job.RuntimeEnv), 1)
	assert.NotNil(t, job.ClusterSelector)
	assert.Nil(t, job.ClusterSpec)

	job = FromCrdToApiJob(&JobExistingClusterSubmitterTest)
	fmt.Printf("jobReferenceCluster = %#v\n", job)
	assert.Equal(t, "test", job.Name)
	assert.Equal(t, "test", job.Namespace)
	assert.Equal(t, "user", job.User)
	assert.Greater(t, len(job.RuntimeEnv), 1)
	assert.NotNil(t, job.ClusterSelector)
	assert.Nil(t, job.ClusterSpec)
	assert.Equal(t, "image", job.JobSubmitter.Image)
	assert.Equal(t, "2", job.JobSubmitter.Cpu)
}

func TestPopulateService(t *testing.T) {
	service := FromCrdToApiService(&ServiceV1Test, make([]v1.Event, 0))
	fmt.Printf("serviceV1 = %#v\n", service)
	if service.Name != "test" {
		t.Errorf("failed to convert name")
	}
	if service.Namespace != "test" {
		t.Errorf("failed to convert namespace")
	}
	if service.User != "user" {
		t.Errorf("failed to convert user")
	}
	if service.ServeDeploymentGraphSpec == nil {
		t.Errorf("failed to convert v1 serve spec")
	}
	if service.ServeConfig_V2 != "" {
		t.Errorf("unexpected v2 server spec")
	}
	if len(service.ServeDeploymentGraphSpec.ServeConfigs) != 3 {
		t.Errorf("failed to convert serveConfiggs")
	}
	service = FromCrdToApiService(&ServiceV2Test, make([]v1.Event, 0))
	fmt.Printf("serviceV2 = %#v\n", service)
	if service.ServeDeploymentGraphSpec != nil {
		t.Errorf("unexpected v1 serve spec")
	}
	if service.ServeConfig_V2 == "" {
		t.Errorf("failed to convert v2 server spec")
	}
}
