package main

import (
	"context"
	"flag"
	"os"

	"github.com/kubernetes-sigs/kubebuilder/pkg/client"
	"github.com/kubernetes-sigs/kubebuilder/pkg/ctrl"
	"github.com/kubernetes-sigs/kubebuilder/pkg/ctrl/eventhandler"
	"github.com/kubernetes-sigs/kubebuilder/pkg/ctrl/reconcile"
	"github.com/kubernetes-sigs/kubebuilder/pkg/ctrl/source"
	logf "github.com/kubernetes-sigs/kubebuilder/pkg/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var log = logf.Log.WithName("main")

func main() {
	// Init main
	flag.Parse()
	logf.SetLogger(logf.ZapLogger(true))

	// Init Controller
	c := &ctrl.Controller{
		Reconcile: reconcile.ReconcileFunc(func(request reconcile.ReconcileRequest) (reconcile.ReconcileResult, error) {
			log.Info("Got Reconcile", "Request", request)
			return reconcile.ReconcileResult{}, nil
		}),
	}
	cm := &ctrl.ControllerManager{}
	cm.AddController(c, func() {
		c.Watch(&source.KindSource{Type: &corev1.Endpoints{}}, &eventhandler.EnqueueHandler{})
		c.FieldIndexes.IndexField(&corev1.Endpoints{}, "synthetic.targets", func(obj runtime.Object) []string {
			ep := obj.(*corev1.Endpoints)
			var res []string
			for _, subset := range ep.Subsets {
				for _, addr := range subset.Addresses {
					res = append(res, addr.TargetRef.Name)
				}
			}
			return res
		})
		c.Watch(&source.KindSource{Type: &corev1.Pod{}}, &eventhandler.EnqueueMappedHandler{
			ToRequests: eventhandler.ToRequestsFunc(func(evt eventhandler.ToRequestArg) []reconcile.ReconcileRequest {
				pods := &corev1.PodList{}
				c.Client.List(context.TODO(), client.MatchingField("synthetic.targets", evt.Meta.GetName()), pods)
				return nil
			}),
		})
	})

	// Start main
	if err := cm.Start(make(chan struct{})); err != nil {
		log.Error(err, "Failed to Start ControllerManager")
		os.Exit(1)
	}
}
