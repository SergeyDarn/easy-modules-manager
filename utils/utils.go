package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	SUCCESS_COLOR = lipgloss.Color("#75d57f")
	WARNING_COLOR = lipgloss.Color("192")
	DANGER_COLOR  = lipgloss.Color("204")
)

func CheckError(err error, readableErrMessage string) {
	if err == nil {
		return
	}

	ThrowError(readableErrMessage + ": " + err.Error())
}

func ThrowError(err string) {
	panic(PrepareDangerOutput(err))
}

func GetPathRootDir(path string) string {
	splitPath := strings.Split(path, string(os.PathSeparator))

	if len(splitPath) == 0 {
		return ""
	}

	return splitPath[0]
}

func PrepareSuccessOutput(output string) string {
	return PrepareColorOutput(output, SUCCESS_COLOR)
}

func PrepareWarningOutput(output string) string {
	return PrepareColorOutput(output, WARNING_COLOR)
}

func PrepareDangerOutput(output string) string {
	return PrepareColorOutput(output, DANGER_COLOR)
}

func PrepareColorOutput(output string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(lipgloss.TerminalColor(color)).Render(output)
}

func ReadJson[T any](file string) T {
	configJson, fileErr := os.ReadFile(file)
	CheckError(fileErr, "Error when opening json configuration file")

	var configJsonParsed T
	jsonErr := json.Unmarshal(configJson, &configJsonParsed)
	CheckError(jsonErr, "Error when parsing json configuration file")

	return configJsonParsed
}
