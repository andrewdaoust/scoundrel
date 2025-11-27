package main
import (
	"fmt"

	"github.com/andrewdaoust/scoundrel/deck"
)

type viewState string

const (
	viewStateRoom   viewState = "room"
	// viewStateChooseAttack viewState = "choose"
	viewStateAttack viewState = "attack"
	viewStateGameOver viewState = "gameover"
)

func (m model) headerView() string {
	return fmt.Sprintf("❤️: %02d\tRemaining: %d\n\n", m.life, len(m.dungeon))
}

func (m model) footerView() string {
	s := ""

	if m.weapon.card.Rank != 0 {
		s += fmt.Sprintf("\n🗡  Power: %d", m.weapon.card.Rank)
	}

	if len(m.weapon.slain) > 0 {
		s += fmt.Sprintf(" (Last slain: %d)\n", m.weapon.slain[len(m.weapon.slain)-1].Rank)
	}
	s += "\nPress q to quit."
	return s
}

func (m model) roomView() string {
	s := m.headerView()

	for i, card := range m.room {
		cursor := " "
		if m.selection == i {
			cursor = ">"
		}
		
		var symbol string
		switch card.Suit {
		case deck.Heart:
			symbol = "❤️"
		case deck.Diamond:
			symbol = "🗡️"
		default:
			symbol = "🐍"
		}
		s += fmt.Sprintf("%s %s%d\n", cursor, symbol, attackStrength(card))
	}

	if m.skippable {
		cursor := " "
		if m.selection == len(m.room) {
			cursor = ">"
		}
		s += fmt.Sprintf("\n%s Skip this room\n", cursor)
	}

	s += m.footerView()
	return s
}

func (m model) chooseAttackView() string {
	s := m.headerView()

	cursor := map[bool]string{true: ">", false: " "}

	s += fmt.Sprintf("%s Fight with 👊\n", cursor[m.attackTypeSelection == 0])
	s += fmt.Sprintf("%s Fight with 🗡️ %d\n", cursor[m.attackTypeSelection == 1], attackStrength(m.weapon.card))
	s += fmt.Sprintf("\n%s Cancel\n", cursor[m.attackTypeSelection == 2])

	s += m.footerView()
	return s
}

func (m model) gameOverView() string {
	s := "💀 Game Over 💀\n\n"
	s += fmt.Sprintf("Score: %d\n\n", m.score())
	s += "Press enter to play again. Press q to quit."
	return s
}
