package main

import (
	"fmt"

	buildv1client "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	projectv1client "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	routev1client "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

func main() {
	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Determine the Namespace referenced by the current context in the
	// kubeconfig file.
	namespace, _, err := kubeconfig.Namespace()
	if err != nil {
		panic(err)
	}

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// Create an OpenShift build/v1 client.
	buildclient, err := buildv1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	projectclient, err := projectv1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// List all Builds in our current Namespace.
	projects, err := projectclient.Projects().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Projects in namespace %s:\n", namespace)
	for _, project := range projects.Items {
		fmt.Printf("  %s\n", project.Name)
		log.Info(project.Name)

	}

	routeclient, err := routev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	routes, err := routeclient.Routes("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Routes in namespace %s:\n", namespace)
	for _, route := range routes.Items {
		fmt.Printf("  %s\n", route.Name)
		fmt.Println(route.Spec.Host)
	}

	// List all Pods in our current Namespace.
	pods, err := coreclient.Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("  %s\n", pod.Name)

	}

	// List all Builds in our current Namespace.
	builds, err := buildclient.Builds(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Builds in namespace %s:\n", namespace)
	for _, build := range builds.Items {
		fmt.Printf("  %s\n", build.Name)
	}
}
