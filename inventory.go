package main

import (
	"maps"
	"slices"
)

type Inventory map[Owner]Ownership

type Owner string

type Ownership struct {
	Namespaces []string `json:"k8s_namespaces,omitempty"`
	Repos      []string `json:"repos,omitempty"`
}

func (i Inventory) Sorted(yield func(Owner, Ownership) bool) {
	for _, s := range slices.Sorted(maps.Keys(i)) {
		if !yield(s, i[s]) {
			return
		}
	}
}
