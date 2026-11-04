package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SergeyDarn/easy-modules-manager/utils"

	"github.com/charmbracelet/log"
)

type JsonConfig struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
}

var modulesDirPermissions os.FileMode = 0o777

func getModulesDir() string {
	return utils.GetModulesDir()
}

func ReadConfigJson() JsonConfig {
	configFile := utils.GetConfigFile()
	return utils.ReadJson[JsonConfig](configFile)
}

func CreateModulesDir() {
	err := os.MkdirAll(getModulesDir(), modulesDirPermissions)
	utils.CheckError(err, "Error when creating modules folder")
}

func RemoveModulesDir() {
	err := os.RemoveAll(getModulesDir())
	utils.CheckError(err, "Error while trying to delete modules folder")
	log.Debug(utils.ColorDanger("Modules folder deleted before installation"))
}

func InstallModules(modules map[string]string, parallelInstall bool) {
	log.Infof("Installing modules into %s", getModulesDir())
	fmt.Println()

	start := time.Now()

	if !parallelInstall {
		for name, url := range modules {
			installModule(name, url)
		}
	} else {
		var waitGroup sync.WaitGroup

		for name, url := range modules {
			waitGroup.Add(1)

			go func() {
				installModule(name, url)
				defer waitGroup.Done()
			}()
		}

		waitGroup.Wait()
	}

	log.Infof(
		utils.ColorSuccess("Installation of %d modules took %s"),
		len(modules),
		time.Since(start),
	)
}

func installModule(moduleName string, moduleUrl string) {
	if !utils.IsGitUrl(moduleUrl) {
		return
	}

	moduleDir := getModuleDir(moduleName)
	err, isModuleNotCloned := checkModuleDirStatus(moduleDir)

	if isModuleNotCloned {
		utils.GitClone(moduleName, moduleUrl, moduleDir)
		return
	}

	utils.CheckError(err, "Error while reading module "+moduleName+" folder")

	gitStatus := utils.GitDirStatus(moduleDir)
	if gitStatus.String() != "" {
		log.Warnf(
			utils.ColorWarning(
				"\nThere are unsaved changes for module \"%s\" - skipping it\n"+
					"\n%s",
			),
			moduleName,
			gitStatus.String(),
		)
		return
	}

	err = os.RemoveAll(moduleDir)
	utils.CheckError(err, "Error while trying to delete module folder for "+moduleName)

	utils.GitClone(moduleName, moduleUrl, moduleDir)
}

func ShowChangedModules() {
	modules, err := os.ReadDir(getModulesDir())
	utils.CheckError(err, "Error reading modules folder")
	changedModules := []string{}

	for _, module := range modules {
		gitStatus := utils.GitDirStatus(
			getModuleDir(module.Name()),
		)

		if gitStatus.String() != "" {
			changedModules = append(changedModules, module.Name())
		}
	}

	if len(changedModules) == 0 {
		utils.LogSuccess("You have no unsaved modules")
		return
	}

	log.Infof(
		utils.ColorWarning("\nUnsaved Modules (%d):\n\n%s"),
		len(changedModules),
		strings.Join(changedModules, "\n"),
	)
}

func checkModuleDirStatus(moduleDir string) (error, bool) {
	_, err := os.Stat(moduleDir)
	return err, os.IsNotExist(err)
}

func getModuleDir(moduleName string) string {
	return filepath.Join(getModulesDir(), moduleName)
}
