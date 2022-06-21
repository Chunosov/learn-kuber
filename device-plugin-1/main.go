package main

import (
	"flag"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	// We don't use flags here but glog (used in dpm) uses them internally
	// and says "ERROR: logging before flag.Parse" if we don't parse them
	flag.Parse()

	lister := &mylister{}

	manager := dpm.NewManager(lister)

	manager.Run()
}
