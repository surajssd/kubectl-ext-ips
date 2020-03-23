package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", os.ExpandEnv("$HOME/.kube/config"), "kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the tabwriter to print the output in tabular format.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, '\t', 0)
	fmt.Fprintf(w, "NODE\tPRIVATE IP\tPUBLIC IP\n")

	for _, node := range nodes.Items {
		var internalIP string
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				internalIP = addr.Address
				break
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", node.Name, internalIP, node.Labels["lokomotive.alpha.kinvolk.io/public-ipv4"])
	}

	w.Flush()
}
