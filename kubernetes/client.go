package kubernetes

import (
	"fmt"

	"github.com/aporeto-inc/kubepox"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// Client is the Trireme representation of the Client.
type Client struct {
	kubeClient *client.Client
	namespace  string
	localNode  string
}

// NewClient Generate and initialize a Trireme Client object
func NewClient(kubeconfig string, namespace string) (*Client, error) {
	Client := &Client{}
	err := Client.InitKubernetesClient(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Coultn't initialize Kubernetes Client: %v", err)
	}
	return Client, nil
}

// InitKubernetesClient Initialize the Kubernetes client
func (k *Client) InitKubernetesClient(kubeconfig string) error {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("Error opening Kubeconfig: %v", err)
	}

	myClient, err := client.New(config)
	if err != nil {
		return fmt.Errorf("Error creating REST Kube Client: %v", err)
	}
	k.kubeClient = myClient
	return nil
}

// PodRules return the list of all the IngressRules that apply to the pod.
func (k *Client) PodRules(podName string, namespace string) (*[]extensions.NetworkPolicyIngressRule, error) {
	// Step1: Get all the rules associated with this Pod.
	targetPod, err := k.kubeClient.Pods(namespace).Get(podName)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get pod %v from Kubernetes API: %v", podName, err)
	}

	allPolicies, err := k.kubeClient.Extensions().NetworkPolicies(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list all the NetworkPolicies from Kubernetes API: %v ", err)
	}

	allRules, err := kubepox.ListIngressRulesPerPod(targetPod, allPolicies)
	if err != nil {
		return nil, fmt.Errorf("Couldn't process the list of rules for pod %v : %v", podName, err)
	}
	return allRules, nil
}

// PodLabels returns the list of all label associated with a pod.
func (k *Client) PodLabels(podName string, namespace string) (map[string]string, error) {
	targetPod, err := k.kubeClient.Pods(namespace).Get(podName)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes labels for pod %v : %v ", podName, err)
	}
	return targetPod.GetLabels(), nil
}

// LocalPods return a PodList with all the pods scheduled on the local node
func (k *Client) LocalPods(namespace string) (*api.PodList, error) {
	// TODO: Generate ListOptions to match on the local node
	return k.kubeClient.Pods(namespace).List(api.ListOptions{})
}

// AllNamespaces return a list of all existing namespaces
func (k *Client) AllNamespaces() (*api.NamespaceList, error) {
	return k.kubeClient.Namespaces().List(api.ListOptions{})
}
