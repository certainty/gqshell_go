package introspector

import (
	"github.com/certainty/gqshell_go/internal/application/tui/style"
	tea "github.com/charmbracelet/bubbletea"
)

type Bubble struct {
	activeBox    int
	height       int
	heightMargin int
	width        int
	widthMargin  int
	style        *style.Styles
	boxes        []tea.Model
}

func NewBubble(height, wm, width, hm int, style *style.Styles) *Bubble {
	b := &Bubble{
		activeBox:    0,
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        style,
		boxes:        make([]tea.Model, 4),
	}
	// initialise the other components
	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) headerView() string {
	// TODO better header, tabs?
	return ""
}

func (b *Bubble) View() string {
	header := b.headerView()
	return header + b.boxes[b.activeBox].View()
}
