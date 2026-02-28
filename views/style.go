package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/ziriraha/odv/lib"
)

func GetRepoStyle(repoName string) lipgloss.Style {
	var repoColor string
	switch repoName {
	case lib.WorkspaceRepo:
		repoColor = "1" // Red
	case "community":
		repoColor = "3" // Yellow
	case "enterprise":
		repoColor = "2" // Green
	case "upgrade":
		repoColor = "4" // Blue
	default:
		repoColor = "7" // Default (white)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(repoColor)).Bold(true)
}

func RenderRepoName(repoName string) string {
	return GetRepoStyle(repoName).Render(repoName)
}

func RenderRepoLetter(repoName string) string {
	letter := repoName[0:1]
	return GetRepoStyle(repoName).Render(letter)
}

var (
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	FaintStyle   = lipgloss.NewStyle().Faint(true)
	BoldStyle    = lipgloss.NewStyle().Bold(true)

	// Header
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	// Spinner
	SpinnerStyle = list.DefaultStyles().Spinner

	// List
	ListTitleStyle        = HeaderStyle
	ListItemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	ListSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	ListPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	ListHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	ListCancelStyle       = WarningStyle

	// Diff
	DiffAddedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))            // Green
	DiffModifiedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))            // Yellow
	DiffDeletedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))            // Red
	DiffUntrackedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))          // Gray
	DiffConflictStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true) // Red, bold

	AheadStyle       = DiffAddedStyle
	BehindStyle      = DiffDeletedStyle
	LocalBranchStyle = DiffUntrackedStyle
)

var (
	Checkmark = SuccessStyle.Render("✓")
	Cross     = ErrorStyle.Render("✗")
)

func colorizeIndicator(s string) string {
	switch s {
	case "A", "R":
		return DiffAddedStyle.Render(s)
	case "M":
		return DiffModifiedStyle.Render(s)
	case "D":
		return DiffDeletedStyle.Render(s)
	case "?":
		return DiffUntrackedStyle.Render(s)
	default:
		return s
	}
}

func ColorizeStatusIndicator(status string) string {
	parts := strings.Split(status, "")
	if len(parts) != 2 {
		return status
	}

	switch status {
	case "UU", "AA", "DD", "AU", "UA", "DU", "UD":
		return DiffConflictStyle.Render(status)
	case "!!":
		return DiffUntrackedStyle.Render(status)
	default:
		return colorizeIndicator(parts[0]) + colorizeIndicator(parts[1])
	}
}
