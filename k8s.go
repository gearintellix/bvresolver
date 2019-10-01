package bvresolver

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8sClient struct {
	stop      chan bool
	namespace string
	clientSet *kubernetes.Clientset
}

func newK8sClient(namespace string) (cln k8sClient, err error) {
	cln = k8sClient{
		namespace: namespace,
	}

	conf, err := rest.InClusterConfig()
	if err != nil {
		err = errors.New("failed to getting cluster config")
		return cln, err
	}

	cln.clientSet, err = kubernetes.NewForConfig(conf)
	if err != nil {
		err = errors.New("")
	}

	return cln, err
}

func (ox *k8sClient) watchEndpoints(service string, stop chan bool) error {
	// a, b := ox.clientSet.CoreV1().Endpoints("").List(metav1.ListOptions{})

	api := ox.clientSet.CoreV1()
	wtc, err := api.Endpoints(ox.namespace).Watch(metav1.ListOptions{
		LabelSelector: service,
	})
	if err != nil {
		return err
	}

	go func() {
		<-stop
		wtc.Stop()
	}()

	for val := range wtc.ResultChan() {
		enp, ok := val.Object.(*v1.Endpoints)
		if !ok {
			return errors.New("unexpected type")
		}

		switch val.Type {
		case watch.Added:
			fmt.Printf("[ADD] %+x", enp.Subsets)

		case watch.Deleted:
			fmt.Printf("[DEL] %+x", enp.Subsets)

		case watch.Modified:
			fmt.Printf("[MOD] %+x", enp.Subsets)

		case watch.Error:
			return errors.New("error has occured")
		}
	}

	return nil
}
