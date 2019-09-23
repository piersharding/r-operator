package models

import (
	"encoding/json"

	"github.com/appscode/go/log"
	dtypes "github.com/piersharding/r-operator/types"
	"github.com/piersharding/r-operator/utils"
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

// RappIngress generates the Ingress description for
// the R cluster
func RappIngress(context dtypes.RContext) (*v1beta1.Ingress, error) {

	const rappIngress = `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: rapp-{{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    app.kubernetes.io/name: rapp
    app.kubernetes.io/instance: "{{ .Name }}"
    app.kubernetes.io/managed-by: MetaController
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/x-forwarded-prefix: "true"
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
    nginx.ingress.kubernetes.io/affinity: "cookie"
    nginx.ingress.kubernetes.io/session-cookie-name: "rappsticky"
    nginx.ingress.kubernetes.io/session-cookie-expires: "172800"
    nginx.ingress.kubernetes.io/session-cookie-max-age: "172800"
spec:
  rules:
  - host: {{ .Ingress }}
    http:
      paths:
      - path: /
        backend:
          serviceName:  rapp-{{ .Name }}
          servicePort: 8080
`
	result, err := utils.ApplyTemplate(rappIngress, context)
	if err != nil {
		log.Debugf("ApplyTemplate Error: %+v\n", err)
		return nil, err
	}
	ingress := &v1beta1.Ingress{}
	if err := json.Unmarshal([]byte(result), ingress); err != nil {
		return nil, err
	}
	return ingress, err
}
