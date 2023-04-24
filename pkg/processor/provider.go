package processor

import (
	v1 "k8s.io/api/core/v1"
	"path/filepath"
)

type ConfigProvider interface {
	GetConfPath(pod *v1.Pod) string
}

var _ ConfigProvider = &FileBeatConfigProvider{}

type FileBeatConfigProvider struct {
}

const FilebeatConfPath = "/usr/share/filebeat/prospectors.d"

func (*FileBeatConfigProvider) GetConfPath(pod *v1.Pod) string {
	return filepath.Join(FilebeatConfPath, string(pod.UID)+".yml")
}

type DevNullConfigProvider struct {
}

func (*DevNullConfigProvider) GetConfPath(pod *v1.Pod) string {
	return filepath.Join("fake-host", string(pod.UID)+".yml")
}
