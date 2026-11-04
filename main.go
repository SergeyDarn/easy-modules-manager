package main

import (
	"github.com/SergeyDarn/easy-modules-manager/manager"
	"github.com/charmbracelet/log"
)

func main() {
	manager.Init(log.DebugLevel)
}
