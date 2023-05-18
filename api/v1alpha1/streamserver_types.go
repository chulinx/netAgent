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

// StreamServerSpec defines the desired state of StreamServer
type StreamServerSpec struct {
	Name       string `json:"name,omitempty"`
	ListenPort int32  `json:"listenPort,omitempty"`
	Proxy      Proxy  `json:"proxy,omitempty"`
}

type Proxy struct {
	NameSpace string `json:"nameSpace,omitempty"`
	Service   string `json:"service,omitempty"`
	Port      int32  `json:"port,omitempty"`
}

// StreamServerStatus defines the observed state of StreamServer
type StreamServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ListenPort",type="string",JSONPath=".spec.listenPort",description="The Nginx server listen_port"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// StreamServer is the Schema for the streamservers API
type StreamServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StreamServerSpec   `json:"spec,omitempty"`
	Status StreamServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StreamServerList contains a list of StreamServer
type StreamServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StreamServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StreamServer{}, &StreamServerList{})
}
