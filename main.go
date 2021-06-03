package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc := dynamic.NewForConfigOrDie(config)
	gvrNode := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	nodes, err := dc.Resource(gvrNode).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, obj := range nodes.Items {
		fmt.Printf("%+v\n", obj.GetName())
	}

	mc := metadata.NewForConfigOrDie(config)
	nodemeta, err := mc.Resource(gvrNode).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range nodemeta.Items {
		fmt.Println(n.GetName())
	}
}
