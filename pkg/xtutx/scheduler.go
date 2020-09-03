package xtutx

import (
	"k8s.io/kubernetes/pkg/api/v1/pod"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/apimachinery/pkg/runtime"
	"strconv"
)

const Name = "xtutx"

type Scheduler struct {}

var _ framework.QueueSortPlugin = &Scheduler{}

func (*Scheduler) Name() string {
	return Name
}

func (*Scheduler) Less(pInfo1, pInfo2 *framework.PodInfo) bool {
	p1 := GetPodPriority(pInfo1.Pod)
	p2 := GetPodPriority(pInfo2.Pod)
	return (p1 > p2) || (p1 == p2 && compQOS(pInfo1.Pod, pInfo2.Pod))
}

func compQOS(p1, p2 *v1.Pod) bool {
	p1QOS, p2QOS := v1qos.GetPodQOS(p1), v1qos.GetPodQOS(p2)
	if p1QOS == v1.PodQOSGuaranteed {
		return true
	}
	if p1QOS == v1.PodQOSBurstable {
		if p2QOS == v1.PodQOSGuaranteed {
			return false
		}
		return true
	}
	return p2QOS == v1.PodQOSBestEffort
}

func GetPodPriority(pod *v1.Pod) int {
	groupPriority, ok := pod.Labels["groupPriority"]
	if !ok {
		return -1
	}
	groupPriorityValue, err := strconv.Atoi(groupPriority)
	if (err != nil) {
		return -1
	}
	return groupPriorityValue
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &Scheduler{}, nil
}