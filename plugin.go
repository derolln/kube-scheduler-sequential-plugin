package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/component-helpers/scheduling/corev1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name                = "SequentialScheduling"
	AttemptedAnnotation = "scheduling.sequential/attempted"
)

type SequentialScheduling struct {
	handle framework.Handle
}

var _ framework.QueueSortPlugin = &SequentialScheduling{}
var _ framework.PreFilterPlugin = &SequentialScheduling{}

func (s *SequentialScheduling) Name() string {
	return Name
}

func (s *SequentialScheduling) Less(pInfo1, pInfo2 *framework.QueuedPodInfo) bool {
	pod1 := pInfo1.Pod
	pod2 := pInfo2.Pod

	p1 := corev1.PodPriority(pod1)
	p2 := corev1.PodPriority(pod2)

	if pod1 == nil || pod2 == nil {
		klog.ErrorS(nil, "Invalid pod info", "pod1", pod1, "pod2", pod2)
		return false
	}

	if p1 == p2 && pod1.CreationTimestamp.Equal(&pod2.CreationTimestamp) {
		return pod1.UID < pod2.UID
	}

	return (p1 > p2) || (p1 == p2 && pod1.CreationTimestamp.Before(&pod2.CreationTimestamp))
	// original implementation : (p1 > p2) || (p1 == p2 && pInfo1.Timestamp.Before(pInfo2.Timestamp)
}

func (s *SequentialScheduling) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (s *SequentialScheduling) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	if _, ok := pod.Annotations[AttemptedAnnotation]; ok {
		klog.InfoS("Pod already attempted scheduling, proceeding", "pod", klog.KObj(pod))
		return nil, framework.NewStatus(framework.Success)
	}

	patch := []byte(fmt.Sprintf(`{"metadata":{"annotations":{%q:"true"}}}`, AttemptedAnnotation))
	_, err := s.handle.ClientSet().CoreV1().Pods(pod.Namespace).Patch(ctx, pod.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return nil, framework.AsStatus(fmt.Errorf("failed to annotate pod %s/%s: %w", pod.Namespace, pod.Name, err))
	}

	//klog.V(4).InfoS("First scheduling attempt, sending pod to unschedulable queue", "pod", klog.KObj(pod))
	klog.InfoS("First scheduling attempt, sending pod to unschedulable queue", "pod", klog.KObj(pod))
	return nil, framework.NewStatus(framework.Unschedulable, "first scheduling attempt, requeueing")
}

func New(_ context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	klog.InfoS("Creating new SequentialScheduling plugin v2")

	return &SequentialScheduling{
		handle: h,
	}, nil
}
