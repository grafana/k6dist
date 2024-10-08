// Package registry contains extension registry related internal helpers.
package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ToModules creates a module list from the registry.
func (reg Registry) ToModules() Modules {
	mods := make(Modules, 0, len(reg))

	for _, ext := range reg {
		if len(ext.Versions) == 0 {
			continue
		}

		mods = append(mods, Module{Path: ext.Module, Version: ext.Versions[0], Cgo: ext.Cgo})
	}

	return mods
}

// LoadRegistry loads registry from URL of filesystem path.
func LoadRegistry(ctx context.Context, source string) (reg Registry, err error) {
	if strings.HasPrefix(source, "http:") || strings.HasPrefix(source, "https:") {
		reg, err = loadRegistryHTTP(ctx, source)
	} else {
		reg, err = loadRegistryFS(source)
	}

	if err != nil {
		return nil, err
	}

	sort.Slice(reg, func(i, j int) bool {
		if reg[i].Module == k6ModulePath {
			return true
		}

		if reg[j].Module == k6ModulePath {
			return false
		}

		return reg[i].Module < reg[j].Module
	})

	for idx := range reg {
		reg[idx].Versions = []string{reg[idx].Versions[0], ""}
	}

	return reg, nil
}

func loadRegistryFS(source string) (Registry, error) {
	file, err := os.Open(filepath.Clean(source)) //nolint:forbidigo
	if err != nil {
		return nil, err
	}

	defer file.Close() //nolint:errcheck

	decoder := json.NewDecoder(file)

	var reg Registry

	if err := decoder.Decode(&reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func loadRegistryHTTP(ctx context.Context, source string) (Registry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", os.ErrNotExist, resp.Status) //nolint:forbidigo
	}

	defer resp.Body.Close() //nolint:errcheck

	decoder := json.NewDecoder(resp.Body)

	var reg Registry

	if err := decoder.Decode(&reg); err != nil {
		return nil, err
	}

	return reg, nil
}

// AddLatest adds latest versions to extensions.
func (reg Registry) AddLatest(modules Modules) bool {
	regAsMap := make(map[string]*Extension, len(reg))

	for idx := range reg {
		regAsMap[reg[idx].Module] = &reg[idx]
	}

	changed := false

	for _, mod := range modules {
		ext, found := regAsMap[mod.Path]
		if found {
			ext.Versions[1] = mod.Version
			changed = changed || (ext.Versions[0] != mod.Version)
		}

		changed = changed || !found
	}

	if changed {
		return true
	}

	modulesAsMap := make(map[string]*Module, len(modules))

	for idx := range modules {
		modulesAsMap[modules[idx].Path] = &modules[idx]
	}

	for _, ext := range reg {
		if _, found := modulesAsMap[ext.Module]; !found {
			return true
		}
	}

	return false
}
