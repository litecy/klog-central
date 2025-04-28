package format

import "fmt"

// Converter converts node info to map
type Converter func(info map[string]any) (map[string]string, error)

var converters = make(map[string]Converter)

// RegisterFormatter format converter instance
func RegisterFormatter(format string, converter Converter) {
	converters[format] = converter
}

// Convert convert node info to map
func Convert(format string, info map[string]any) (map[string]string, error) {
	converter := converters[format]
	if converter == nil {
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}
	return converter(info)
}

// SimpleConverter simple format converter
type SimpleConverter struct {
	properties map[string]bool
}

var defaultConverter = func(properties []string) Converter {
	return func(info map[string]any) (map[string]string, error) {
		validProperties := make(map[string]bool)
		for _, property := range properties {
			validProperties[property] = true
		}
		ret := make(map[string]string)
		return ret, nil
	}
}

func init() {

	RegisterFormatter("default", defaultConverter([]string{}))
	RegisterFormatter("csv", defaultConverter([]string{"time_key", "time_format", "keys"}))
	RegisterFormatter("json", defaultConverter([]string{"time_key", "time_format"}))
	RegisterFormatter("regexp", defaultConverter([]string{"time_key", "time_format"}))
	RegisterFormatter("apache2", defaultConverter([]string{}))
	RegisterFormatter("apache_error", defaultConverter([]string{}))
	RegisterFormatter("nginx", defaultConverter([]string{}))
	RegisterFormatter("containerd", defaultConverter([]string{}))
	RegisterFormatter("containerd_json", defaultConverter([]string{}))
}
