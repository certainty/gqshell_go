package tui

import (
	"fmt"
	"github.com/certainty/gqshell_go/internal/application/tui/bubbles/selection"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
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
	b := &Bubble{
		styles:      style.DefaultStyles(),
		menuEntries: make([]MenuEntry, len(cfg.endpoints)),
		boxes:       make([]tea.Model, 2),
		height:      20,
		width:       20,
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
	name := ""
	return b.styles.Header.Copy().Width(w).Render(name)
}

func (b *Bubble) viewForBox(i int) string {
	isActive := i == b.activeBox
	switch box := b.boxes[i].(type) {
	case *selection.Bubble:
		// Menu
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
	str := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.ErrorTitle.Render("Bummer"),
		s.ErrorBody.Render(b.error),
	)
	h := b.height -
		s.App.GetVerticalFrameSize() -
		lipgloss.Height(b.headerView()) -
		s.RepoBody.GetVerticalFrameSize() +
		3 // TODO: this is repo header height -- get it dynamically
	return s.Error.Copy().Height(h).Render(str)
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
	return b.styles.App.Render(s.String())
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

	p := tea.NewProgram(NewBubble(cfg), tea.WithAltScreen(), tea.WithoutCatchPanics())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
