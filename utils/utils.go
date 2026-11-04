package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	successColor = lipgloss.Color("#75d57f")
	warningColor = lipgloss.Color("192")
	dangerColor  = lipgloss.Color("204")
)

func CheckError(err error, readableErrMessage string) {
	if err == nil {
		return
	}

	ThrowError(readableErrMessage + ": " + err.Error())
}

func ThrowError(err string) {
	panic(ColorDanger(err))
}

func GetPathRootDir(path string) string {
	splitPath := strings.Split(path, string(os.PathSeparator))

	if len(splitPath) == 0 {
		return ""
	}

	return splitPath[0]
}

func LogSuccess(text string) {
	log.Info(ColorSuccess(text))
}

func LogWarning(text string) {
	log.Warn(ColorWarning(text))
}

func LogError(text string) {
	log.Error(ColorDanger(text))
}

func ColorText(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(lipgloss.TerminalColor(color)).Render(text)
}

func ColorDanger(text string) string {
	return ColorText(text, dangerColor)
}

func ColorSuccess(text string) string {
	return ColorText(text, successColor)
}

func ColorWarning(text string) string {
	return ColorText(text, warningColor)
}

func ReadJson[T any](file string) T {
	configJson, fileErr := os.ReadFile(file)
	CheckError(fileErr, "Error when opening json configuration file")

	var configJsonParsed T
	jsonErr := json.Unmarshal(configJson, &configJsonParsed)
	CheckError(jsonErr, "Error when parsing json configuration file")

	return configJsonParsed
}
