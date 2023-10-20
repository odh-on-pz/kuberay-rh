package common

const (

	// Default application name
	DefaultServeAppName = "default"
	// Belows used as label key
	RayServiceLabelKey               = "ray.io/service"
	RayClusterLabelKey               = "ray.io/cluster"
	RayNodeTypeLabelKey              = "ray.io/node-type"
	RayNodeGroupLabelKey             = "ray.io/group"
	RayNodeLabelKey                  = "ray.io/is-ray-node"
	RayIDLabelKey                    = "ray.io/identifier"
	RayClusterServingServiceLabelKey = "ray.io/serve"
	RayServiceClusterHashKey         = "ray.io/cluster-hash"

	// In KubeRay, the Ray container must be the first application container in a head or worker Pod.
	RayContainerIndex = 0

	// Batch scheduling labels
	// TODO(tgaddair): consider making these part of the CRD
	RaySchedulerName     = "ray.io/scheduler-name"
	RayPriorityClassName = "ray.io/priority-class-name"

	// Ray GCS FT related annotations
	RayFTEnabledAnnotationKey         = "ray.io/ft-enabled"
	RayExternalStorageNSAnnotationKey = "ray.io/external-storage-namespace"

	// Finalizers for GCS fault tolerance
	GCSFaultToleranceRedisCleanupFinalizer = "ray.io/gcs-ft-redis-cleanup-finalizer"

	EnableAgentServiceKey  = "ray.io/enableAgentService"
	EnableAgentServiceTrue = "true"

	EnableRayClusterServingServiceTrue  = "true"
	EnableRayClusterServingServiceFalse = "false"

	KubernetesApplicationNameLabelKey = "app.kubernetes.io/name"
	KubernetesCreatedByLabelKey       = "app.kubernetes.io/created-by"

	// Use as separator for pod name, for example, raycluster-small-size-worker-0
	DashSymbol = "-"

	// Use as default port
	DefaultClientPort = 10001
	// For Ray >= 1.11.0, "DefaultRedisPort" actually refers to the GCS server port.
	// However, the role of this port is unchanged in Ray APIs like ray.init and ray start.
	// This is the port used by Ray workers and drivers inside the Ray cluster to connect to the Ray head.
	DefaultRedisPort                = 6379
	DefaultDashboardPort            = 8265
	DefaultMetricsPort              = 8080
	DefaultDashboardAgentListenPort = 52365
	DefaultServingPort              = 8000

	ClientPortName               = "client"
	RedisPortName                = "redis"
	DashboardPortName            = "dashboard"
	MetricsPortName              = "metrics"
	DashboardAgentListenPortName = "dashboard-agent"
	ServingPortName              = "serve"

	// The default AppProtocol for Kubernetes service
	DefaultServiceAppProtocol = "tcp"

	// The default application name
	ApplicationName = "kuberay"

	// The default name for kuberay operator
	ComponentName = "kuberay-operator"

	// The defaule RayService Identifier.
	RayServiceCreatorLabelValue = "rayservice"

	// Use as container env variable
	RAY_CLUSTER_NAME                        = "RAY_CLUSTER_NAME"
	RAY_IP                                  = "RAY_IP"
	FQ_RAY_IP                               = "FQ_RAY_IP"
	RAY_PORT                                = "RAY_PORT"
	RAY_ADDRESS                             = "RAY_ADDRESS"
	REDIS_PASSWORD                          = "REDIS_PASSWORD"
	RAY_DASHBOARD_ENABLE_K8S_DISK_USAGE     = "RAY_DASHBOARD_ENABLE_K8S_DISK_USAGE"
	RAY_EXTERNAL_STORAGE_NS                 = "RAY_external_storage_namespace"
	RAY_GCS_RPC_SERVER_RECONNECT_TIMEOUT_S  = "RAY_gcs_rpc_server_reconnect_timeout_s"
	RAY_TIMEOUT_MS_TASK_WAIT_FOR_DEATH_INFO = "RAY_timeout_ms_task_wait_for_death_info"
	RAY_GCS_SERVER_REQUEST_TIMEOUT_SECONDS  = "RAY_gcs_server_request_timeout_seconds"
	RAY_SERVE_KV_TIMEOUT_S                  = "RAY_SERVE_KV_TIMEOUT_S"
	SERVE_CONTROLLER_PIN_ON_NODE            = "RAY_INTERNAL_SERVE_CONTROLLER_PIN_ON_NODE"
	RAY_USAGE_STATS_KUBERAY_IN_USE          = "RAY_USAGE_STATS_KUBERAY_IN_USE"
	RAYCLUSTER_DEFAULT_REQUEUE_SECONDS_ENV  = "RAYCLUSTER_DEFAULT_REQUEUE_SECONDS_ENV"
	RAYCLUSTER_DEFAULT_REQUEUE_SECONDS      = 300

	// This KubeRay operator environment variable is used to determine if random Pod
	// deletion should be enabled. Note that this only takes effect when autoscaling
	// is enabled for the RayCluster. This is a feature flag for v0.6.0, and will be
	// removed if the default behavior is stable enoguh.
	ENABLE_RANDOM_POD_DELETE = "ENABLE_RANDOM_POD_DELETE"

	// This KubeRay operator environment variable is used to determine if the Redis
	// cleanup Job should be enabled. This is a feature flag for v1.0.0.
	ENABLE_GCS_FT_REDIS_CLEANUP = "ENABLE_GCS_FT_REDIS_CLEANUP"

	// Ray core default configurations
	DefaultWorkerRayGcsReconnectTimeoutS = "600"

	LOCAL_HOST = "127.0.0.1"
	// Ray FT default readiness probe values
	DefaultReadinessProbeInitialDelaySeconds = 10
	DefaultReadinessProbeTimeoutSeconds      = 1
	DefaultReadinessProbePeriodSeconds       = 5
	DefaultReadinessProbeSuccessThreshold    = 1
	DefaultReadinessProbeFailureThreshold    = 10

	// Ray FT default liveness probe values
	DefaultLivenessProbeInitialDelaySeconds = 30
	DefaultLivenessProbeTimeoutSeconds      = 1
	DefaultLivenessProbePeriodSeconds       = 5
	DefaultLivenessProbeSuccessThreshold    = 1
	DefaultLivenessProbeFailureThreshold    = 120

	// Ray health check related configurations
	RayAgentRayletHealthPath  = "api/local_raylet_healthz"
	RayDashboardGCSHealthPath = "api/gcs_healthz"

	// Finalizers for RayJob
	RayJobStopJobFinalizer = "ray.io/rayjob-finalizer"
)

type ServiceType string

const (
	HeadService    ServiceType = "headService"
	ServingService ServiceType = "serveService"
)
