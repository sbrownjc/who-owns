package main

type Inventory map[Owner]Ownership

type Owner string

type Ownership struct {
	Namespaces []string `json:"k8s_namespaces,omitempty"`
	Repos      []string `json:"repos,omitempty"`
}
