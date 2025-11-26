package internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/muesli/reflow/truncate"
	"golang.org/x/term"
)

func getTerminalWidth() int {
    width, _, err := term.GetSize(0)
    if err != nil { return 80 }
    return width
}

type MultiSpinner struct {
    spinners []*LineSpinner
    mu       sync.Mutex
    done     chan struct{}
    frames   []string
	onClose  []func()
}

type LineSpinner struct {
    text   string
    done   bool
    result string
    line   int

}

func NewMultiSpinner() *MultiSpinner { 
	return &MultiSpinner{
		done: make(chan struct{}),
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (m *MultiSpinner) Add(text string) *LineSpinner {
    ls := &LineSpinner{ text: text, line: len(m.spinners) }
    m.spinners = append(m.spinners, ls)
    return ls
}

func (m *MultiSpinner) UpdateText(ls *LineSpinner, text string) {
	m.mu.Lock()
	ls.text = text
	m.mu.Unlock()
}

func (m *MultiSpinner) print(frame int) {
	termWidth := getTerminalWidth()

	m.mu.Lock()
    fmt.Printf("\033[%dA", len(m.spinners))

    for _, s := range m.spinners {
        var line string
		text := s.text
		truncatedText := truncate.String(text, uint(termWidth))
		if tl := len(truncatedText); tl < len(text) {
			text = truncatedText[:tl-6] + "..."
		}

        if s.done { line = strings.Join([]string{text, s.result}, " ")
        } else { line = fmt.Sprintf("%s %s", text, m.frames[frame%len(m.frames)]) }


        fmt.Printf("\r%s\033[K\n", line)
    }
	m.mu.Unlock()
}

func (m *MultiSpinner) Start() {
    for range m.spinners { fmt.Println("") }

    go func() {
        for i := 0; true; i++ {
            select {
            case <-m.done: return
            default:
				m.print(i)
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
}

func (m *MultiSpinner) Stop(ls *LineSpinner, message string) {
    m.mu.Lock()
    ls.done = true
    ls.result = message
    m.mu.Unlock()
}

func (m *MultiSpinner) Done(ls *LineSpinner) {
	m.mu.Lock()
	ls.done = true
	ls.result = color.New(color.FgGreen, color.Bold).Sprint("✓")
	m.mu.Unlock()
}

func (m *MultiSpinner) Fail(ls *LineSpinner) {
	m.mu.Lock()
	ls.done = true
	ls.result = color.New(color.FgRed, color.Bold).Sprint("✗")
	m.mu.Unlock()
}

func (m *MultiSpinner) AddOnClose(f func()) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.onClose = append(m.onClose, f)
}

func (m *MultiSpinner) Close() {
	close(m.done)
	m.print(0)
    for _, f := range m.onClose {
        f()
    }
}
