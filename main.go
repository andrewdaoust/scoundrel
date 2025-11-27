package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/andrewdaoust/scoundrel/deck"
)

type model struct {
	dungeon   []deck.Card
	room      []deck.Card
	life      int
	weapon    weapon
	skippable bool
	lastCard  deck.Card

	selection           int
	attackTypeSelection int
	viewState           viewState

	// Terminal dimensions
	width  int
	height int
}

type weapon struct {
	card  deck.Card
	slain []deck.Card
}

func initModel() model {
	m := model{
		dungeon: newDungeon(),
		room:    []deck.Card{},
		life:    20,
		weapon: weapon{
			card:  deck.Card{Rank: 0},
			slain: []deck.Card{},
		},
		skippable: true,

		selection:           0,
		attackTypeSelection: 1,
		viewState:           viewStateRoom,
	}

	m.drawToRoom(4)

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle window size changes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			m.up()
		case "down", "j":
			m.down()
		case "enter", "return":
			// Handle selection based on current view state
			switch m.viewState {
			case viewStateAttack:
				m.playAttack()
			case viewStateRoom:
				m.playRoom()
			case viewStateGameOver:
				m = initModel()
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	switch m.viewState {
	case viewStateRoom:
		return m.roomView()
	case viewStateAttack:
		return m.chooseAttackView()
	case viewStateGameOver:
		content := m.gameOverView()
		// Center the game over content if we have terminal dimensions
		if m.width > 0 && m.height > 0 {
			content = centerHorizontally(content, m.width)
			content = centerVertically(content, m.height)
		}
		return content
	default:
		return "Unknown view state"
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
