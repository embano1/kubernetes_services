// --- Type definitions ---
package main

// --- API service specs ---
type servicelist struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Items      []service
}

type service struct {
	Metadata struct {
		Name string
	} `json:"metadata"`
	Spec spec
}

type spec struct {
	Ports []struct {
		Protocol   string `json:"protocol"`
		Port       int    `json:"port"`
		TargetPort int    `json:"targetPort"`
		NodePort   int    `json:"nodePort"`
	}
	Type string `json:"type"`
}

// --- API node specs ---
type nodelist struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Items      []node
}

type node struct {
	Status struct {
		Addresses []struct {
			Type    string `json:"type"`
			Address string `json:"address"`
		}
	}
}
