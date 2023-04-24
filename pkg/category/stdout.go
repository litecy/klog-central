package category

import (
	"errors"
	"fmt"
	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
)

var _ LogType = &Stdout{}

func init() {
	RegisterLogType(&Stdout{})
}

type Stdout struct {
}

func (s *Stdout) Name() string {
	return LogTypeStdout
}

func (s *Stdout) GetPattern(pod *v1.Pod, item *entity.ConfigItemForProcess) error {
	if pod == nil || item == nil {
		return errors.New("pod or item is nil")
	}

	var c *v1.Container

	for _, container := range pod.Spec.Containers {
		if container.Name == item.Container || item.Container == "" {
			c = &container
			break
		}
	}

	if c == nil {
		for _, container := range pod.Spec.InitContainers {
			if container.Name == item.Container || item.Container == "" {
				c = &container
				break
			}
		}

		if c == nil {
			for _, container := range pod.Spec.EphemeralContainers {
				if container.Name == item.Container || item.Container == "" {
					cc := v1.Container(container.EphemeralContainerCommon)
					c = &cc
					break
				}
			}
		}
	}

	if c == nil {
		return errors.New("container not found: " + item.Container)
	}

	logPath := fmt.Sprintf("/var/log/pods/%s_%s_%s/%s/*.log", pod.Namespace, pod.Name, pod.UID, c.Name)

	item.ContainerPath = logPath

	// /var/lib/kubelet/pods/fec63131-cd3a-402d-9f7a-d036b9d368ed/volumes/kubernetes.io~empty-dir/log-dir
	return nil
}
