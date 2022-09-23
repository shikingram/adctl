version: '3.2'
{{- if or .Values.mysql.enabled  .Values.mysql.adminer.enabled }}
services:
{{- end}}
{{- if .Values.mysql.enabled }}
  mysql:
    image: {{ .Values.mysql.image }}:{{ .Values.mysql.tag }}
    restart: always
    volumes:
    {{- range .Values.mysql.volumes }}
    - {{ . }}
    {{- end }}
    {{- if .Values.mysql.storage.mysql_config}}
    - {{ .Values.mysql.storage.mysql_config }}:/var/lib/mysql
    {{- end}}
    ports: {{ if eq (len .Values.mysql.export_port) 0 }}[]{{- end}}
    {{- range .Values.mysql.export_port }}
    - {{ .node_port }}:{{ .port }}
    {{- end }}
    environment:
    - MYSQL_ROOT_PASSWORD={{ .Values.mysql.MYSQL_ROOT_PASSWORD }}
    networks:
      - {{ .Release.Name }}
{{- end }}
{{- if .Values.mysql.adminer.enabled }}
  adminer:
    image: {{ .Values.mysql.adminer.image }}:{{.Values.mysql.adminer.tag }}
    restart: always
    ports: {{ if eq (len .Values.mysql.adminer.export_port) 0 }}[]{{- end}}
    {{- range .Values.mysql.adminer.export_port }}
    - {{ .node_port }}:{{ .port }}
    {{- end }}
    networks:
      - {{ .Release.Name }}
{{- end}}
networks:
  {{ .Release.Name }}:
    external: true
