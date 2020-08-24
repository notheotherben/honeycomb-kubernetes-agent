package metrics

const (
	LabelNodeName      = PrefixNode + "name"
	LabelNamespaceName = PrefixNamespace + "name"
	LabelPodUid        = PrefixPod + "uid"
	LabelPodName       = PrefixPod + "name"
	LabelContainerName = PrefixContainer + "name"
	LabelContainerId   = PrefixContainer + "id"
	LabelVolumeName    = PrefixVolume + "name"

	StatusPodMessage        = PrefixPod + "status.message"
	StatusPodPhase          = PrefixPod + "status.phase"
	StatusPodReason         = PrefixPod + "status.reason"
	StatusContainerReady    = PrefixContainer + "status.ready"
	StatusContainerRestarts = PrefixContainer + "status.restarts"
	StatusContainerState    = PrefixContainer + "status.state"
	StatusContainerExitCode = PrefixContainer + "status.exitcode"
	StatusContainerMessage  = PrefixContainer + "status.message"
	StatusContainerReason   = PrefixContainer + "status.reason"

	MeasureUptime                = "uptime"
	MeasureCpuTime               = "cpu.time"
	MeasureCpuUsage              = "cpu.usage"
	MeasureMemoryAvailable       = "memory.available"
	MeasureMemoryUsage           = "memory.usage"
	MeasureMemoryRSS             = "memory.rss"
	MeasureMemoryWorkingSet      = "memory.working_set"
	MeasureMemoryPageFaults      = "memory.page_faults"
	MeasureMemoryMajorPageFaults = "memory.major_page_faults"
	MeasureFilesystemAvailable   = "filesystem.available"
	MeasureFilesystemCapacity    = "filesystem.capacity"
	MeasureFilesystemUsage       = "filesystem.usage"
	MeasureNetworkBytesReceive   = "network.bytes.receive"
	MeasureNetworkBytesSend      = "network.bytes.send"
	MeasureNetworkErrorsReceive  = "network.errors.receive"
	MeasureNetworkErrorsSend     = "network.errors.send"
	MeasureVolumeAvailable       = "volume.available"
	MeasureVolumeCapacity        = "volume.capacity"
	MeasureVolumeInodesTotal     = "volume.inodes.total"
	MeasureVolumeInodesFree      = "volume.inodes.free"
	MeasureVolumeInodesUsed      = "volume.inodes.used"

	PrefixK8s       = "k8s."
	PrefixMetrics   = "metrics."
	PrefixLabel     = PrefixK8s + "label."
	PrefixCluster   = PrefixK8s + "cluster."
	PrefixNode      = PrefixK8s + "node."
	PrefixNamespace = PrefixK8s + "namespace."
	PrefixPod       = PrefixK8s + "pod."
	PrefixContainer = PrefixK8s + "container."
	PrefixVolume    = PrefixK8s + "volume."

	MetricSourceName       = "source"
	MetricSourceType       = "source.type"
	KubernetesResourceType = PrefixK8s + "resource.type"
)
