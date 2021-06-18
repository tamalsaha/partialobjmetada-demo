package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"kmodules.xyz/client-go/discovery"
	"log"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main__() {
	vp, err := semver.NewVersion("v1.2.3-alpha.0+buil9")
	if err != nil {
		panic(err)
	}
	v := *vp
	v, err = v.SetPrerelease("")
	if err != nil {
		panic(err)
	}
	v, err = v.SetMetadata("")
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Original())

	v.IncPatch()
}

type Object interface {
	metav1.Object
	runtime.Object
}


func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kc := kubernetes.NewForConfigOrDie(config)

	m := discovery.NewRestMapper(kc.Discovery())
	rsm := discovery.NewResourceMapper(m)
	rsm.Reset()

	dc := dynamic.NewForConfigOrDie(config)
	gvrNode := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	nodes, err := dc.Resource(gvrNode).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			fmt.Println(err.Error())
		}
		panic(err)
	}

	for _, obj := range nodes.Items {
		copytest(&obj)
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

func copytest(obj client.Object) {
	o2 := obj.DeepCopyObject().(client.Object)
	o2.SetManagedFields(nil)

	dobj, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(dobj))
	fmt.Println("--------------------------------")
	d2, _ := json.MarshalIndent(o2, "", "  ")
	fmt.Println(string(d2))
}
