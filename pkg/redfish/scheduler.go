package redfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/ironic"
	"github.com/stmcginnis/gofish"
)

// Node is ...
type Node struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// Redfish is ...
type Redfish struct {
	cfg config.Config
}

// New Redfish Instance
func New(cfg config.Config) Redfish {
	r := Redfish{
		cfg: cfg,
	}
	return r
}

// Start ...
func (r Redfish) Start(ctx context.Context, errors chan<- error) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

loop:
	for {
		nodes, err := r.loadNodes()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, node := range nodes {
			bm, err := r.loadRedfishInfo(node)
			if err != nil {
				continue
			}
			r.createIronicNode(bm)
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			break loop
		}
	}
}

func (r Redfish) loadNodes() (ips []string, err error) {
	d, err := ioutil.ReadFile(r.cfg.NetboxNodesPath)
	if err != nil {
		return
	}

	t := make([]Node, 0)
	if err = json.Unmarshal(d, &t); err != nil {
		return
	}

	for _, node := range t {
		ips = append(ips, node.Targets...)
	}

	return
}

func (r Redfish) loadRedfishInfo(nodeIP string) (i ironic.InspectorCallbackData, err error) {
	fmt.Println(nodeIP)
	cfg := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("https://%s", nodeIP),
		Username:  r.cfg.IronicUser,
		Password:  r.cfg.IronicPassword,
		Insecure:  true,
		BasicAuth: false,
	}
	c, err := gofish.Connect(cfg)
	if err != nil {
		return
	}
	defer c.Logout()
	service := c.Service
	chassis, err := service.Chassis()
	if err != nil {
		return
	}
	for _, chass := range chassis {
		n, err := chass.NetworkAdapters()
		if err != nil {
			continue
		}
		if len(n) == 0 {
			continue
		}
		i.Interfaces = make([]ironic.Interface, len(n))
		f, err := n[0].NetworkDeviceFunctions()
		i.Interfaces[0].MacAddress = f[0].Ethernet.MACAddress
		fmt.Println(f[0].Ethernet.MACAddress)
		fmt.Printf("Chassis: %#v\n\n", chass.Manufacturer)
	}
	return
}

func (r Redfish) createIronicNode(i ironic.InspectorCallbackData) (err error) {
	return ironic.CreateIronicNode(i, r.cfg.IronicInspectorCallback)
}

func (r Redfish) checkIronicNodeCreation(i ironic.InspectorCallbackData) {

}

func (r Redfish) checkIronicNodeExists(i ironic.InspectorCallbackData) {

}

func (r Redfish) updateNetbox(i ironic.InspectorCallbackData) {
	// update provision_state
}
