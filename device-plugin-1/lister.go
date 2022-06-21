package main

import (
	"fmt"
	"os"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

// Lister is an interface between plugin imlementation and kubernetes Device Plugin Manager.
// See call of NewManager function in main.go
// Manager will use it to obtain resource namespace,
// monitor available resources and instantate a new plugin for them.
//
type mylister struct{}

// GetResourceNamespace must return namespace (vendor ID) of implemented Lister.
// e.g. for resources in format "color.example.com/<color>" that would be "color.example.com".
//
func (l *mylister) GetResourceNamespace() string {
	fmt.Printf("API: GetResourceNamespace\n")

	return "myplugin.io"
}

// NewPlugin instantiates a plugin implementation.
// It is given the last name of the resource,
// e.g. for resource name "color.example.com/red" that would be "red".
// It must return valid implementation of a PluginInterface.
//
func (l *mylister) NewPlugin(resourceLastName string) dpm.PluginInterface {
	fmt.Printf("API: NewPlugin %s\n", resourceLastName)

	return &myplugin{
		deviceType: resourceLastName,
		nodeName:   os.Getenv("NODE_NAME"),
	}
}

// Discover notifies plugin manager with a list of currently available resources in its namespace.
// e.g. if "color.example.com/red" and "color.example.com/blue" are available in the system,
// it would pass PluginNameList{"red", "blue"} to given channel.
// In case list of resources is static, it would use the channel only once and then return.
// In case the list is dynamic, it could block and pass a new list each times resources changed.
// If blocking is used, it should check whether the channel is closed, i.e. Discover should stop.
//
func (l *mylister) Discover(pluginListCh chan dpm.PluginNameList) {
	pluginListCh <- []string{
		deviceType1,
		deviceType2,
		deviceType3,
	}
}
