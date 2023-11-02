package machineset

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	"github.com/Azure/ARO-RP/pkg/operator/controllers/base"
)

const (
	controllerName    = "MachineSet"
	controllerEnabled = "aro.machineset.enabled"
)

type Reconciler struct {
	base.AROController
}

// MachineSet reconciler watches MachineSet objects for changes, evaluates total worker replica count, and reverts changes if needed.
func NewReconciler(log *logrus.Entry, client client.Client) *Reconciler {
	r := &Reconciler{
		AROController: base.AROController{
			Log:         log.WithField("controller", controllerName),
			Client:      client,
			Name:        controllerName,
			EnabledFlag: controllerEnabled,
		},
	}
	r.Reconciler = r
	return r
}

func (r *Reconciler) ReconcileEnabled(ctx context.Context, request ctrl.Request, instance *arov1alpha1.Cluster) (ctrl.Result, error) {
	var err error

	modifiedMachineset := &machinev1beta1.MachineSet{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: request.Name, Namespace: machineSetsNamespace}, modifiedMachineset)
	if err != nil {
		r.Log.Error(err)
		r.SetDegraded(ctx, err)

		return reconcile.Result{}, err
	}

	machinesets := &machinev1beta1.MachineSetList{}
	selector, _ := labels.Parse("machine.openshift.io/cluster-api-machine-role=worker")
	err = r.Client.List(ctx, machinesets, &client.ListOptions{
		Namespace:     machineSetsNamespace,
		LabelSelector: selector,
	})
	if err != nil {
		r.Log.Error(err)
		r.SetDegraded(ctx, err)

		return reconcile.Result{}, err
	}

	// Count amount of total current worker replicas
	replicaCount := 0
	for _, machineset := range machinesets.Items {
		// If there are any custom machinesets in the list, bail and don't requeue
		if !strings.Contains(machineset.Name, instance.Spec.InfraID) {
			r.ClearDegraded(ctx)

			return reconcile.Result{}, nil
		}
		if machineset.Spec.Replicas != nil {
			replicaCount += int(*machineset.Spec.Replicas)
		}
	}

	if replicaCount < minSupportedReplicas {
		r.Log.Infof("Found less than %v worker replicas. The MachineSet controller will attempt scaling.", minSupportedReplicas)
		// Add replicas to the object, and call Update
		modifiedMachineset.Spec.Replicas = to.Int32Ptr(int32(minSupportedReplicas-replicaCount) + *modifiedMachineset.Spec.Replicas)
		err := r.Client.Update(ctx, modifiedMachineset)
		if err != nil {
			r.Log.Error(err)
			r.SetDegraded(ctx, err)

			return reconcile.Result{}, err
		}
	}

	r.ClearConditions(ctx)
	return reconcile.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	machineSetPredicate := predicate.NewPredicateFuncs(func(o client.Object) bool {
		role := o.GetLabels()["machine.openshift.io/cluster-api-machine-role"]
		return strings.EqualFold("worker", role)
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1beta1.MachineSet{}, builder.WithPredicates(machineSetPredicate)).
		Named(r.GetName()).
		Complete(r)
}
