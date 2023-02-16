/*
Copyright 2021 RadonDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MysqlClusterSpec defines the desired state of MysqlCluster
type MysqlClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is the number of pods.
	// +optional
	// +kubebuilder:validation:Enum=0;1;2;3;5
	// +kubebuilder:default:=3
	Replicas *int32 `json:"replicas,omitempty"`

	// Username of new user to create.
	// Only be a combination of letters, numbers or underlines. The length can not exceed 26 characters.
	// +optional
	// +kubebuilder:default:="radondb_usr"
	// +kubebuilder:validation:Pattern="^[A-Za-z0-9_]{2,26}$"
	User string `json:"user,omitempty"`

	// MySQLConfig `ConfigMap` name of MySQL config.
	// +optional
	MySQLConfig *string `json:"mysqlConfig,omitempty"`

	//Compute resources of a MySQL container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Containing CA (ca.crt) and server cert (tls.crt), server private key (tls.key) for SSL
	// +optional
	CustomTLSSecret *corev1.SecretProjection `json:"customTLSSecret,omitempty"`

	// Defines a PersistentVolumeClaim for MySQL data.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes
	// +kubebuilder:validation:Required
	Storage corev1.PersistentVolumeClaimSpec `json:"storage"`

	// Represents the MySQL version that will be run. The available version can be found here:
	// This field should be set even if the Image is set to let the operator know which mysql version is running.
	// Based on this version the operator can take decisions which features can be used.
	// +optional
	// +kubebuilder:default:="5.7"
	MysqlVersion string `json:"mysqlVersion,omitempty"`

	// DatabaseInitSQL defines a ConfigMap containing custom SQL that will
	// be run after the cluster is initialized. This ConfigMap must be in the same
	// namespace as the cluster.
	// +optional
	DatabaseInitSQL *DatabaseInitSQL `json:"databaseInitSQL,omitempty"`

	// XenonOpts is the options of xenon container.
	// +optional
	// +kubebuilder:default:={image: "radondb/xenon:v2.3.0", admitDefeatHearbeatCount: 5, electionTimeout: 10000, resources: {limits: {cpu: "100m", memory: "256Mi"}, requests: {cpu: "50m", memory: "128Mi"}}}
	Xenon XenonOpts `json:"xenonOpts,omitempty"`

	// Backup is the options of backup container.
	// +optional
	Backup BackupOpts `json:"backupOpts,omitempty"`

	// Monitoring is the options of metrics container.
	// +optional
	// +kubebuilder:default:={image: "prom/mysqld-exporter:v0.12.1", resources: {limits: {cpu: "100m", memory: "128Mi"}, requests: {cpu: "10m", memory: "32Mi"}}, enabled: false}
	Monitoring MonitoringSpec `json:"MonitoringSpec,omitempty"`

	// ImagePullPolicy is used to determine when Kubernetes will attempt to
	// pull (download) container images.
	// More info: https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy
	// +kubebuilder:validation:Enum={Always,Never,IfNotPresent}
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Tolerations of a MySQL pod. Changing this value causes MySQL to restart.
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Scheduling constraints of MySQL pod. Changing this value causes
	// MySQL to restart.
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Priority class name for the MySQL pods. Changing this value causes
	// MySQL to restart.
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/
	// +optional
	PriorityClassName *string `json:"priorityClassName,omitempty"`

	// The number of pods from that set that must still be available after the
	// eviction, even in the absence of the evicted pod
	// +optional
	// +kubebuilder:default:="50%"
	MinAvailable string `json:"minAvailable,omitempty"`

	// Specifies a data source for bootstrapping the MySQL cluster.
	// +optional
	DataSource *DataSource `json:"dataSource,omitempty"`

	// Run this cluster as a read-only copy of an existing cluster or archive.
	// +optional
	Standby *MySQLStandbySpec `json:"standby,omitempty"`

	// If true, when the data is inconsistent, Xenon will automatically rebuild the invalid node.
	// +optional
	// +kubebuilder:default:=false
	EnableAutoRebuild bool `json:"enableAutoRebuild,omitempty"`

	// TODO:The specification of a proxy that connects to MySQL.
	// +optional
	// Proxy *MySQLProxySpec `json:"proxy,omitempty"`

	// LogOpts is the options of log settings.
	// +optional
	Log LogOpts `json:"logOpts,omitempty"`

	// Specification of the service that exposes the MySQL leader instance.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`
}

type BackupOpts struct {
	// Image is the image of backup container.
	// +optional
	Image string `json:"image,omitempty"`
	// Changing this value causes MySQL
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// MysqlClusterStatus defines the observed state of MysqlCluster
type MysqlClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ReadyNodes represents number of the nodes that are in ready state.
	ReadyNodes int `json:"readyNodes,omitempty"`
	// State
	State ClusterState `json:"state,omitempty"`
	// Conditions contains the list of the cluster conditions fulfilled.
	Conditions []ClusterCondition `json:"conditions,omitempty"`
	// Nodes contains the list of the node status fulfilled.
	Nodes []NodeStatus `json:"nodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.readyNodes
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="The cluster status"
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="Current",type="integer",JSONPath=".status.readyNodes",description="The number of current replicas"
// +kubebuilder:printcolumn:name="Leader",type="string",JSONPath=".status.nodes[?(@.raftStatus.role == 'LEADER')].name",description="Name of the leader node"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=mysql
// MysqlCluster is the Schema for the mysqlclusters API
type MysqlCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MysqlClusterSpec   `json:"spec,omitempty"`
	Status MysqlClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// MysqlClusterList contains a list of MysqlCluster
type MysqlClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MysqlCluster `json:"items"`
}

// ClusterState defines cluster state.
type ClusterState string

const (
	// ClusterInitState indicates whether the cluster is initializing.
	ClusterInitState ClusterState = "Initializing"
	// ClusterUpdateState indicates whether the cluster is being updated.
	ClusterUpdateState ClusterState = "Updating"
	// ClusterReadyState indicates whether all containers in the pod are ready.
	ClusterReadyState ClusterState = "Ready"
	// ClusterCloseState indicates whether the cluster is closed.
	ClusterCloseState ClusterState = "Closed"
	// ClusterScaleInState indicates whether the cluster replicas is decreasing.
	ClusterScaleInState ClusterState = "ScaleIn"
	// ClusterScaleOutState indicates whether the cluster replicas is increasing.
	ClusterScaleOutState ClusterState = "ScaleOut"
)

// ClusterConditionType defines type for cluster condition type.
type ClusterConditionType string

const (
	// ConditionInit indicates whether the cluster is initializing.
	ConditionInit ClusterConditionType = "Initializing"
	// ConditionUpdate indicates whether the cluster is being updated.
	ConditionUpdate ClusterConditionType = "Updating"
	// ConditionReady indicates whether all containers in the pod are ready.
	ConditionReady ClusterConditionType = "Ready"
	// ConditionClose indicates whether the cluster is closed.
	ConditionClose ClusterConditionType = "Closed"
	// ConditionError indicates whether there is an error in the cluster.
	ConditionError ClusterConditionType = "Error"
	// ConditionScaleIn indicates whether the cluster replicas is decreasing.
	ConditionScaleIn ClusterConditionType = "ScaleIn"
	// ConditionScaleOut indicates whether the cluster replicas is increasing.
	ConditionScaleOut ClusterConditionType = "ScaleOut"
)

// ClusterCondition defines type for cluster conditions.
type ClusterCondition struct {
	// Type of cluster condition, values in (\"Initializing\", \"Ready\", \"Error\").
	Type ClusterConditionType `json:"type"`
	// Status of the condition, one of (\"True\", \"False\", \"Unknown\").
	Status corev1.ConditionStatus `json:"status"`

	// The last time this Condition type changed.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// One word, camel-case reason for current status of the condition.
	Reason string `json:"reason,omitempty"`
	// Full text reason for current status of the condition.
	Message string `json:"message,omitempty"`
}

// NodeStatus defines type for status of a node into cluster.
type NodeStatus struct {
	// Name of the node.
	Name string `json:"name"`
	// Full text reason for current status of the node.
	Message string `json:"message,omitempty"`
	// RaftStatus is the raft status of the node.
	RaftStatus RaftStatus `json:"raftStatus,omitempty"`
	// Conditions contains the list of the node conditions fulfilled.
	Conditions []NodeCondition `json:"conditions,omitempty"`
}

type RaftStatus struct {
	// Role is one of (LEADER/CANDIDATE/FOLLOWER/IDLE/INVALID)
	Role string `json:"role,omitempty"`
	// Leader is the name of the Leader of the current node.
	Leader string `json:"leader,omitempty"`
	// Nodes is a list of nodes that can be identified by the current node.
	Nodes []string `json:"nodes,omitempty"`
}

// NodeCondition defines type for representing node conditions.
type NodeCondition struct {
	// Type of the node condition.
	Type NodeConditionType `json:"type"`
	// Status of the node, one of (\"True\", \"False\", \"Unknown\").
	Status corev1.ConditionStatus `json:"status"`
	// The last time this Condition type changed.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

// The index of the NodeStatus.Conditions.
type NodeConditionsIndex uint8

const (
	IndexLagged NodeConditionsIndex = iota
	IndexLeader
	IndexReadOnly
	IndexReplicating
)

// NodeConditionType defines type for node condition type.
type NodeConditionType string

const (
	// NodeConditionLagged represents if the node is lagged.
	NodeConditionLagged NodeConditionType = "Lagged"
	// NodeConditionLeader represents if the node is leader or not.
	NodeConditionLeader NodeConditionType = "Leader"
	// NodeConditionReadOnly repesents if the node is read only or not
	NodeConditionReadOnly NodeConditionType = "ReadOnly"
	// NodeConditionReplicating represents if the node is replicating or not.
	NodeConditionReplicating NodeConditionType = "Replicating"
)

// DatabaseInitSQL defines a ConfigMap containing custom SQL that will
// be run after the cluster is initialized. This ConfigMap must be in the same
// namespace as the cluster.
type DatabaseInitSQL struct {
	// Name is the name of a ConfigMap
	// +required
	Name string `json:"name"`

	// Key is the ConfigMap data key that points to a SQL string
	// +required
	Key string `json:"key"`
}

type XenonOpts struct {
	// To specify the image that will be used for xenon container.
	// +optional
	// +kubebuilder:default:="radondb/xenon:v2.3.0"
	Image string `json:"image,omitempty"`

	// High available component admit defeat heartbeat count.
	// +optional
	// +kubebuilder:default:=5
	AdmitDefeatHearbeatCount *int32 `json:"admitDefeatHearbeatCount,omitempty"`

	// High available component election timeout. The unit is millisecond.
	// +optional
	// +kubebuilder:default:=10000
	ElectionTimeout *int32 `json:"electionTimeout,omitempty"`

	// If true, when the data is inconsistent, Xenon will automatically rebuild the invalid node.
	// +optional
	// +kubebuilder:default:=false
	EnableAutoRebuild bool `json:"enableAutoRebuild,omitempty"`

	// The compute resource requirements.
	// +optional
	// +kubebuilder:default:={limits: {cpu: "100m", memory: "256Mi"}, requests: {cpu: "50m", memory: "128Mi"}}
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type MonitoringSpec struct {
	// +optional
	Exporter *ExporterSpec `json:"exporter,omitempty"`
}

type ExporterSpec struct {

	// Projected secret containing custom TLS certificates to encrypt output from the exporter
	// web server
	// +optional
	CustomTLSSecret *corev1.SecretProjection `json:"customTLSSecret,omitempty"`

	// To specify the image that will be used for metrics container.
	// +optional
	// +kubebuilder:default:="prom/mysqld-exporter:v0.12.1"
	Image string `json:"image,omitempty"`

	// Changing this value causes MySQL and the exporter to restart.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type DataSource struct {
	// Bootstraping from remote data source
	// +optional
	Remote *RemoteDataSource `json:"remote,omitempty"`
}

type RemoteDataSource struct {
	//
	SourceConfig *corev1.SecretProjection `json:"sourceConfig,omitempty"`
}

type LogOpts struct {
	// To specify the image that will be used for log container.
	// +optional

}

type MySQLStandbySpec struct {
	// Whether or not the MySQL cluster should be read-only. When this is
	// true, the cluster will be read-only. When this is false, the cluster will
	// run as writable.
	// +optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// The name of the MySQL cluster to follow for binlog.
	// +optional
	ClusterName string `json:"clusterName,omitempty"`

	// Network address of the MySQL server to follow via via binlog replication.
	// +optional
	Host string `json:"host,omitempty"`

	// Network port of the MySQL server to follow via binlog replication.
	// +optional
	// +kubebuilder:validation:Minimum=1024
	Port *int32 `json:"port,omitempty"`
}

type ServiceSpec struct {
	// The port on which this service is exposed when type is NodePort or
	// LoadBalancer. Value must be in-range and not in use or the operation will
	// fail. If unspecified, a port will be allocated if this Service requires one.
	// - https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +optional
	NodePort *int32 `json:"nodePort,omitempty"`

	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
	//
	// +optional
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum={ClusterIP,NodePort,LoadBalancer}
	Type string `json:"type"`
}

func init() {
	SchemeBuilder.Register(&MysqlCluster{}, &MysqlClusterList{})
}