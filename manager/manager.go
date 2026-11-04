package manager

import (
	"flag"
	"maps"

	"github.com/SergeyDarn/easy-modules-manager/modules"
	"github.com/SergeyDarn/easy-modules-manager/utils"
	"github.com/charmbracelet/log"
)

var ModulesMap = map[string]string{}

func Init() {
	log.SetLevel(log.DebugLevel)

	utils.InitEnv()

	showChangedModules := flag.Bool("show-changed-modules", false, "run command to show modules with unsaved git changes")
	parallelInstall := flag.Bool("parallel-install", true, "install modules in parallel (true/false)")
	safeInstall := flag.Bool("safe-install", true, "if this is set to false, modules folder will be deleted on start. Default version - each module is checked separately, and only if module has no unsaved changes, it's deleted and then reinstalled")
	flag.Parse()

	if *showChangedModules {
		modules.ShowChangedModules()
		return
	}

	configJson := modules.ReadConfigJson()
	dependencies := configJson.Dependencies
	maps.Copy(dependencies, configJson.DevDependencies)

	for key, value := range dependencies {
		if utils.IsGitUrl(value) {
			ModulesMap[key] = value
		}
	}

	if len(ModulesMap) == 0 {
		log.Info(utils.PrepareWarningOutput("No git modules to install."))
		return
	}

	if !*safeInstall {
		modules.RemoveModulesDir()
	}

	modules.CreateModulesDir()
	modules.InstallModules(ModulesMap, *parallelInstall)
}
