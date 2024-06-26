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

// Package helper contains hacky alternative implementations of resources that would otherwise not be able to be converted to yaml
package helper

import fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"

// CompositionRequest is a drop-in for RunFunctionRequest to replace ExtraResources
type CompositionRequest struct {
	*fnv1beta1.RunFunctionRequest

	ExtraResources map[string]*fnv1beta1.Resources `protobuf:"bytes,6,rep,name=extra_resources,json=extraResources,proto3" json:"extraResources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

// CompositionResponse is a drop-in for RunFunctionResponse to replace Requirements
type CompositionResponse struct {
	fnv1beta1.RunFunctionResponse

	Requirements *Requirements `json:"requirements,omitempty"`
}

// Requirements is a drop-in for the version provided by fnv1beta1
type Requirements struct {
	// fnv1beta1.Requirements

	ExtraResources map[string]*ResourceSelector `protobuf:"bytes,1,rep,name=extra_resources,json=extraResources,proto3" json:"extraResources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

// ResourceSelector allows to select Kubernetes Resources
type ResourceSelector struct {
	APIVersion string `protobuf:"bytes,1,opt,name=api_version,json=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Match      Match  `json:"match,omitempty"`
}

// Match specifies the method used for the Lookup
type Match struct {
	MatchName   string                 `protobuf:"bytes,3,opt,name=match_name,json=matchName,proto3,oneof" json:"matchName,omitempty"`
	MatchLabels *fnv1beta1.MatchLabels `protobuf:"bytes,4,opt,name=match_labels,json=matchLabels,proto3,oneof" json:"matchLabels,omitempty"`
}
