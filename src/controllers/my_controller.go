package controllers

import (
	"context"
	"github.com/Luukvdm/operator-template/src/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type MyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresource,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresource/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresource/finalizers,verbs=update

func (reconciler *MyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconcile started for MyResource CRD")

	return ctrl.Result{}, nil
}

func (r *MyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MyResource{}).
		Complete(r)
}
