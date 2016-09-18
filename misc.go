// Copyright 2013 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/engine-api/types/swarm"
)

// Version returns version information about the docker server.
//
// See https://goo.gl/ND9R8L for more details.
func (c *Client) Version() (*Env, error) {
	resp, err := c.do(http.MethodGet, "/version", doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env Env
	if err := env.Decode(resp.Body); err != nil {
		return nil, err
	}
	return &env, nil
}

// DockerInfo contains information about the Docker server
//
// See https://goo.gl/bHUoz9 for more details.
type DockerInfo struct {
	ID                 string
	Containers         int
	ContainersRunning  int
	ContainersPaused   int
	ContainersStopped  int
	Images             int
	Driver             string
	DriverStatus       [][2]string
	SystemStatus       [][2]string
	Plugins            PluginsInfo
	MemoryLimit        bool
	SwapLimit          bool
	KernelMemory       bool
	CPUCfsPeriod       bool `json:"CpuCfsPeriod"`
	CPUCfsQuota        bool `json:"CpuCfsQuota"`
	CPUShares          bool
	CPUSet             bool
	IPv4Forwarding     bool
	BridgeNfIptables   bool
	BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
	Debug              bool
	OomKillDisable     bool
	ExperimentalBuild  bool
	NFd                int
	NGoroutines        int
	SystemTime         string
	ExecutionDriver    string
	LoggingDriver      string
	CgroupDriver       string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	OSType             string
	Architecture       string
	IndexServerAddress string
	NCPU               int
	MemTotal           int64
	DockerRootDir      string
	HTTPProxy          string `json:"HttpProxy"`
	HTTPSProxy         string `json:"HttpsProxy"`
	NoProxy            string
	Name               string
	Labels             []string
	ServerVersion      string
	ClusterStore       string
	ClusterAdvertise   string
	Swarm              swarm.Info
}

// PluginsInfo is a struct with the plugins registered with the docker daemon
//
// for more information, see: https://goo.gl/bHUoz9
type PluginsInfo struct {
	// List of Volume plugins registered
	Volume []string
	// List of Network plugins registered
	Network []string
	// List of Authorization plugins registered
	Authorization []string
}

// Info returns system-wide information about the Docker server.
//
// See https://goo.gl/ElTHi2 for more details.
func (c *Client) Info() (*DockerInfo, error) {
	resp, err := c.do(http.MethodGet, "/info", doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info DockerInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ParseRepositoryTag gets the name of the repository and returns it splitted
// in two parts: the repository and the tag. It ignores the digest when it is
// present.
//
// Some examples:
//
//     localhost.localdomain:5000/samalba/hipache:latest -> localhost.localdomain:5000/samalba/hipache, latest
//     localhost.localdomain:5000/samalba/hipache -> localhost.localdomain:5000/samalba/hipache, ""
//     busybox:latest@sha256:4a731fb46adc5cefe3ae374a8b6020fc1b6ad667a279647766e9a3cd89f6fa92 -> busybox, latest
func ParseRepositoryTag(repoTag string) (repository string, tag string) {
	parts := strings.SplitN(repoTag, "@", 2)
	repoTag = parts[0]
	n := strings.LastIndex(repoTag, ":")
	if n < 0 {
		return repoTag, ""
	}
	if tag := repoTag[n+1:]; !strings.Contains(tag, "/") {
		return repoTag[:n], tag
	}
	return repoTag, ""
}
