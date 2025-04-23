package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GPUTracker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	GPUNodes string `json:"gpu_nodes,omitempty"`
}

func (in *GPUTracker) DeepCopyObject() runtime.Object {
	out := new(GPUTracker)
	*out = *in
	out.ObjectMeta = *in.ObjectMeta.DeepCopy()
	return out
}

type GPUTrackerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GPUTracker `json:"items"`
}

func (in *GPUTrackerList) DeepCopyObject() runtime.Object {
	out := new(GPUTrackerList)
	*out = *in
	out.Items = make([]GPUTracker, len(in.Items))
	for i := range in.Items {
		out.Items[i] = *in.Items[i].DeepCopyObject().(*GPUTracker)
	}
	return out
}
