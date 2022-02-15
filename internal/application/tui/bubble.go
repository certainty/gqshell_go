package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/certainty/gqshell_go/internal/application/tui/bubbles/selection"
	"github.com/certainty/gqshell_go/internal/application/tui/style"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	startState appState = iota
	errorState
	loadedState
	quittingState
	quitState
)

type Config struct {
	endpoints []EndpointConfig
}

type EndpointConfig struct {
	name string
}

type MenuEntry struct {
	endpointName string
}

type appState int

// the main application state
type Bubble struct {
	state         appState
	width         int
	height        int
	styles        *style.Styles
	error         string
	lastResize    tea.WindowSizeMsg
	menuEntries   []MenuEntry
	boxes         []tea.Model
	activeBox     int
	activeEndoint int
}

func NewBubble(cfg Config) *Bubble {
	width, height, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)

	}

	b := &Bubble{
		styles:      style.DefaultStyles(),
		menuEntries: make([]MenuEntry, len(cfg.endpoints)),
		boxes:       make([]tea.Model, 2),
		height:      height,
		width:       width,
	}

	for i, endpoint := range cfg.endpoints {
		b.menuEntries[i] = MenuEntry{
			endpointName: endpoint.name,
		}
	}

	b.activeEndoint = 0
	b.activeBox = 0
	//b.state = startState
	b.state = loadedState
	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		case "tab", "shift+tab":
			b.activeBox = (b.activeBox + 1) % 2
		}
	case tea.WindowSizeMsg:
		b.lastResize = msg
		b.width = msg.Width
		b.height = msg.Height
		if b.state == loadedState {
			for i, bx := range b.boxes {
				m, cmd := bx.Update(msg)
				b.boxes[i] = m
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	case selection.SelectedMsg:
		b.activeBox = 1
	case selection.ActiveMsg:
		cmds = append(cmds, func() tea.Msg {
			return b.lastResize
		})
	}
	if b.state == loadedState {
		ab, cmd := b.boxes[b.activeBox].Update(msg)
		b.boxes[b.activeBox] = ab
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return b, tea.Batch(cmds...)
}

func (b Bubble) headerView() string {
	w := b.width - b.styles.App.GetHorizontalFrameSize()
	name := "GraphqShell"
	return b.styles.Header.Copy().Width(w).Render(name)
}

func (b *Bubble) viewForBox(i int) string {
	isActive := i == b.activeBox

	switch box := b.boxes[i].(type) {
	case *selection.Bubble:
		var s lipgloss.Style
		s = b.styles.Menu
		if isActive {
			s = s.Copy().BorderForeground(b.styles.ActiveBorderColor)
		}
		return s.Render(box.View())
	default:
		panic(fmt.Sprintf("unknown box type %T", box))
	}
}

func (b Bubble) errorView() string {
	s := b.styles

	errorMessage := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.ErrorTitle.Render("Error"),
		s.ErrorBody.Render(b.error),
	)

	h := b.height -
		s.App.GetVerticalFrameSize() -
		lipgloss.Height(b.headerView()) -
		lipgloss.Height(b.footerView()) -
		s.EndpointBody.GetVerticalFrameSize() + 3

	return s.Error.Copy().Height(h).Render(errorMessage)
}

func (b Bubble) View() string {
	s := strings.Builder{}
	s.WriteString(b.headerView())
	s.WriteRune('\n')

	switch b.state {
	case loadedState:
		lb := b.viewForBox(0)
		rb := b.viewForBox(1)
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lb, rb))
	case errorState:
		s.WriteString(b.errorView())
	}

	s.WriteRune('\n')
	s.WriteString(b.footerView())
	return b.styles.App.Render(s.String())
}

func (self Bubble) footerView() string {
	renderedHelp := &strings.Builder{}
	var helpEntries []selection.HelpEntry
	if self.state != errorState {
		helpEntries = []selection.HelpEntry{
			{Key: "tab", Value: "section"},
		}
		if box, ok := self.boxes[self.activeBox].(selection.HelpableBubble); ok {
			help := box.Help()
			for _, helpEntry := range help {
				helpEntries = append(helpEntries, helpEntry)
			}
		}
	}

	helpEntries = append(helpEntries, selection.HelpEntry{Key: "q", Value: "quit"})
	for i, helpEntry := range helpEntries {
		fmt.Fprint(renderedHelp, helpEntryRender(helpEntry, self.styles))
		if i != len(helpEntries)-1 {
			fmt.Fprint(renderedHelp, self.styles.HelpDivider)
		}
	}

	help := renderedHelp.String()
	section := self.styles.FooterSection.Render("TestSection")
	gap := lipgloss.NewStyle().Width(self.width - lipgloss.Width(help) - lipgloss.Width(section) - self.styles.App.GetHorizontalFrameSize()).Render("")

	footer := lipgloss.JoinHorizontal(lipgloss.Top, help, gap, section)
	return self.styles.Footer.Render(footer)
}

func helpEntryRender(h selection.HelpEntry, s *style.Styles) string {
	return fmt.Sprintf("%s %s", s.HelpKey.Render(h.Key), s.HelpValue.Render(h.Value))
}

func Start() {
	endpoints := make([]EndpointConfig, 2)
	endpoints[0] = EndpointConfig{name: "local"}
	endpoints[1] = EndpointConfig{name: "somehost"}

	cfg := Config{
		endpoints: endpoints,
	}

	program := tea.NewProgram(NewBubble(cfg), tea.WithAltScreen(), tea.WithoutCatchPanics())

	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
