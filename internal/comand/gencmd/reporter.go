package gencmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mrbryside/go-generate/internal/generator/handlergen"
)

// Styles
var (
	style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#00b33c")).
		PaddingRight(1).
		PaddingLeft(1)

	greenSuccessText = lipgloss.NewStyle().
				Bold(true).
				MarginLeft(1).
				Foreground(lipgloss.Color("#00b33c")).
				PaddingLeft(1).
				PaddingRight(1)

	greenSuccessSummary = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00b33c"))

	redErrorSummary = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff4d4d"))

	redErrorText = lipgloss.NewStyle().
			Bold(true).
			MarginLeft(1).
			Foreground(lipgloss.Color("#ff4d4d")).
			PaddingLeft(1).
			PaddingRight(1)

	underLineText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff9e6"))
)

// Refactored printReports to return accumulated string result
func printReports(rs []handlergen.Report) {
	var sb strings.Builder
	for _, r := range rs {
		sb.WriteString(printReport(r))
	}
	fmt.Println(sb.String())
}

// Refactored printReport to accumulate results in a string builder
func printReport(r handlergen.Report) string {
	var sb strings.Builder

	// Summary texts
	successLenText := renderSummaryText(len(r.HandlerTemplateSuccessRoute), "SUCCESS", greenSuccessSummary)
	failedLenText := renderSummaryText(len(r.HandlerTemplateErrorRoute), "FAILED", redErrorSummary)

	// Build header
	sb.WriteString(fmt.Sprintf("\n%s %s • %s | %s\n", style.Render("PATH"), r.BasePathOfJsonSpec, successLenText, failedLenText))
	sb.WriteString(underLineText.Render("==================================================") + "\n")

	// Build success routes
	for _, h := range r.HandlerTemplateSuccessRoute {
		sb.WriteString(renderRouteText("✔", h.Name, greenSuccessText) + "\n")
	}

	// Build error routes
	for _, h := range r.HandlerTemplateErrorRoute {
		sb.WriteString(renderRouteText("✗", h.HandlerNameTemplateData.Name, redErrorText) + "\n")
		sb.WriteString(renderErrorDetails(h.Errors))
	}

	// Build specific file errors for debugging (handler.go and routes_gen.go)
	for _, e := range r.PathToGenerateError {
		if strings.Contains(e.Path, "handler.go") {
			sb.WriteString(renderRouteText("✗", "handler", redErrorText) + "\n")
			sb.WriteString(renderErrorDetails([]error{e.Error}))
		}
		if strings.Contains(e.Path, "routes_gen.go") {
			sb.WriteString(renderRouteText("✗", "routes", redErrorText) + "\n")
			sb.WriteString(renderErrorDetails([]error{e.Error}))
		}
	}

	return sb.String()
}

// Updated renderSummaryText to return a string
func renderSummaryText(count int, label string, style lipgloss.Style) string {
	return fmt.Sprintf("%v %v", style.Render(fmt.Sprintf("%v", count)), style.Render(label))
}

// Updated renderRouteText to return a string
func renderRouteText(icon string, name string, style lipgloss.Style) string {
	return lipgloss.NewStyle().Bold(true).Render(style.Render(fmt.Sprintf("%s", icon)) + name)
}

// Updated renderErrorDetails to return a string
func renderErrorDetails(errors []error) string {
	var sb strings.Builder
	errorDetailStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3333"))
	for _, e := range errors {
		sb.WriteString(errorDetailStyle.Render(fmt.Sprintf("    • %s", e.Error())) + "\n")
	}
	return sb.String()
}
