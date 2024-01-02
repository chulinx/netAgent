/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualServerSpec defines the desired state of VirtualServer
type VirtualServerSpec struct {
	ListenPort int32      `json:"listenPort,omitempty"`
	ServerName string     `json:"serverName,omitempty"`
	Proxys     []Location `json:"proxys,omitempty"`
}

// VirtualServerStatus defines the observed state of VirtualServer
type VirtualServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Status",xDescriptors="urn:alm:descriptor:com.tectonic.ui:text"
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ListenPort",type="string",JSONPath=".spec.listenPort",description="The Nginx server listen_port"
//+kubebuilder:printcolumn:name="ServerName",type="string",JSONPath=".spec.serverName",description="The Nginx server server_name"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualServer is the Schema for the virtual-servers API
type VirtualServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualServerSpec   `json:"spec,omitempty"`
	Status VirtualServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VirtualServerList contains a list of VirtualServer
type VirtualServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualServer `json:"items"`
}

// Location is Nginx Location
type Location struct {
	Name   string `json:"name,omitempty"`
	Path   string `json:"path,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	// NameSpace empty is current namespace
	NameSpace string `json:"nameSpace,omitempty"`
	Service   string `json:"service,omitempty"`
	Port      int32  `json:"port,omitempty"`
	// ProxyRedirect set proxy_redirect     off;
	ProxyRedirect    bool              `json:"proxyRedirect,omitempty"`
	ProxyHttpVersion string            `json:"proxyHttpVersion,omitempty"`
	ProxyHeaders     map[string]string `json:"proxyHeaders,omitempty"`
}

func init() {
	SchemeBuilder.Register(&VirtualServer{}, &VirtualServerList{})
}
