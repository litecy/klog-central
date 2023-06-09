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
				"klc.klog.vibly.vip/logs-config.1.a":    `{"format":"json","type":"file","path":"/var/log/app.log","name":"app.log"}`,
				"klc.klog.vibly.vip/logs-config.1.b":    `{"format":"regexp","type":"stdout","path":"stdout","name":"stdout"}`,
				"klc.klog.vibly.vip/logs-config-data.3": `eyJmb3JtYXQiOiAianNvbiIsICJ0eXBlIjogImZpbGUiLCAicGF0aCI6ICIvdmFyL2xvZy9zdGF0dXMuY29udHJvbCIsICJuYW1lIjogInN0YXR1cyIgfQ==`,
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
		entity.ConfigItem{Format: "json", Type: "file", Path: "/var/log/app.log", Name: "app.log", Container: "a"},
		entity.ConfigItem{Format: "regexp", Type: "stdout", Path: "stdout", Name: "stdout", Container: "b"},
		entity.ConfigItem{Format: "json", Type: "file", Path: "/var/log/status.control", Name: "status"},
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
				"klc.klog.vibly.vip/logs-config.1": `invalid json`,
			},
		},
	}

	configs, err = CheckKLCConfig(ctx, testPod)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}
