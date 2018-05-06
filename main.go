package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	buildv1client "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	projectv1client "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	routev1client "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	"regexp"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	ClusterEnvironment := strings.ToLower(os.Getenv("ClusterEnvironment"))

	if ClusterEnvironment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func findRoute(route string, f string, r string) (v string, err error) {

	//match, _ := regexp.MatchString("p([a-z]+)ch", route.Spec.Host)

	match, _ := regexp.MatchString(f, route)
	if match {
		va := strings.Replace(route, f, r, 1)
		return va, nil
	}
	return "", errors.New("Empty")

}

func main() {
	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Determine the Namespace referenced by the current context in the
	// kubeconfig file.
	/*
		namespace, _, err := kubeconfig.Namespace()
		if err != nil {
			panic(err)
		}
	*/

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

	log.Info("--> Projects")
	for _, project := range projects.Items {
		log.WithFields(log.Fields{
			"Name": project.Name,
		}).Info("Projects")

	}

	routeclient, err := routev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	routes, err := routeclient.Routes("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	log.Info("--> Routes")
	for _, route := range routes.Items {

		log.WithFields(log.Fields{
			"Name":      route.Name,
			"Host":      route.Spec.Host,
			"Namespace": route.Namespace,
			//"GetSelfLink": route.SelfLink,
		}).Info("Routes")

		//match, _ := regexp.MatchString("p([a-z]+)ch", route.Spec.Host)

		va, err := findRoute(route.Spec.Host, "bit", "big")
		if err != nil {
			//panic(err)
		}
		if va != "" {
			log.Warn("--> correcting Route: ", va)
			route.Spec.Host = va

			log.WithFields(log.Fields{
				"Name":      route.Name,
				"Host":      route.Spec.Host,
				"Namespace": route.Namespace,
			}).Info("Routes")

			//c.Client.Build().Builds(namespace).Update(build)
			//routes, err := routeclient.Routes(route.Namespace).Patch(route.Name, api.JSONPatchType, route)

			/*
			   patchBytes := []byte(`{"spec":{"paused":true,"replicas":0,"revisionHistoryLimit":0}}`)
			   	return reaper.appsClient.Apps().DeploymentConfigs(namespace).Patch(name, types.StrategicMergePatchType, patchBytes)

			*/
			log.Warn("..........  patch route !!!!!!!!")
			patchBytes := []byte(fmt.Sprintf("{\"spec\":{\"host\":\"%s\"}}", route.Spec.Host))

			//fmt.Printf("{\"spec\":{\"host\":\"%v\"}}", route.Spec.Host)
			//fmt.Print(patchBytes)
			log.Info(fmt.Sprintf("%s", patchBytes))

			_, err := routeclient.Routes(route.Namespace).Patch(route.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				panic(err)
			} else {
				log.Info("--> patch successfull :-) !!!")
			}

		}

	}

	// List all Pods in our current Namespace.
	pods, err := coreclient.Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	log.Info("--> Pods")

	podLogger := log.WithFields(log.Fields{
		"Get": "pods",
	})

	podLogger.Info("I'll be logged with common and other field")
	for _, pod := range pods.Items {

		log.WithFields(log.Fields{
			"Name":      pod.Name,
			"Namespace": pod.Namespace,
		}).Info("Pods")

	}

	// List all Builds in our current Namespace.
	builds, err := buildclient.Builds("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	log.Info("--> Builds")
	for _, build := range builds.Items {
		log.WithFields(log.Fields{
			"Name":      build.Name,
			"Namespace": build.Namespace,
		}).Info("Builds")
	}
}
