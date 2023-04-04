package filter

import (
	"context"
	"testing"

	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckKLCConfig(t *testing.T) {
	// create a test pod with valid klc config annotations
	testPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Annotations: map[string]string{
				"klc.klog.uiey.vip/logs-config.1":      `{"format":"json","type":"file","path":"/var/log/app.log","name":"app.log"}`,
				"klc.klog.uiey.vip/logs-config.2":      `{"format":"regexp","type":"stdout","path":"stdout","name":"stdout"}`,
				"klc.klog.uiey.vip/logs-config-data.3": `eyJmb3JtYXQiOiAiandzb24iLCAidHlwZSI6ICJmaWxlIiwgInBhdGgiOiAiL3Zhci9sb2cvc3RhdHVzLmNvbnRyb2wiLCAibmFtZSI6ICJzdGF0dXMiIH0=`,
			},
		},
	}

	ctx := context.TODO()
	configs, err := CheckKLCConfig(ctx, testPod)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check if the parsed configs match the expected ones
	expectedConfigs := entity.ConfigItems{
		entity.ConfigItem{Format: "json", Type: "file", Path: "/var/log/app.log", Name: "app.log"},
		entity.ConfigItem{Format: "regexp", Type: "stdout", Path: "stdout", Name: "stdout"},
		entity.ConfigItem{Format: "json", Type: "file", Path: "/var/log/status.controller", Name: "status"},
	}
	if len(*configs) != len(expectedConfigs) {
		t.Errorf("unexpected number of configs: got %d, want %d", len(*configs), len(expectedConfigs))
	}

	for i, c := range *configs {
		if c != expectedConfigs[i] {
			t.Errorf("unexpected config item: got %+v, want %+v", c, expectedConfigs[i])
		}
	}

	// create a test pod with invalid klc config annotations
	testPod = v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Annotations: map[string]string{
				"klc.klog.uiey.vip/logs-config.1": `invalid json`,
			},
		},
	}

	configs, err = CheckKLCConfig(ctx, testPod)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}
