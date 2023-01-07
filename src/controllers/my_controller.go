package controllers

import (
	"context"
	"fmt"
	"github.com/Luukvdm/operator-template/src/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type MyResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myresource.luukvdm.github.com,resources=myresources/finalizers,verbs=update

func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx, "Request.Namespace", req.Namespace, "Request.Name", req.Name)

	// Get the resource
	myRes := &v1alpha1.MyResource{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, myRes)
	if err != nil {
		if errors.IsNotFound(err) {
			// If the resource somehow doesn't exist anymore, return without an error
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Do something with the resource
	// You can get the properties from the spec, like this:
	// myRes.Spec.FieldA
	l.Info(fmt.Sprintf("Reconciling for %s, with fielda: '%s' and fieldb: '%s'", req.Name, myRes.Spec.FieldA, myRes.Spec.FieldB))

	return ctrl.Result{}, nil
}

func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MyResource{}).
		Complete(r)
}
