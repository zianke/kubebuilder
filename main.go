package main

import (
	"log"
	"os"
	"flag"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	kcache "k8s.io/client-go/tools/cache" // $cashmoney
	corev1 "k8s.io/api/core/v1"
	"github.com/kubernetes-sigs/kubebuilder/pkg/informer"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatalf("shit: %v", err)
	}

	discoClient := discovery.NewDiscoveryClientForConfigOrDie(cfg)
	groupReses, err := discovery.GetAPIGroupResources(discoClient)
	if err != nil {
		log.Fatalf("gaaaaah: %v", err)
	}
	discoMapper := discovery.NewRESTMapper(groupReses, dynamic.VersionInterfaces)

	cache := informer.NewInformerCache(discoMapper, cfg, scheme.Scheme)
	podInformer, err := cache.InformerFor(&corev1.Pod{})
	if err != nil {
		log.Fatalf("darnit: %v", err)
	}

	go cache.Start(make(chan struct{}))

	podInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, obj interface{}) {
			log.Printf("haha: %v", obj.(*corev1.Pod).Name)
		},
	})

	select{}
}
