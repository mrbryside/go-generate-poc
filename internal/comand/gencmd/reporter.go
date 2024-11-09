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
	// Summary texts
	successLenText := renderSummaryTextWithLen(len(r.HandlerTemplateSuccessRoute), "SUCCESS", greenSuccessSummary)
	failedLenText := renderSummaryTextWithLen(len(r.HandlerTemplateErrorRoute), "FAILED", redErrorSummary)
	failedTextOnly := renderSummaryText("FAILED", redErrorSummary)

	var sb strings.Builder
	if r.MandaToryError.Error != nil {
		sb.WriteString(fmt.Sprintf("\n%s %s %s \n", style.Render("PATH"), r.MandaToryError.Path, failedTextOnly))
		sb.WriteString(underLineText.Render("==================================================") + "\n")
		sb.WriteString(renderMandatoryError(r.MandaToryError.Error))
		return sb.String()
	}

	swaggerReportText := ""
	if r.SwagGenerateReport.Error == nil {
		swaggerReportText = renderSummaryText("SWAGGER SUCCESS", greenSuccessSummary)
	}
	if r.SwagGenerateReport.Error != nil {
		swaggerReportText = renderSummaryText("SWAGGER FAILED", redErrorSummary)
	}

	// Build header
	sb.WriteString(fmt.Sprintf("\n%s %s • %s | %s | %s\n", style.Render("PATH"), r.BasePathOfJsonSpec, successLenText, failedLenText, swaggerReportText))
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
		sb.WriteString(renderRouteText("✗", e.Path, redErrorText) + "\n")
		sb.WriteString(renderErrorDetails([]error{e.Error}))
	}

	// for debugger only
	if r.SwagGenerateReport.Error != nil {
		sb.WriteString(renderRouteText("✗", "swaggo", redErrorText) + "\n")
		sb.WriteString(renderErrorDetails([]error{r.SwagGenerateReport.Error}))
	}

	return sb.String()
}

// Updated renderSummaryTextWithLen to return a string
func renderSummaryTextWithLen(count int, label string, style lipgloss.Style) string {
	return fmt.Sprintf("%v %v", style.Render(fmt.Sprintf("%v", count)), style.Render(label))
}

// Updated renderSummaryTextWithLen to return a string
func renderSummaryText(label string, style lipgloss.Style) string {
	return style.Render(label)
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

func renderMandatoryError(e error) string {
	var sb strings.Builder
	errorDetailStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3333"))
	sb.WriteString(errorDetailStyle.Render(fmt.Sprintf("• %s", e.Error())) + "\n")
	return sb.String()
}
