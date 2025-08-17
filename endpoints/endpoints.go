package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Endpoint struct {
	Endpoint string   `json:"endpoint"`
	Command  []string `json:"command"`
}

type Manager struct {
	mu        sync.RWMutex
	endpoints map[string][]string
	lastMod   time.Time
}

const endpointFile = "endpoints.json"

var manager = &Manager{
	endpoints: make(map[string][]string),
}

func Init() error {
	if err := load(); err != nil {
		return err
	}

	f, err := os.Stat(endpointFile)
	if err != nil {
		return err
	}

	manager.lastMod = f.ModTime()
	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	reloadIfNeeded()

	manager.mu.RLock()
	cmd, ok := manager.endpoints[r.URL.Path]
	manager.mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	err := exec.Command(cmd[0], cmd[1:]...).Run()
	if err != nil {
		http.Error(w,
			fmt.Sprintf("command failed: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func load() error {
	f, err := os.Open(endpointFile)
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(endpointFile, []byte(`{"endpoints":[]}`), 0644)
	}
	if err != nil {
		return fmt.Errorf("failed to open endpoint file: %w", err)
	}
	defer f.Close()

	var data struct {
		Endpoints []Endpoint `json:"endpoints"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	tmp := make(map[string][]string, len(data.Endpoints))
	for _, e := range data.Endpoints {
		if len(e.Command) == 0 {
			continue
		}
		tmp[e.Endpoint] = e.Command
	}

	manager.mu.Lock()
	manager.endpoints = tmp
	manager.mu.Unlock()

	return nil
}

func reloadIfNeeded() {
	f, err := os.Stat(endpointFile)
	if err != nil {
		return
	}

	manager.mu.RLock()
	needsReload := f.ModTime().After(manager.lastMod)
	manager.mu.RUnlock()

	if !needsReload {
		return
	}

	if err := load(); err != nil {
		fmt.Printf("failed to reload endpoints: %s", err)
		return
	}

	manager.mu.Lock()
	manager.lastMod = f.ModTime()
	manager.mu.Unlock()
}
