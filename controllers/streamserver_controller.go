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
	"github.com/chulinx/netAgent/pkg/nginx"
	"github.com/chulinx/zlxGo/log"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	streamserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
)

// StreamServerReconciler reconciles a StreamServer object
type StreamServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.chulinx,resources=streamservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.chulinx,resources=streamservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.chulinx,resources=streamservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the StreamServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *StreamServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info("Reconciler netAgent stream server", zap.Any("name", req.NamespacedName.Name), zap.Any("namespace", req.NamespacedName.Namespace))
	streamServer := streamserverv1alpha1.StreamServer{}
	// manager Nginx Stream server Config lifecycle
	var err error
	if err != r.reconcileNginxStreamServerConfig(ctx, streamServer, req) {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil

}

func (r *StreamServerReconciler) reconcileNginxStreamServerConfig(ctx context.Context, s streamserverv1alpha1.StreamServer, req ctrl.Request) error {
	err := r.Get(ctx, req.NamespacedName, &s)
	nginxManager := nginx.NewStreamServerManager(s)
	if err != nil {
		// delete virtual host config file
		if errors.IsNotFound(err) {
			log.Info("streamServer not found", zap.Any("name", req.NamespacedName.Name))
			err := nginxManager.RemoveVirtualServerManager(req.NamespacedName)
			if err != nil {
				log.Error("remove stream host config failed", zap.Any("err", err))
				return err
			}
			return nil
		}
		return err
	}
	// create or update virtual host config file and reload nginx
	err = nginxManager.CreateOrUpdate()
	log.Info("create or update stream host config", zap.Any("listPort", s.Name),
		zap.Any("service", s.Spec.Service),
		zap.Any("port", s.Spec.Port))
	if err != nil {
		log.Error("create or update stream host config failed", zap.Any("err", err))
		return err
	}
	s.Status.State = "success"
	err = r.Update(ctx, &s)
	// Todo
	// add or update delete service NodePort/LoadBalance/ClusterIP Port field
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StreamServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&streamserverv1alpha1.StreamServer{}).
		Complete(r)
}
