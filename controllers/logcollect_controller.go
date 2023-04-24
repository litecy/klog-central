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
	"errors"
	"github.com/go-logr/logr"
	"github.com/litecy/klog-central/pkg/filter"
	"github.com/litecy/klog-central/pkg/processor"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// LogCollectReconciler reconciles a LogCollect object
type LogCollectReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Processor *processor.KConfig
}

//+kubebuilder:rbac:groups=klog.vibly.vip,resources=logcollects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=klog.vibly.vip,resources=logcollects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=klog.vibly.vip,resources=logcollects/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the LogCollect object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *LogCollectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var pod v1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	kcfg, err := filter.CheckKLCConfig(ctx, pod)
	if err != nil {
		// some errors here
		return ctrl.Result{}, nil
	}

	if kcfg == nil || len(*kcfg) == 0 {
		// no suitable log config find
		return ctrl.Result{}, nil
	}

	err = r.Processor.Process(ctx, kcfg, &pod)
	if err != nil {
		logger.Info("handle pod with logs change failed", "pod", req.NamespacedName, "annotations", pod.ObjectMeta.Annotations, "kcfg", kcfg, "err", err)
		return reconcile.Result{
			Requeue: true,
		}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogCollectReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if r.Processor == nil {
		return errors.New("kconfig is nil, please setup kconfig first")
	}

	err := r.Processor.Setup()
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}
