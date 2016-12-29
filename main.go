// Queries Kubernetes API <address> to get all services in namespace <default>
// Prints them as a table with the NodeIP:NodePort

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// --- Define some global variables ---
var (
	// Flags
	proxymode bool
	allsvc    bool
	host      string
	namespace string
	port      int

	// Types declared in types.go
	services servicelist
	nodes    nodelist

	// Write output table to
	tabledata [][]string

	// Various
	apiservice string
	apinode    string
)

// --- Error handler func ---
func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// --- Command line parsing and help ---
func flags() {
	// TODO flag.BoolVar(&proxymode, "c", false, "Use kube.config instead of proxy (default false)")
	flag.BoolVar(&allsvc, "a", false, "List all services incl. those type != NodePort (default false)")
	flag.StringVar(&host, "h", "localhost", "Kubernetes API Server")
	flag.IntVar(&port, "p", 8001, "Port API Server is listing on")
	flag.StringVar(&namespace, "n", "default", "Namespace to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s queries the Kubernetes API for services exposed on NodePorts.", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nUsage of %s: \n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
}

// --- Get a node ip we can connect to NodePort(s) ---
func getvalidnode(nodes *nodelist) (ip string, err error) {
	// TODO node might not be ready/ cordoned, etc. check for readiness
Loop:
	for _, node := range nodes.Items {
		for _, ipaddr := range node.Status.Addresses {

			if ipaddr.Type == "InternalIP" {
				ip = ipaddr.Address
				break Loop
			}
		}
	}
	return
}

// --- Generate a nicer output ---
func gentable(tabledata [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Port", "Target Port", "NodePort", "Connect"})

	for _, v := range tabledata {
		table.Append(v)
	}
	table.Render()
	return
}

func main() {

	// Command line parsing
	flags()

	// Generate server strings
	apiservice = fmt.Sprintf("http://%s:%v/api/v1/namespaces/%s/services", host, port, namespace)
	apinode = fmt.Sprintf("http://%s:%v/api/v1/nodes", host, port)

	// Connect and read http GET
	fmt.Printf("Trying to connect to: http://%v:%v (namespace: %v)\n", host, port, namespace)
	respservice, err := http.Get(apiservice)
	must(err)
	if respservice.StatusCode != 200 {
		fmt.Println("Invalid status code", respservice.Status)
		os.Exit(1)
	} else {
		fmt.Printf("Connected.\n\n")
	}

	respnode, err := http.Get(apinode)
	must(err)
	if respnode.StatusCode != 200 {
		fmt.Println("Invalid status code", respnode.Status)
		os.Exit(1)
	}

	defer respservice.Body.Close()
	defer respnode.Body.Close()

	// Decode json streams from GET
	decservice := json.NewDecoder(respservice.Body)
	err = decservice.Decode(&services)
	must(err)

	decnode := json.NewDecoder(respnode.Body)
	err = decnode.Decode(&nodes)
	must(err)

	// Get a node ip so we can connect to the NodePort from outside
	nodeip, err := getvalidnode(&nodes)
	must(err)

	// Build table data
	for _, service := range services.Items {
		if service.Spec.Type != "NodePort" && !allsvc {
			continue
		}

		for _, spec := range service.Spec.Ports {
			tabledata = append(tabledata, []string{service.Metadata.Name, strconv.Itoa(spec.Port), strconv.Itoa(spec.TargetPort), strconv.Itoa(spec.NodePort), "http://" + nodeip + ":" + strconv.Itoa(spec.NodePort)})
		}

	}

	// Print output as table
	gentable(tabledata)
	os.Exit(0)

}
