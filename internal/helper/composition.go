package helper

import fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"

type CompositionRequest struct {
	fnv1beta1.RunFunctionRequest

	ExtraResources map[string]*fnv1beta1.Resources `protobuf:"bytes,6,rep,name=extra_resources,json=extraResources,proto3" json:"extraResources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type CompositionResponse struct {
	fnv1beta1.RunFunctionResponse

	Requirements *Requirements `json:"requirements,omitempty"`
}

type Requirements struct {
	//fnv1beta1.Requirements

	ExtraResources map[string]*ResourceSelector `protobuf:"bytes,1,rep,name=extra_resources,json=extraResources,proto3" json:"extraResources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type ResourceSelector struct {
	ApiVersion string `protobuf:"bytes,1,opt,name=api_version,json=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Match      Match  `json:"match,omitempty"`
}

type Match struct {
	MatchName   string                 `protobuf:"bytes,3,opt,name=match_name,json=matchName,proto3,oneof" json:"matchName,omitempty"`
	MatchLabels *fnv1beta1.MatchLabels `protobuf:"bytes,4,opt,name=match_labels,json=matchLabels,proto3,oneof" json:"matchLabels,omitempty"`
}