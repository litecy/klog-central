package category

import (
	"errors"
	"github.com/litecy/klog-central/pkg/entity"
	"github.com/litecy/klog-central/pkg/util"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
)

var _ LogType = &File{}

func init() {
	RegisterLogType(&File{})
}

type File struct {
}

func (f *File) Name() string {
	return LogTypeFile
}

func (f *File) GetPattern(pod *v1.Pod, item *entity.ConfigItemForProcess) error {
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

	mountName := ""
	relPath := ""
	vm := c.VolumeMounts
	for _, v := range vm {

		var err error
		var r string
		r, err = filepath.Rel(v.MountPath, item.Path)
		if err != nil {
			// not under mount path
			continue
		}

		if strings.HasSuffix(relPath, "..") || strings.HasPrefix(relPath, ".") {
			// not under mount path
			continue
		}

		mountName = v.Name
		relPath = r

		break
	}

	logPath := filepath.Join("/var/lib/kubelet/pods", string(pod.UID), "volumes", "kubernetes.io~empty-dir", mountName, relPath)

	item.ContainerPath = logPath

	util.PutIfNotEmpty(item.Meta, "k8s.container.name", c.Name)

	// /var/lib/kubelet/pods/fec63131-cd3a-402d-9f7a-d036b9d368ed/volumes/kubernetes.io~empty-dir/log-dir
	return nil
}
