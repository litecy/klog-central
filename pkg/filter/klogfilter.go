package filter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

const (
	AnnotationKLogCentralLogConfigKey = "klc.klog.uiey.vip/logs-config."

	AnnotationKLogCentralLogConfigKeyData = "klc.klog.uiey.vip/logs-config-data."
)

func CheckKLCConfig(ctx context.Context, pod v1.Pod) (*entity.ConfigItems, error) {
	logger := log.FromContext(ctx)
	configs := make(entity.ConfigItems, 0)
	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, AnnotationKLogCentralLogConfigKey) {
			// get config from annotation

			data := v

			if strings.HasPrefix(k, AnnotationKLogCentralLogConfigKeyData) {
				// config is BASE64, try to decode first
				dec, errD := base64.StdEncoding.DecodeString(data)
				if errD != nil {
					logger.Error(errD, "decode failed", "pod", pod.Name, "key", k, "value", data)
					continue
				}

				data = string(dec)
			}

			var item entity.ConfigItem
			errUn := json.Unmarshal([]byte(data), &item)
			if errUn != nil {
				logger.Error(errUn, "parse json failed", "pod", pod.Name, "key", k, "value", data)
				continue
			}
			configs = append(configs, item)
		}
	}

	return &configs, nil
}
