{{- define "spiringo.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "spiringo.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "spiringo.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
