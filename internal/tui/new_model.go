package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepPurpose step = iota
	stepPass
	stepConfirm
	stepDone
)

type Model struct {
	name    string
	comment string
	keyPath string
	pubPath string
	ttl     string

	purpose textinput.Model
	pass1   textinput.Model
	pass2   textinput.Model

	cur     step
	aborted bool
}

func NewModel(name, comment, keyPath, pubPath, ttl string) Model {
	p := textinput.New()
	p.Placeholder = "work"
	p.CharLimit = 64
	p.SetValue(name)
	pw1 := textinput.New()
	pw1.Placeholder = "passphrase"
	pw1.EchoMode = textinput.EchoPassword
	pw1.EchoCharacter = '•'
	pw2 := textinput.New()
	pw2.Placeholder = "confirm"
	pw2.EchoMode = textinput.EchoPassword
	pw2.EchoCharacter = '•'
	return Model{name: name, comment: comment, keyPath: keyPath, pubPath: pubPath, ttl: ttl, purpose: p, pass1: pw1, pass2: pw2, cur: stepPurpose}
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) View() string {
	s := lipgloss.NewStyle().Margin(1, 2)
	switch m.cur {
	case stepPurpose:
		return s.Render(fmt.Sprintf("Create new SSH key\n\nPurpose name: %s\n\n(enter to continue, ctrl+c to abort)", m.purpose.View()))
	case stepPass:
		return s.Render(fmt.Sprintf("Set passphrase for %q\n\n%s\n%s\n\n(enter to continue, ctrl+c to abort)", m.purpose.Value(), m.pass1.View(), m.pass2.View()))
	case stepConfirm:
		return s.Render(fmt.Sprintf("About to create key:\n  Name: %s\n  Comment: %s\n  Path: %s\n  TTL: %s\n\n(enter to create, esc to go back)", m.name, m.comment, m.keyPath, displayTTL(m.ttl)))
	case stepDone:
		return s.Render("Done.\n")
	}
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			switch m.cur {
			case stepPurpose:
				val := strings.TrimSpace(m.purpose.Value())
				if val == "" {
					return m, nil
				}
				m.name = val
				m.cur = stepPass
				return m, nil
			case stepPass:
				if m.pass1.Value() == "" || m.pass1.Value() != m.pass2.Value() {
					return m, nil
				}
				m.cur = stepConfirm
				return m, nil
			case stepConfirm:
				m.cur = stepDone
				return m, tea.Quit
			}
		case "esc":
			if m.cur == stepConfirm {
				m.cur = stepPass
				return m, nil
			}
		}
	}
	// forward to active input
	switch m.cur {
	case stepPurpose:
		var cmd tea.Cmd
		m.purpose, cmd = m.purpose.Update(msg)
		return m, cmd
	case stepPass:
		var cmd tea.Cmd
		m.pass1, cmd = m.pass1.Update(msg)
		m.pass2, _ = m.pass2.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) Aborted() bool      { return m.aborted }
func (m *Model) Passphrase() string { return m.pass1.Value() }
func (m *Model) ZeroPass()          { m.pass1.SetValue(""); m.pass2.SetValue("") }

func displayTTL(ttl string) string {
	if ttl == "0" {
		return "none"
	}
	return ttl
}
