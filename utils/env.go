package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvVariable int

const (
	envRoot EnvVariable = iota
	envConfigFile
	envModulesDir
	envSshKeyPath
	envSshKeyPassword
)

var envMap = map[EnvVariable]string{
	envRoot:           "ENV_ROOT",
	envConfigFile:     "CONFIG_FILE",
	envModulesDir:     "MODULES_DIR",
	envSshKeyPath:     "SSH_KEY_PATH",
	envSshKeyPassword: "SSH_KEY_PASSWORD",
}

func InitEnv() {
	envRoot := getEnv(envRoot)

	err := godotenv.Load(filepath.Join(envRoot, "go.env.local"))
	if err != nil {
		LogWarning("FYI: couldn't open go.env.local")
	}

	err = godotenv.Load(filepath.Join(envRoot, "go.env"))
	CheckError(err, "Error loading go.env file")
}

func GetModulesDir() string {
	return getEnv(envModulesDir)
}

func GetConfigFile() string {
	return getEnv(envConfigFile)
}

func GetSshKeyPath() string {
	return getEnv(envSshKeyPath)
}

func GetSshKeyPassword() string {
	return getEnv(envSshKeyPassword)
}

func getEnv(env EnvVariable) string {
	return os.Getenv(envMap[env])
}
