/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/base64"
	goerror "errors"
	"fmt"
	crdv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	nginx_manager "github.com/chulinx/netAgent/pkg/nginx"
	"github.com/chulinx/zlxGo/log"
	"github.com/chulinx/zlxGo/stringfile"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// VirtualServerReconciler reconciles a VirtualServer object
type VirtualServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtualServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *VirtualServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info("Reconciler netAgent", zap.Any("name", req.NamespacedName.Name), zap.Any("namespace", req.NamespacedName.Namespace))
	virtualServer := virtualserverv1alpha1.VirtualServer{}
	// manager Nginx Virtual Host Config lifecycle
	var err error
	if !goerror.Is(err, r.reconcileNginxVirtualServerConfig(ctx, virtualServer, req)) {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *VirtualServerReconciler) reconcileNginxVirtualServerConfig(ctx context.Context, vs virtualserverv1alpha1.VirtualServer, req ctrl.Request) error {
	// get service and secret
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error("get cluster config failed", zap.Any("err", err))
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("get clientSet failed", zap.Any("err", err))
	}
	//clientSet := ok8s.ClientSet("/Users/zhangxiang/.kube/dev-config")
	var (
		nameSpace, secretName string
	)

	secretSplit := strings.Split(vs.Spec.TlsSecret, "/")
	if len(secretSplit) == 2 {
		nameSpace = secretSplit[0]
		secretName = secretSplit[1]
	} else {
		nameSpace = req.Namespace
		secretName = vs.Spec.TlsSecret
	}

	secret, err := clientSet.CoreV1().Secrets(nameSpace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		log.Error("get service list failed", zap.Any("err", err))
		return err
	}
	tlsCrt, err1 := base64.StdEncoding.DecodeString(string(secret.Data["tls.crt"]))
	tlsKey, err2 := base64.StdEncoding.DecodeString(string(secret.Data["tls.key"]))
	if err1 != nil || err2 != nil {
		log.Error("decode tls.crt or tls.key failed", zap.Any("err1", err1), zap.Any("err2", err2))
		return goerror.New(fmt.Sprintf("decode tls.crt or tls.key failed %s,%s", err1, err2))
	}
	err3 := stringfile.RewriteFile(tlsCrt, fmt.Sprintf("/opt/%s-%s.crt", nameSpace, secretName))
	err4 := stringfile.RewriteFile(tlsKey, fmt.Sprintf("/opt/%s-%s.key", nameSpace, secretName))
	if err3 != nil || err4 != nil {
		log.Error("rewrite tls.crt or tls.key failed", zap.Any("err3", err3), zap.Any("err4", err4))
		return goerror.New(fmt.Sprintf("rewrite tls.crt or tls.key failed %s,%s", err3, err4))
	}
	err = r.Get(ctx, req.NamespacedName, &vs)
	nginxManager := nginx_manager.NewVirtualServerManager(vs)
	if err != nil {
		// delete virtual host config file
		if errors.IsNotFound(err) {
			log.Info("virtualServer not found", zap.Any("name", req.NamespacedName.Name))
			err := nginxManager.RemoveVirtualServerManager(req.NamespacedName)
			if err != nil {
				log.Error("remove virtual host config failed", zap.Any("err", err))
				return err
			}
			return nil
		}
		return err
	}
	// create or update virtual host config file and reload nginx
	err = nginxManager.CreateOrUpdate()
	log.Info("create or update virtual host config", zap.Any("name", vs.Name),
		zap.Any("server_name", vs.Spec.ServerName),
		zap.Any("port", vs.Spec.ListenPort))
	if err != nil {
		log.Error("create or update virtual host config failed", zap.Any("err", err))
		return err
	}
	vs.Status.State = "success"
	err = r.Update(ctx, &vs)
	// Todo
	// add or update delete service NodePort/LoadBalance/ClusterIP Port field
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1alpha1.VirtualServer{}).
		Complete(r)
}
