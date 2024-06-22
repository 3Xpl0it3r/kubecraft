/*
Copyright 2024 The github.com/3Xpl0it3r/kubecraft Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// Bar defines Bar deployment
type Bar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BarSpec   `json:"spec"`
	Status BarStatus `json:"status"`
}

// BarSpec describes the specification of Bar applications using kubernetes as a cluster manager
type BarSpec struct {
	// TODO, write your code
}

// BarStatus describes the current status of Bar applications
type BarStatus struct {
	// TODO, write your code
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BarList carries a list of Bar objects
type BarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Bar `json:"items"`
}
