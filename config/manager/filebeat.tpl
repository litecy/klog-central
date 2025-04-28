# first version of filebeat tpl
{{range .}}
- type: log
  enabled: true
  paths:
      - {{ .HostPath }}{{ .ContainerPath }}
  scan_frequency: 10s
  fields_under_root: true

  {{if eq .Format "json"}}
  json.keys_under_root: true
  json.add_error_key: true
  json.overwrite_keys: true
  {{end}}
  fields:
      {{range $key, $value := .Meta}}
      {{ $key }}: {{ $value }}
      {{end}}
      {{range $key, $value := .AddFields}}
      {{ $key }}: {{ $value }}
      {{end}}
  processors:
    - drop_fields:
        fields: ["log"]
  {{if eq .Format "containerd_json"}}
    - dissect:
        tokenizer: "%{log_timestamp} %{std} %{capital-letter} %{content_json}"
        field: "message"
        target_prefix: ""
        overwrite_keys: true
    - drop_fields:
        fields: ["message"]
    - add_fields:
        when:
          equals:
            std: stderr
        target: log
        fields:
          level: "error"
    - add_fields:
        when:
          equals:
            std: stdout
        target: log
        fields:
          level: "info"
    - decode_json_fields:
        fields: ["content_json"]
        target: ""
        overwrite_keys: true
        ignore_failure: true
    - timestamp:
        field: log_timestamp
        layouts:
          - '2006-01-02T15:04:05Z'
          - '2006-01-02T15:04:05.999Z'
          - '2006-01-02T15:04:05.999-07:00'
        test:
          - '2019-06-22T16:33:51Z'
          - '2019-11-18T04:59:51.123Z'
          - '2020-08-03T07:10:20.123456+02:00'
          - '2023-02-02T07:49:45.576941217Z'
    - drop_fields:
        fields: ["log_timestamp","capital-letter", "std"]
        ignore_missing: true

  {{end}}
  {{if eq .Format "containerd"}}
    - dissect:
        tokenizer: "%{log_timestamp} %{std} %{capital-letter} %{content}"
        field: "message"
        target_prefix: ""
        overwrite_keys: true
    - drop_fields:
        fields: ["message"]
    - add_fields:
        when:
          equals:
            std: stderr
        target: log
        fields:
          level: "error"
    - add_fields:
        when:
          equals:
            std: stdout
        target: log
        fields:
          level: "info"

    - timestamp:
        field: log_timestamp
        layouts:
          - '2006-01-02T15:04:05Z'
          - '2006-01-02T15:04:05.999Z'
          - '2006-01-02T15:04:05.999-07:00'
        test:
          - '2019-06-22T16:33:51Z'
          - '2019-11-18T04:59:51.123Z'
          - '2020-08-03T07:10:20.123456+02:00'
          - '2023-02-02T07:49:45.576941217Z'
    - drop_fields:
        fields: ["log_timestamp","capital-letter", "std"]
        ignore_missing: true

  {{end}}
    - rename:
        fields:
         - from: "content"
           to: "message"
         - from: "level"
           to: "log.level"
        ignore_missing: true
        fail_on_error: false
    - drop_fields:
        fields: ["input", "host", "ecs","agent", "host.name"]
        ignore_missing: true
  tail_files: false
  close_inactive: 2m
  close_eof: false
  close_removed: true
  clean_removed: true
  close_renamed: false

{{end}}
