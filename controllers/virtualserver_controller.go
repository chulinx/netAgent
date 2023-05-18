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
	crdv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	nginx_manager "github.com/chulinx/netAgent/pkg/nginx"
	"github.com/chulinx/zlxGo/log"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VirtualServerReconciler reconciles a VirtualServer object
type VirtualServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	//NginxManager *nginx_manager.Manager
}

//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.chulinx,resources=virtualservers/finalizers,verbs=update

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
	if err != r.reconcileNginxVirtualServerConfig(ctx, virtualServer, req) {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *VirtualServerReconciler) reconcileNginxVirtualServerConfig(ctx context.Context, vs virtualserverv1alpha1.VirtualServer, req ctrl.Request) error {
	err := r.Get(ctx, req.NamespacedName, &vs)
	nginxManager := nginx_manager.NewManager(vs)
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
	err = nginxManager.CreateOrUpdateVirtualServerManager()
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