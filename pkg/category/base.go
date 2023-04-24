package category

import (
	"errors"
	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
)

const (
	LogTypeStdout = "stdout"
	LogTypeFile   = "file"
)

type LogType interface {
	Name() string
	GetPattern(pod *v1.Pod, item *entity.ConfigItemForProcess) error
}

var logTypes = make(map[string]LogType)

// RegisterLogType bind logtype to name
func RegisterLogType(lt LogType) {
	logTypes[lt.Name()] = lt
}

func GetLogPattern(name string, pod *v1.Pod, item *entity.ConfigItemForProcess) error {
	if lt, has := logTypes[name]; !has {
		return errors.New("log type not found: " + name)
	} else {
		return lt.GetPattern(pod, item)
	}
}
