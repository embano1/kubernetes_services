// Calls the Kubernetes API <address> to get all services in namespace <default>
// Prints them out together with the NodePort

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {

	hostPtr := flag.String("s", "localhost", "API Host")
	portPtr := flag.String("p", "8001", "Port API is listing on")
	flag.Parse()

	var address string = "http://" + *hostPtr + ":" + *portPtr + "/api/v1/namespaces/default/services"
	fmt.Printf("Trying to connect to: %v\n", address)

	resp, err := http.Get(address)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Make sure %s is correct or start kubectl proxy on local machine\n", address)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Invalid status code", resp.Status)
		os.Exit(1)
	} else {
		fmt.Printf("Connected.\n\n")
	}

	defer resp.Body.Close()

	// Types declared in types.go
	var services ServiceList

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&services)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("Found services (excl. svc \"kubernetes\"):\n\n")
	}

	for _, service := range services.Items {

		// Don´t list the builtin "kubernetes" service
		if service.Metadata.Name == "kubernetes" {
			continue
		}

		fmt.Printf("Service Name: %s\n", service.Metadata.Name)

		for _, spec := range service.Spec.Ports {
			fmt.Printf("Nodeport: %v\nPort: %v\nTargetport: %v\n\n", spec.NodePort, spec.Port, spec.TargetPort)
		}

	}

}
