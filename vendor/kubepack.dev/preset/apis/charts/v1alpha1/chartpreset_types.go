/*
Copyright 2022.

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
	"k8s.io/apimachinery/pkg/runtime"
	"kmodules.xyz/client-go/apiextensions"
	"kubepack.dev/preset/crds"
)

const (
	ResourceKindChartPreset = "ChartPreset"
	ResourceChartPreset     = "chartpreset"
	ResourceChartPresets    = "chartpresets"
)

// ChartPresetSpec defines the desired state of ChartPreset
type ChartPresetSpec struct {
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *runtime.RawExtension `json:"values,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChartPreset is the Schema for the chartpresets API
type ChartPreset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ChartPresetSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// ChartPresetList contains a list of ChartPreset
type ChartPresetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChartPreset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChartPreset{}, &ChartPresetList{})
}

func (_ ChartPreset) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourceChartPresets))
}
