package main

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name = "SequentialScheduling"
)

type SequentialScheduling struct {
	handle framework.Handle
}

var _ framework.QueueSortPlugin = &SequentialScheduling{}

func (s *SequentialScheduling) Name() string {
	return Name
}

func (s *SequentialScheduling) Less(pInfo1, pInfo2 *framework.QueuedPodInfo) bool {
	pod1 := pInfo1.Pod
	pod2 := pInfo2.Pod

	if pod1 == nil || pod2 == nil {
		klog.ErrorS(nil, "Invalid pod info", "pod1", pod1, "pod2", pod2)
		return false
	}

	if pod1.CreationTimestamp.Equal(&pod2.CreationTimestamp) {
		return pod1.UID < pod2.UID
	}

	return pod1.CreationTimestamp.Before(&pod2.CreationTimestamp)
}

func New(_ context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	klog.InfoS("Creating new SequentialScheduling plugin")

	return &SequentialScheduling{
		handle: h,
	}, nil
}
