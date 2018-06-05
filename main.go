package main

import (
	"flag"
	"os"

	"github.com/kubernetes-sigs/kubebuilder/pkg/informer"
	logf "github.com/kubernetes-sigs/kubebuilder/pkg/log"
	corev1 "k8s.io/api/core/v1"
	kcache "k8s.io/client-go/tools/cache" // $cashmoney
)

func main() {
	flag.Parse()

	/*
		cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			log.Error(err, "could not initialize kubernetes client")
			os.Exit(1)
		}

		// NB: this should really only be done once for improved startup time

		discoClient := discovery.NewDiscoveryClientForConfigOrDie(cfg)
		groupReses, err := discovery.GetAPIGroupResources(discoClient)
		if err != nil {
			log.Error(err, "could not fetch API discovery information")
			os.Exit(1)
		}
		discoMapper := discovery.NewRESTMapper(groupReses, dynamic.VersionInterfaces)*/

	log := logf.Log.WithName("main")
	logf.SetLogger(logf.ZapLogger(true))
	log.Info("initializing...")

	cache := &informer.SelfPopulatingInformers{}
	podInformer, err := cache.InformerFor(&corev1.Pod{})
	if err != nil {
		log.Error(err, "could not initialize informer", "kind", "Pod")
		os.Exit(1)
	}

	go cache.Start(make(chan struct{}))

	podInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, obj interface{}) {
			log.Info("got pod update", "object name", obj.(*corev1.Pod).Name)
		},
	})

	select {}
}
