package models

import (
	"encoding/json"

	"github.com/appscode/go/log"
	dtypes "github.com/piersharding/r-operator/types"
	"github.com/piersharding/r-operator/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

// RappService generates the Service description for
// the Rapp Notebook server
func RappService(context dtypes.RContext) (*v1.Service, error) {

	const rappService = `
apiVersion: v1
kind: Service
metadata:
  name: rapp-{{ .Name }}
  labels:
    app.kubernetes.io/name: rapp
    app.kubernetes.io/instance: "{{ .Name }}"
    app.kubernetes.io/managed-by: MetaController
spec:
  selector:
    app.kubernetes.io/name:  rapp
    app.kubernetes.io/instance: "{{ .Name }}"
  type: {{ .ServiceType }}
  ports:
  - name: http
    port: 8080
    targetPort: http
    protocol: TCP
`
	result, err := utils.ApplyTemplate(rappService, context)
	if err != nil {
		log.Debugf("ApplyTemplate Error: %+v\n", err)
		return nil, err
	}
	service := &v1.Service{}
	if err := json.Unmarshal([]byte(result), service); err != nil {
		return nil, err
	}
	return service, err
}

// RappDeployment generates the Deployment description for
// the Jupyter Notebook
func RappDeployment(context dtypes.RContext) (*appsv1.Deployment, error) {

	const rappDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rapp-{{ .Name }}
  labels:
    app.kubernetes.io/name: rapp
    app.kubernetes.io/instance: "{{ .Name }}"
    app.kubernetes.io/managed-by: MetaController
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: rapp
      app.kubernetes.io/instance: "{{ .Name }}"
  replicas: {{ .Replicas }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: rapp
        app.kubernetes.io/instance: "{{ .Name }}"
        app.kubernetes.io/managed-by: MetaController
    spec:
    {{- with .PullSecrets }}
      imagePullSecrets:
      {{range $val := .}}
      - name: {{ $val.name }}
      {{end}}
      {{- end }}      
      containers:
      - name: rapp
        image: "{{ .Image }}"
        imagePullPolicy: {{ .PullPolicy }}
        command:
          - /start-rapp.sh
        env:
          - name: HOST
            value: "0.0.0.0"
          - name: PORT
            value: "8080"
{{- with .Env }}
{{ toYaml . | indent 10 }}
{{- end }}
        ports:
        - name: http
          containerPort: 8080
        volumeMounts:
        - mountPath: /start-rapp.sh
          subPath: start-rapp.sh
          name: rapp-script
        - mountPath: /var/tmp
          readOnly: false
          name: localdir
{{- with .VolumeMounts }}
{{ toYaml . | indent 8 }}
{{- end }}
        readinessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 10
          timeoutSeconds: 10
          periodSeconds: 20
          failureThreshold: 3
      volumes:
      - configMap:
          name: rapp-configs-{{ .Name }}
          defaultMode: 0777
        name: rapp-script
      - hostPath:
          path: /var/tmp
          type: DirectoryOrCreate
        name: localdir

{{- with .Volumes }}
{{ toYaml . | indent 6 }}
{{- end }}
{{- with .NodeSelector }}
      nodeSelector:
{{ toYaml . | indent 8 }}
{{- end }}
{{- with .Affinity }}
      affinity:
{{ toYaml . | indent 8 }}
{{- end }}
{{- with .Tolerations }}
      tolerations:
{{ toYaml . | indent 8 }}
{{- end }}

`
	if context.Daemon {
		log.Infof("Adding Daemon affinity rules")
		if context.Affinity == nil {
			log.Debugf("context.Affinity.podAntiAffinity does not exist")
			context.Affinity = map[string]interface{}{}
		}
		if _, ok := context.Affinity.(map[string]interface{})["podAntiAffinity"]; !ok {
			log.Debugf("context.Affinity.podAntiAffinity does not exist")
			context.Affinity.(map[string]interface{})["podAntiAffinity"] = map[string]interface{}{}
		}
		cAp := context.Affinity.(map[string]interface{})["podAntiAffinity"]

		if _, ok := cAp.(map[string]interface{})["requiredDuringSchedulingIgnoredDuringExecution"]; !ok {
			log.Debugf("context.Affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution does not exist")
			cAp.(map[string]interface{})["requiredDuringSchedulingIgnoredDuringExecution"] = []interface{}{}
		}
		cAp.(map[string]interface{})["requiredDuringSchedulingIgnoredDuringExecution"] =
			append(cAp.(map[string]interface{})["requiredDuringSchedulingIgnoredDuringExecution"].([]interface{}),
				map[string]interface{}{
					"labelSelector": map[string][]map[string]interface{}{
						"matchExpressions": []map[string]interface{}{
							map[string]interface{}{
								"key":      "app.kubernetes.io/instance",
								"operator": "In",
								"values":   []string{context.Name}}}},
					"topologyKey": "kubernetes.io/hostname"})
	}

	result, err := utils.ApplyTemplate(rappDeployment, context)
	if err != nil {
		log.Debugf("ApplyTemplate Error: %+v\n", err)
		return nil, err
	}

	deployment := &appsv1.Deployment{}
	if err := json.Unmarshal([]byte(result), deployment); err != nil {
		return nil, err
	}
	return deployment, err
}
