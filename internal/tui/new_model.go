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
	stepPass1
	stepPass2
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

	cur       step
	aborted   bool
	confirmed bool
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	hintStyle  = lipgloss.NewStyle().Faint(true)
	frame      = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))
)

func NewModel(name, comment, keyPath, pubPath, ttl string) *Model {
	p := textinput.New()
	p.Placeholder = "work"
	p.CharLimit = 64
	p.SetValue(name)

	pw1 := textinput.New()
	pw1.Placeholder = "passphrase"
	pw1.EchoMode = textinput.EchoPassword
	pw1.EchoCharacter = '•'

	pw2 := textinput.New()
	pw2.Placeholder = "confirm passphrase"
	pw2.EchoMode = textinput.EchoPassword
	pw2.EchoCharacter = '•'

	m := &Model{
		name: name, comment: comment, keyPath: keyPath, pubPath: pubPath, ttl: ttl,
		purpose: p, pass1: pw1, pass2: pw2,
	}
	if name == "" {
		m.cur = stepPurpose
		m.purpose.Focus()
	} else {
		m.cur = stepPass1
		m.pass1.Focus()
	}
	return m
}

func (m *Model) Init() tea.Cmd { return textinput.Blink }

func (m *Model) View() string {
	switch m.cur {
	case stepPurpose:
		return frame.Render(
			titleStyle.Render("Create new SSH key") + "\n\n" +
				"Purpose name:\n" + m.purpose.View() + "\n\n" +
				hintStyle.Render("(Enter to continue · Ctrl+C to cancel)"),
		)
	case stepPass1:
		return frame.Render(
			titleStyle.Render(fmt.Sprintf("Set passphrase for %q", strings.TrimSpace(m.purpose.Value()))) + "\n\n" +
				"Passphrase:\n" + m.pass1.View() + "\n\n" +
				hintStyle.Render("(Enter to confirm · Ctrl+C to cancel)"),
		)
	case stepPass2:
		return frame.Render(
			titleStyle.Render("Confirm passphrase") + "\n\n" +
				"Confirm passphrase:\n" + m.pass2.View() + "\n\n" +
				hintStyle.Render("(Enter to continue · Esc to go back · Ctrl+C to cancel)"),
		)
	case stepConfirm:
		return frame.Render(
			titleStyle.Render("Confirm") + "\n\n" +
				fmt.Sprintf("Name: %s\nComment: %s\nPath: %s\nTTL: %s\n\n",
					m.name, m.comment, m.keyPath, displayTTL(m.ttl)) +
				hintStyle.Render("(Enter to create · Esc to edit · Ctrl+C to cancel)"),
		)
	case stepDone:
		return frame.Render(titleStyle.Render("Done") + "\n")
	}
	return ""
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "esc":
			switch m.cur {
			case stepPass2:
				m.cur = stepPass1
				m.pass2.Blur()
				m.pass1.Focus()
			case stepConfirm:
				m.cur = stepPass2
			}
			return m, nil
		case "enter":
			switch m.cur {
			case stepPurpose:
				val := strings.TrimSpace(m.purpose.Value())
				if val == "" {
					return m, nil
				}
				m.name = val
				m.purpose.Blur()
				m.pass1.Focus()
				m.cur = stepPass1
				return m, nil
			case stepPass1:
				if m.pass1.Value() == "" {
					return m, nil
				}
				m.pass1.Blur()
				m.pass2.Focus()
				m.cur = stepPass2
				return m, nil
			case stepPass2:
				if m.pass2.Value() == "" || m.pass2.Value() != m.pass1.Value() {
					m.pass2.SetValue("")
					return m, nil
				}
				m.cur = stepConfirm
				return m, nil
			case stepConfirm:
				// ✅ set confirmed before quitting
				m.confirmed = true
				m.cur = stepDone
				return m, tea.Quit
			}
		}
	}

	// forward to active input
	switch m.cur {
	case stepPurpose:
		var cmd tea.Cmd
		m.purpose, cmd = m.purpose.Update(msg)
		return m, cmd
	case stepPass1:
		var cmd tea.Cmd
		m.pass1, cmd = m.pass1.Update(msg)
		return m, cmd
	case stepPass2:
		var cmd tea.Cmd
		m.pass2, cmd = m.pass2.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) Aborted() bool      { return m.aborted }
func (m *Model) Confirmed() bool    { return m.confirmed }
func (m *Model) Passphrase() string { return m.pass1.Value() }
func (m *Model) ZeroPass()          { m.pass1.SetValue(""); m.pass2.SetValue("") }
func displayTTL(ttl string) string {
	if ttl == "0" {
		return "none"
	}
	return ttl
}
