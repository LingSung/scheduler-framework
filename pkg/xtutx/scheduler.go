package xtutx

import (
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/apimachinery/pkg/runtime"
	"strconv"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const Name = "xtutx"

type Scheduler struct {
	handle framework.FrameworkHandle
}

var _ framework.QueueSortPlugin = &Scheduler{}
var _ framework.PreFilterPlugin = &Scheduler{}

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

func (s *Scheduler) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) *framework.Status {
	podGroup, podGroupExist := pod.Labels["podGroup"]
	if !podGroupExist {
		return framework.NewStatus(framework.Success, "podGroup label doesn't exist")
	}

	minAvailable, minAvailableExist := pod.Labels["minAvailable"]
	if !minAvailableExist {
		return framework.NewStatus(framework.Success, "minAvailable label doesn't exist")
	}

	minAvailableValue, err := strconv.Atoi(minAvailable)
	if err != nil {
		return framework.NewStatus(framework.Unschedulable, "atoi error")
	}

	totalNumOfPod := s.getTotalNumofPod(pod.Namespace, podGroup)
	if totalNumOfPod < minAvailableValue {
		return framework.NewStatus(framework.Unschedulable, "total number of pods is not enough")
	} 

	return framework.NewStatus(framework.Success, "PreFliter done")
}

func (s *Scheduler) getTotalNumofPod(namespace string, podgroup string) int {
	Podlist, err := s.handle.ClientSet().CoreV1().Pods(namespace).List(metav1.ListOptions{})

	if err != nil {
		return 0
	}

	total := 0
	for _, pod := range Podlist.Items {
		if podGroup, ok := pod.Labels["podGroup"]; ok && podGroup == podgroup {
			total++
		}
	}

	return total
}

// PreFilterExtensions ...
func (*Scheduler) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &Scheduler{
		handle: handle,
	}, nil
}