package filter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/litecy/klog-central/pkg/entity"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	AnnotationKLogCentralLogConfigKey = "klc.klog.vibly.vip/logs-config."

	AnnotationKLogCentralLogConfigKeyData = "klc.klog.vibly.vip/logs-config-data."

	AnnotationKLogCentralLogMetaPrefixKey = "klc.klog.vibly.vip/logs-meta-prefix"
)

type AnnoKey struct {
	Key       string
	Prefix    string
	Index     int
	Container string
}

func parseAnnoKeySection(key string) *AnnoKey {
	annoKey := &AnnoKey{
		Key:       key,
		Prefix:    "",
		Index:     0,
		Container: "",
	}

	prefix := ""
	if strings.HasPrefix(key, AnnotationKLogCentralLogConfigKey) {
		prefix = AnnotationKLogCentralLogConfigKey
	} else if strings.HasPrefix(key, AnnotationKLogCentralLogConfigKeyData) {
		prefix = AnnotationKLogCentralLogConfigKeyData
	} else if strings.EqualFold(key, AnnotationKLogCentralLogMetaPrefixKey) {
		prefix = ""
		return annoKey
	} else {
		return nil
	}

	annoKey.Prefix = prefix

	keyData := strings.Replace(key, prefix, "", 1)
	keyDataSec := strings.SplitN(keyData, ".", 2)
	if len(keyDataSec) >= 1 {
		annoKey.Index, _ = strconv.Atoi(keyDataSec[0])
	}
	if len(keyDataSec) >= 2 {
		annoKey.Container = keyDataSec[1]
	}
	return annoKey
}

func CheckKLCConfig(ctx context.Context, pod v1.Pod) (*entity.ConfigItems, error) {
	logger := log.FromContext(ctx)
	configs := make(entity.ConfigItems, 0)

	var keys = make([]*AnnoKey, 0)
	for k, _ := range pod.Annotations {
		sec := parseAnnoKeySection(k)
		if sec != nil {
			keys = append(keys, sec)
		}
	}

	metaPrefix := ""

	for _, k := range keys {
		if k.Key == AnnotationKLogCentralLogMetaPrefixKey {
			// 处理meta
			metaPrefix = pod.Annotations[k.Key]
			break
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Index == keys[j].Index {
			return strings.Compare(keys[i].Container, keys[j].Container) < 0
		}
		return keys[i].Index < keys[j].Index
	})

	for _, k := range keys {
		// get config from annotation
		if k.Prefix == "" {
			// meta
			continue
		}

		data := pod.Annotations[k.Key]

		if k.Prefix == AnnotationKLogCentralLogConfigKeyData {
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

		item.Name = metaPrefix + item.Name

		item.Container = k.Container

		configs = append(configs, item)
	}

	return &configs, nil
}
