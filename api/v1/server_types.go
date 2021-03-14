/*


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

package v1

import (
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServerSpec defines the desired state of Server
type ServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Hostname of server.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=100
	// +kubebuilder:vaildation:Pattern=^(?=.*[a-zA-Z])([a-zA-Z0-9().-])+$
	// +kubebuilder:validation:Required
	Hostname string `json:"hostname,omitempty"`

	// Description of server.
	// +kubebuilder:validation:MaxLength=250
	Description string `json:"description,omitempty"`

	// OS ID used for server creation.
	// +kubebuilder:validation:Required
	OS ServerOS `json:"os,omitempty"`

	// Server type used for creation.
	// +kubebuilder:validation:Required
	Type ServerType `json:"type,omitempty"`

	// Location ID where the server is created.
	// +kubebuilder:validation:Required
	Location LocationID `json:"location,omitempty"`

	// Whether or not to install SSH Keys marked as default in additionl to any SSH keys speficied on this resource.
	// Defaults to true.
	InstallDefaultSSHKeys *bool `json:"installDefaultSshKeys"`

	// A list of SSH key IDs (BMC resource ID) that will be installed on the server in addition default SSH keys if enabled.
	// +kubebuilder:validation:Optional
	SSHKeyIDs []string `json:"sshKeyIds,omitempty"`

	// The type of networks where this server should be attached.
	// +kubebuilder:validation:Optional
	NetworkType NetworkType `json:"networkType,omitempty"`
}

// NetworkType represents the type of networking configuraiton a server should use.
// Only one of the following network types may be specified.
// If none of the following network types are specified, the default one is PublicAndPrivate.
// +kubebuilder:validation:Enum=PUBLIC_AND_PRIVATE;PRIVATE_ONLY
type NetworkType string

const (
	PublicAndPrivate NetworkType = `PUBLIC_AND_PRIVATE`
	PrivateOnly      NetworkType = `PRIVATE_ONLY`
)

// LocationID identifies a BMC region.
// Only one of the following locations may be specified.
// If none of the following locations are specified, the default one is Phoenix.
// +kubebuilder:validation:Enum=PHX;ASH;SGP;NLD
type LocationID string

const (
	Phoenix   LocationID = `PHX`
	Ashburn   LocationID = `ASH`
	Singapore LocationID = `SGP`
	Amsterdam LocationID = `NLD`
)

// ServerOS describes the operating system image for this server.
// Only one of the following server OSs may be specified.
// If none of the following OSs are specified, the default one is UbuntuBionic.
// +kubebuilder:validation:Enum=ubuntu/bionic;centos/centos7
type ServerOS string

const (
	UbuntuBionic  ServerOS = `ubuntu/bionic`
	CentosCentos7 ServerOS = `centos/centos7`
)

// ServerType describes the hardware to allocate for this server.
// Only one of the following server types may be specified.
// If none of the following types are specified, the default one is S1C1Small.
// +kubebuilder:validation:Enum=s1.c1.small;s1.c1.medium;s1.c2.medium;s1.c2.large;d1.c1.small;d1.c2.small;d1.c3.small;d1.c4.small;d1.c1.medium;d1.c2.medium;d1.c3.medium;d1.c4.medium;d1.c1.large;d1.c2.large;d1.c3.large;d1.c4.large;d1.m1.medium;d1.m2.medium;d1.m3.medium;d1.m4.medium
type ServerType string

const (
	S1C1Small  ServerType = `s1.c1.small`
	S1C1Medium ServerType = `s1.c1.medium`
	S1C2Medium ServerType = `s1.c2.medium`
	S1C2Large  ServerType = `s1.c2.large`
	D1C1Small  ServerType = `d1.c1.small`
	D1C2Small  ServerType = `d1.c2.small`
	D1C3Small  ServerType = `d1.c3.small`
	D1C4Small  ServerType = `d1.c4.small`
	D1C1Medium ServerType = `d1.c1.medium`
	D1C2Medium ServerType = `d1.c2.medium`
	D1C3Medium ServerType = `d1.c3.medium`
	D1C4Medium ServerType = `d1.c4.medium`
	D1C1Large  ServerType = `d1.c1.large`
	D1C2Large  ServerType = `d1.c2.large`
	D1C3Large  ServerType = `d1.c3.large`
	D1C4Large  ServerType = `d1.c4.large`
	D1M1Medium ServerType = `d1.m1.medium`
	D1M2Medium ServerType = `d1.m2.medium`
	D1M3Medium ServerType = `d1.m3.medium`
	D1M4Medium ServerType = `d1.m4.medium`
)

// ServerPricingModel describes the pricing model used for a specific server resource.
// One on of the following pricing models may be specified.
// If none of the following types are specified, the default one is PMHourly.
// +kubebuilder:validation:Enum=HOURLY;ONE_MONTH_RESERVATION;TWELVE_MONTHS_RESERVATION;TWENTY_FOUR_MONTHS_RESERVATION;THIRTY_SIX_MONTHS_RESERVATION
type ServerPricingModel string

const (
	PMHourly                      ServerPricingModel = `HOURLY`
	PMOneMonthReservation         ServerPricingModel = `ONE_MONTH_RESERVATION`
	PMTwelveMonthsReservation     ServerPricingModel = `TWELVE_MONTHS_RESERVATION`
	PMTwentyFourMonthsReservation ServerPricingModel = `TWENTY_FOUR_MONTHS_RESERVATION`
	PMThirtySixMonthsReservation  ServerPricingModel = `THIRTY_SIX_MONTHS_RESERVATION`
)

// ServerStatus defines the observed state of Server
type ServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	BMCServerID        string            `json:"id,omitempty"`
	BMCStatus          string            `json:"status,omitempty"`
	CPU                string            `json:"cpu,omitempty"`
	CPUCount           int32             `json:"cpuCount,omitempty"`
	CPUCores           int32             `json:"coresPerCpu,omitempty"`
	CPUFrequency       resource.Quantity `json:"cpuFrequency,omitempty"`
	Ram                string            `json:"ram,omitempty"`
	Storage            string            `json:"storage,omitempty"`
	PrivateIPAddresses []string          `json:"privateIpAddresses,omitempty"`
	PublicIPAddresses  []string          `json:"publicIpAddresses,omitempty"`
}

// +kubebuilder:object:root=true

// Server is the Schema for the servers API
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec,omitempty"`
	Status ServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerList contains a list of Server
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Server{}, &ServerList{})
}
