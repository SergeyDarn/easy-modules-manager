package utils

import (
	"encoding/json"
	"fmt"
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

	ThrowError(PrepareError(err, readableErrMessage))
}

func CheckErrorf(err error, readableErrMessage string, args ...any) {
	CheckError(err, fmt.Sprintf(readableErrMessage, args...))
}

func ThrowError(err string) {
	panic(ColorDanger(err))
}

func ThrowErrorf(err string, args ...any) {
	ThrowError(fmt.Sprintf(err, args...))
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

func LogSuccessf(text string, args ...any) {
	LogSuccess(fmt.Sprintf(text, args...))
}

func LogWarning(text string) {
	log.Warn(ColorWarning(text))
}

func LogWarningf(text string, args ...any) {
	LogWarning(fmt.Sprintf(text, args...))
}

func LogError(text string) {
	log.Error(ColorDanger(text))
}

func LogErrorf(text string, args ...any) {
	LogError(fmt.Sprintf(text, args...))
}

func LogReadableError(err error, readableErrMessage string) {
	LogError(PrepareError(err, readableErrMessage))
}

func LogReadableErrorf(err error, readableErrMessage string, args ...any) {
	LogReadableError(err, fmt.Sprintf(readableErrMessage, args...))
}

func PrepareError(err error, readableErrMessage string) string {
	return fmt.Sprintf("%s: %s", readableErrMessage, err.Error())
}

func ColorText(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(lipgloss.TerminalColor(color)).Render(text)
}

func ColorTextf(text string, color lipgloss.Color, args ...any) string {
	return ColorText(fmt.Sprintf(text, args...), color)
}

func ColorDanger(text string) string {
	return ColorText(text, dangerColor)
}

func ColorDangerf(text string, args ...any) string {
	return ColorTextf(text, dangerColor, args...)
}

func ColorSuccess(text string) string {
	return ColorText(text, successColor)
}

func ColorSuccessf(text string, args ...any) string {
	return ColorTextf(text, successColor, args...)
}

func ColorWarning(text string) string {
	return ColorText(text, warningColor)
}

func ColorWarningf(text string, args ...any) string {
	return ColorTextf(text, warningColor, args...)
}

func ReadJson[T any](file string) T {
	configJson, fileErr := os.ReadFile(file)
	CheckError(fileErr, "error when opening json configuration file")

	var configJsonParsed T
	jsonErr := json.Unmarshal(configJson, &configJsonParsed)
	CheckError(jsonErr, "error when parsing json configuration file")

	return configJsonParsed
}
