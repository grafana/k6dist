package registry

import (
	"github.com/grafana/k6foundry"
)

// Module representa a buildable k6 module.
type Module struct {
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
	Cgo     bool   `json:"cgo,omitempty"`
}

// Modules represents list of buildable modules.
type Modules []Module

// ToFoundry converts module list to k6foundry builder parameters.
func (mods Modules) ToFoundry() (string, []k6foundry.Module) {
	var k6Version string

	fmods := make([]k6foundry.Module, 0, len(mods)-1)

	for _, mod := range mods {
		if mod.Path == k6ModulePath {
			k6Version = mod.Version
		} else {
			fmods = append(fmods, k6foundry.Module{Path: mod.Path, Version: mod.Version})
		}
	}

	return k6Version, fmods
}

// Cgo returns true if at least one module requires cgo.
func (mods Modules) Cgo() bool {
	for _, mod := range mods {
		if mod.Cgo {
			return true
		}
	}

	return false
}

const k6ModulePath = "go.k6.io/k6"
