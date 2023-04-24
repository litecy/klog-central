package processor

import (
	"bytes"
	"context"
	"errors"
	"github.com/Masterminds/sprig/v3"
	"github.com/litecy/klog-central/pkg/category"
	"github.com/litecy/klog-central/pkg/entity"
	"github.com/litecy/klog-central/pkg/util"
	v1 "k8s.io/api/core/v1"
	"os"
	"text/template"
)

type KConfig struct {
	confTemplate     *template.Template
	confTemplateFile string
	logHost          string
	cProvider        ConfigProvider
}

func NewKConfig(logHost string, tplFile string, provider ConfigProvider) *KConfig {
	return &KConfig{
		confTemplateFile: tplFile,
		logHost:          logHost,
		cProvider:        provider,
	}
}

func (c *KConfig) Setup() error {
	if c.confTemplateFile == "" {
		return errors.New("conf template file is empty")
	}

	data, err := os.ReadFile(c.confTemplateFile)
	if err != nil {
		return err
	}
	tpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(string(data))
	if err != nil {
		return err
	}

	c.confTemplate = tpl

	return nil
}

func (c *KConfig) Process(ctx context.Context, items *entity.ConfigItems, pod *v1.Pod) (err error) {

	var itemProcesses = make([]entity.ConfigItemForProcess, 0, len(*items))

	for _, item := range *items {
		itemProcess := entity.ConfigItemForProcess{ConfigItem: item, Meta: make(map[string]string), AddFields: make(map[string]string)}

		errIn := category.GetLogPattern(item.Type, pod, &itemProcess)
		if errIn != nil {
			continue
		}

		itemProcess.HostPath = c.logHost

		itemProcesses = append(itemProcesses, itemProcess)

		util.PutIfNotEmpty(itemProcess.Meta, "k8s_pod", pod.Name)
		util.PutIfNotEmpty(itemProcess.Meta, "k8s_namespace", pod.Namespace)
		util.PutIfNotEmpty(itemProcess.Meta, "k8s_node_name", "")
		util.PutIfNotEmpty(itemProcess.Meta, "k8s_node", pod.Spec.NodeName)

		util.PutIfNotEmpty(itemProcess.AddFields, "service", itemProcess.Name)
	}

	// category.GetLogPattern()

	var buf bytes.Buffer

	if err = c.confTemplate.Execute(&buf, itemProcesses); err != nil {
		return err
	}

	conf := c.cProvider.GetConfPath(pod)
	if conf == "" {
		return errors.New("conf path is empty")
	}

	if err = os.WriteFile(conf, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
