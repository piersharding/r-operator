package types

import (
	"github.com/appscode/go/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RContext is the set of parameters to configures this instance
type RContext struct {
	Ingress      string
	Daemon       bool
	Namespace    string
	Name         string
	ServiceType  string
	Port         int
	Replicas     int
	Image        string
	Repository   string
	Tag          string
	PullSecrets  interface{}
	PullPolicy   string
	NodeSelector interface{}
	Affinity     interface{}
	Tolerations  interface{}
	Resources    interface{}
	VolumeMounts interface{}
	Volumes      interface{}
	Env          interface{}
	Password     string
}

// Controller struct root object from payload
type Controller struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ControllerSpec   `json:"spec"`
	Status            ControllerStatus `json:"status"`
}

// ControllerSpec struct defines Spec expected from Kind: Rapp
type ControllerSpec struct {
	Ingress      string      `json:"ingress"`
	Replicas     int         `json:"replicas"`
	Daemon       bool        `json:"daemon"`
	Password     string      `json:"password"`
	Image        string      `json:"image"`
	PullPolicy   string      `json:"imagePullPolicy"`
	PullSecrets  interface{} `json:"imagePullSecrets"`
	NodeSelector interface{} `json:"nodeSelector"`
	Affinity     interface{} `json:"affinity"`
	Tolerations  interface{} `json:"tolerations"`
	Resources    interface{} `json:"resources"`
	VolumeMounts interface{} `json:"volumeMounts"`
	Volumes      interface{} `json:"volumes"`
	Env          interface{} `json:"env"`
}

// ControllerStatus struct response status
type ControllerStatus struct {
	Replicas  int    `json:"replicas"`
	Succeeded int    `json:"succeeded"`
	State     string `json:"state"`
	Resources string `json:"resources"`
}

// SyncRequest struct root request object
type SyncRequest struct {
	Parent   Controller          `json:"parent"`
	Children SyncRequestChildren `json:"children"`
}

// SyncRequestChildren struct children objects returned
type SyncRequestChildren struct {
	// Pods         map[string]*v1.Pod             `json:"Pod.v1"`
	// StatefulSets map[string]*appsv1.StatefulSet `json:"StatefulSet.apps/v1"`
	Deployments map[string]*appsv1.Deployment `json:"Deployment.apps/v1"`
}

// SyncResponse struct root response object
type SyncResponse struct {
	Status   ControllerStatus `json:"status"`
	Children []runtime.Object `json:"children"`
}

// SetConfig setup the configuration
func SetConfig(request *SyncRequest) RContext {

	context := RContext{
		Ingress:      request.Parent.Spec.Ingress,
		Daemon:       request.Parent.Spec.Daemon,
		Namespace:    request.Parent.Namespace,
		Name:         request.Parent.Name,
		ServiceType:  "ClusterIP",
		Port:         8786,
		Replicas:     request.Parent.Spec.Replicas,
		Image:        request.Parent.Spec.Image,
		PullSecrets:  request.Parent.Spec.PullSecrets,
		PullPolicy:   request.Parent.Spec.PullPolicy,
		NodeSelector: request.Parent.Spec.NodeSelector,
		Affinity:     request.Parent.Spec.Affinity,
		Tolerations:  request.Parent.Spec.Tolerations,
		Resources:    request.Parent.Spec.Resources,
		VolumeMounts: request.Parent.Spec.VolumeMounts,
		Volumes:      request.Parent.Spec.Volumes,
		Env:          request.Parent.Spec.Env,
		Password:     request.Parent.Spec.Password}

	if context.Password == "" {
		context.Password = "password"
	}
	log.Debugf("context: %+v", context)
	return context
}
