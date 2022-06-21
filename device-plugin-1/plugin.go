package main

import (
	"context"
	"fmt"
	"time"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	nodeName1   = "kind-worker"
	nodeName2   = "kind-worker2"
	deviceType1 = "mydevice-1"
	deviceType2 = "mydevice-2"
	deviceType3 = "mydevice-3"
)

var (
	// Available devices per node per device type
	deviceConfig = map[string]map[string][]string{
		nodeName1: {
			deviceType1: []string{"/node1/dev1/1"},
		},
		nodeName2: {
			deviceType2: []string{"/node2/dev2/1"},
			deviceType3: []string{"/node2/dev3/1", "/node2/dev3/2", "/node2/dev3/3"},
		},
	}
)

// myplugin implements device plugin API
type myplugin struct {
	deviceType string
	nodeName   string
}

// getDevices returns list of available devices
// In reality the driver should look through some /sys/* dirs on the host
// to detect available devices and their properties and state.
// Our imaginary devices are always available and healthy.
func (p *myplugin) getDevices() []*pluginapi.Device {
	devices := []*pluginapi.Device{}

	if devs, ok := deviceConfig[p.nodeName]; ok {
		if ids, ok := devs[p.deviceType]; ok {
			for _, id := range ids {
				devices = append(devices, &pluginapi.Device{
					ID:     id,
					Health: pluginapi.Healthy,
				})
			}
		}
	}

	return devices
}

// Start is an optional interface that could be implemented by plugin.
// If case Start is implemented, it will be executed by Manager after
// plugin instantiation and before its registration to kubelet. This
// method could be used to prepare resources before they are offered
// to Kubernetes.
func (p *myplugin) Start() error {
	fmt.Printf("Start plugin %s\n", p.deviceType)

	return nil
}

// Stop is an optional interface that could be implemented by plugin.
// If case Stop is implemented, it will be executed by Manager after the
// plugin is unregistered from kubelet. This method could be used to tear
// down resources.
func (p *myplugin) Stop() error {
	fmt.Printf("Stop plugin %s\n", p.deviceType)

	return nil
}

// GetDevicePluginOptions returns options to be communicated with Device Manager
func (p *myplugin) GetDevicePluginOptions(ctx context.Context, e *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	fmt.Printf("API: GetDevicePluginOptions\n")

	return &pluginapi.DevicePluginOptions{}, nil
}

// PreStartContainer is called, if indicated by Device Plugin during registeration phase, before each container start.
// Device plugin can run device specific operations such as resetting the device before making devices available to the container.
// Plugins are not required to provide useful implementations for PreStartContainer
func (p *myplugin) PreStartContainer(ctx context.Context, r *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	fmt.Printf("API: PreStartContainer %v\n", r)

	return &pluginapi.PreStartContainerResponse{}, nil
}

// ListAndWatch returns a stream of List of Devices.
// Whenever a Device state change or a Device disappears, ListAndWatch returns the new list
func (p *myplugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	for {
		s.Send(&pluginapi.ListAndWatchResponse{Devices: p.getDevices()})

		time.Sleep(10 * time.Second)
	}
	// returning a value from this function will unregister the plugin from k8s
}

// Allocate is called during container creation so that the Device Plugin can run device specific operations
// and instruct Kubelet of the steps to make the Device available in the container
func (p *myplugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	fmt.Printf("API: Allocate %v\n", r)

	resps := []*pluginapi.ContainerAllocateResponse{}

	for _, req := range r.ContainerRequests {
		fmt.Printf("Allocate for container %v\n", req)

		devices := []*pluginapi.DeviceSpec{}
		envs := map[string]string{}

		for i, id := range req.DevicesIDs {
			fmt.Printf("Allocate device %v\n", id)

			devices = append(devices, &pluginapi.DeviceSpec{
				// These paths must be real
				// Kubelet checks them before creating a container
				ContainerPath: "/tmp",
				HostPath:      "/tmp",
				Permissions:   "r",
			})

			// Provide some device specific vars into container
			switch p.deviceType {
			case deviceType1:
				envs[fmt.Sprintf("DEV_1_ID_%d", i)] = id
			case deviceType2:
				envs[fmt.Sprintf("DEV_2_ID_%d", i)] = id
			case deviceType3:
				envs[fmt.Sprintf("DEV_3_ID_%d", i)] = id
			}
		}

		resps = append(resps, &pluginapi.ContainerAllocateResponse{
			Envs:    envs,
			Devices: devices,
		})

		fmt.Printf("Devices allocated: %v\n", resps)
	}

	return &pluginapi.AllocateResponse{ContainerResponses: resps}, nil
}

// GetPreferredAllocation returns a preferred set of devices to allocate from a list of available ones.
// The resulting preferred allocation is not guaranteed to be the allocation ultimately performed by the devicemanager.
// It is only designed to help the devicemanager make a more informed allocation decision when possible.
// Plugins are not required to provide useful implementations for GetPreferredAllocation
func (p *myplugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	fmt.Printf("API: GetPreferredAllocation %v\n", r)

	return &pluginapi.PreferredAllocationResponse{}, nil
}
