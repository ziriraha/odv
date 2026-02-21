package views

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type OperationStatus int

const (
	StatusPending OperationStatus = iota
	StatusInProgress
	StatusDone
	StatusFailed
)

type RepoOperationState struct {
	Name      string
	Status    OperationStatus
	Err       error
	Spinner   spinner.Model
	StartTime time.Time
	Duration  time.Duration
}

type RepoOperationDoneMsg struct {
	RepoIndex int
	Err       error
	Duration  time.Duration
}

func NewRepoOperationState(repoName string) RepoOperationState {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle
	return RepoOperationState{
		Name:    repoName,
		Status:  StatusPending,
		Spinner: s,
	}
}

func (state *RepoOperationState) RenderInProgress(message string) string {
	elapsed := time.Since(state.StartTime).Round(time.Millisecond)
	return fmt.Sprintf("%s%s - %s (%s)\n",
		state.Spinner.View(),
		RenderRepoName(state.Name),
		message,
		FaintStyle.Render(elapsed.String()))
}

func (state *RepoOperationState) RenderDone(message string) string {
	return fmt.Sprintf("%s %s - %s (%s)\n",
		Checkmark,
		RenderRepoName(state.Name),
		message,
		FaintStyle.Render(state.Duration.Round(time.Millisecond).String()))
}

func (state *RepoOperationState) RenderFailed(message string) string {
	return fmt.Sprintf("%s %s - %s\n  %s\n",
		Cross,
		RenderRepoName(state.Name),
		message,
		ErrorStyle.Render(fmt.Sprintf("Error: %v", state.Err)))
}

type RepoBranchSpinnerView struct {
	Title          string
	States         []*RepoOperationState
	SkippedIndices map[int]bool
	LaunchOp       func(i int) tea.Cmd
	OnMsg          func(msg tea.Msg, states []*RepoOperationState) tea.Cmd
	RenderRepo     func(i int, state *RepoOperationState) string
}

func (cfg RepoBranchSpinnerView) Run() (failCount int, err error) {
	activeCount := len(cfg.States)
	for range cfg.SkippedIndices {
		activeCount--
	}
	p := tea.NewProgram(repoBranchSpinnerModel{
		totalRepos:     activeCount,
		startTime:      time.Now(),
		states:         cfg.States,
		skippedIndices: cfg.SkippedIndices,
		config:         cfg,
	})
	finalModel, err := p.Run()
	if err != nil {
		return 0, fmt.Errorf("error running program: %w", err)
	}
	if fm, ok := finalModel.(repoBranchSpinnerModel); ok {
		return fm.failCount, nil
	}
	return 0, nil
}

func (cfg RepoBranchSpinnerView) RunOrExit() {
	failCount, err := cfg.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if failCount > 0 {
		os.Exit(1)
	}
}

// Bubbletea model

type repoBranchSpinnerModel struct {
	doneCount      int
	failCount      int
	totalRepos     int
	startTime      time.Time
	states         []*RepoOperationState
	skippedIndices map[int]bool
	config         RepoBranchSpinnerView
}

func (m repoBranchSpinnerModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.states)*2)
	for i := range m.states {
		if m.skippedIndices[i] {
			continue
		}
		cmds = append(cmds, m.states[i].Spinner.Tick)
	}
	for i := range m.states {
		if m.skippedIndices[i] {
			continue
		}
		m.states[i].Status = StatusInProgress
		m.states[i].StartTime = time.Now()
		cmds = append(cmds, m.config.LaunchOp(i))
	}
	return tea.Batch(cmds...)
}

func (m repoBranchSpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RepoOperationDoneMsg:
		state := m.states[msg.RepoIndex]
		state.Duration = msg.Duration
		if msg.Err != nil {
			state.Status = StatusFailed
			state.Err = msg.Err
			m.failCount++
		} else {
			state.Status = StatusDone
		}
		m.doneCount++
		if m.doneCount >= m.totalRepos {
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmds []tea.Cmd
		for i := range m.states {
			if m.states[i].Status == StatusInProgress {
				var cmd tea.Cmd
				m.states[i].Spinner, cmd = m.states[i].Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)

	default:
		if m.config.OnMsg != nil {
			if cmd := m.config.OnMsg(msg, m.states); cmd != nil {
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m repoBranchSpinnerModel) View() string {
	var b strings.Builder

	// Header
	if m.doneCount < m.totalRepos {
		fmt.Fprintf(&b, "%s Progress: %d/%d complete\n",
			HeaderStyle.Render(m.config.Title+"..."),
			m.doneCount, m.totalRepos)
	} else {
		fmt.Fprintf(&b, "%s Completed in %s\n",
			HeaderStyle.Render("✓ "+m.config.Title+" complete!"),
			time.Since(m.startTime).Round(time.Millisecond))
	}

	// Repo lines
	for i, state := range m.states {
		fmt.Fprint(&b, m.config.RenderRepo(i, state))
	}

	// Failure summary
	if m.failCount > 0 {
		fmt.Fprint(&b, WarningStyle.Render(fmt.Sprintf("⚠ %d operation(s) failed", m.failCount))+"\n")
	}

	return b.String()
}
