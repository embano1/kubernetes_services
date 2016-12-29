package main

type ServiceList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Items      []Service
}

type Service struct {
	Metadata Metadata
	Spec     Spec
}

type Spec struct {
	Ports []struct {
		Protocol   string `json:"protocol"`
		Port       int    `json:"port"`
		TargetPort int    `json:"targetPort"`
		NodePort   int    `json:"nodePort"`
	}
}

type Metadata struct {
	Name string
}
