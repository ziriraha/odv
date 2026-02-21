package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type branchName string

func (b branchName) FilterValue() string { return string(b) }

type branchDelegate struct{}

func (d branchDelegate) Height() int                             { return 1 }
func (d branchDelegate) Spacing() int                            { return 0 }
func (d branchDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d branchDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	branch, ok := listItem.(branchName)
	if !ok {
		return
	}

	fn := ListItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return ListSelectedItemStyle.Render("→ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(string(branch)))
}

type branchListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m branchListModel) Init() tea.Cmd { return nil }

func (m branchListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil

	case tea.KeyMsg:
		keypress := msg.String()
		if keypress == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if m.list.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch keypress {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if branch, ok := m.list.SelectedItem().(branchName); ok {
				m.choice = string(branch)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m branchListModel) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return ListCancelStyle.Render("Cancelled.\n")
	}
	return m.list.View()
}

type BranchSelectListView struct {
	Title    string
	Branches []string
}

func (cfg BranchSelectListView) Run() (string, error) {
	const defaultWidth = 80
	const defaultHeight = 20

	items := make([]list.Item, len(cfg.Branches))
	for i, b := range cfg.Branches {
		items[i] = branchName(b)
	}

	l := list.New(items, branchDelegate{}, defaultWidth, defaultHeight)
	l.Title = cfg.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ListTitleStyle
	l.Styles.PaginationStyle = ListPaginationStyle
	l.Styles.HelpStyle = ListHelpStyle

	p := tea.NewProgram(branchListModel{list: l})
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running program: %w", err)
	}

	if m, ok := finalModel.(branchListModel); ok {
		return m.choice, nil
	}
	return "", nil
}
