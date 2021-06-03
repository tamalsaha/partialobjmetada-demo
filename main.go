package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc := dynamic.NewForConfigOrDie(config)
	nodes, err := dc.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, obj := range nodes.Items {
		s, err := status.Compute(&obj)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", *s)
	}

}
