package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"

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

func patchRoute() {

}

func banner() {
	log.Info("|==================================================|")
	log.Info("|========== Correcting Routes CTRL start ==========|")
	log.Info("|==================================================|")
}

func main() {
	// Instantiate loader for kubeconfig file.

	// enable signal trapping
	banner()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGTERM, // Termination Request
			syscall.SIGSEGV, // FullDerp
			syscall.SIGABRT, // Abnormal termination
			syscall.SIGILL,  // illegal instruction
			syscall.SIGFPE)  // floating point
		sig := <-c
		log.Fatalf("Signal (%v) Detected, Shutting Down", sig)
	}()

	for {

		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		// Get a rest.Config from the kubeconfig file.  This will be passed into all
		// the client objects we create.
		restconfig, err := kubeconfig.ClientConfig()
		if err != nil {
			panic(err)
		}

		routeclient, err := routev1client.NewForConfig(restconfig)
		if err != nil {
			panic(err)
		}

		// check with field selector https://stackoverflow.com/questions/41545123/how-to-get-pods-under-the-service-with-client-go-the-client-library-of-kubernete

		routes, err := routeclient.Routes("").List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		log.Info("--> Routes")
		for _, route := range routes.Items {
			/*
				log.WithFields(log.Fields{
					"Name":      route.Name,
					"Host":      route.Spec.Host,
					"Namespace": route.Namespace,
					//"GetSelfLink": route.SelfLink,
				}).Info("Routes")

			*/

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

		log.Info("sleeping 10 Second")
		time.Sleep(10 * time.Second)

	}

}
