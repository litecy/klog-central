package entity

type ConfigItem struct {
	Container string `json:"-"`
	Format    string `json:"format"` // format for log, like json, csv, apache2, regexp, current only support json
	Type      string `json:"type"`   // file or stdout
	Path      string `json:"path"`   // path for log, if type is stdout, then path can be omitted
	Name      string `json:"name"`   // name for log
}

type ConfigItems []ConfigItem
