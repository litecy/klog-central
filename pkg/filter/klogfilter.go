package filter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
	"strings"
)

const (
	AnnotationKLogCentralLogConfigKey = "klc.klog.vibly.vip/logs-config."

	AnnotationKLogCentralLogConfigKeyData = "klc.klog.vibly.vip/logs-config-data."
)

func CheckKLCConfig(ctx context.Context, pod v1.Pod) (*entity.ConfigItems, error) {
	logger := log.FromContext(ctx)
	configs := make(entity.ConfigItems, 0)

	var keys = make([]string, 0)
	for k, _ := range pod.Annotations {
		if strings.HasPrefix(k, AnnotationKLogCentralLogConfigKey) || strings.HasPrefix(k, AnnotationKLogCentralLogConfigKeyData) {
			// get config from annotation
			keys = append(keys, k)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) > 0
	})

	for _, k := range keys {
		if strings.HasPrefix(k, AnnotationKLogCentralLogConfigKey) || strings.HasPrefix(k, AnnotationKLogCentralLogConfigKeyData) {
			// get config from annotation

			data := pod.Annotations[k]

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
				// continue
				return nil, errUn
			}
			configs = append(configs, item)
		}
	}

	return &configs, nil
}
