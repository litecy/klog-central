package util

func PutIfNotEmpty(store map[string]string, key, value string) {
	if key == "" || value == "" || store == nil {
		return
	}
	store[key] = value
}
