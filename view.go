package main

import (
	"fmt"
	"strings"

	"github.com/andrewdaoust/scoundrel/deck"
)

type viewState string

const (
	viewStateRoom viewState = "room"
	// viewStateChooseAttack viewState = "choose"
	viewStateAttack   viewState = "attack"
	viewStateGameOver viewState = "gameover"
)

// layoutView creates a fixed layout with header, dynamic selection area, and footer
func layoutView(header string, selectionLines []string, footer string, width int, height int) string {
	if width <= 0 || height <= 0 {
		return header + strings.Join(selectionLines, "\n") + footer
	}

	var result []string

	// Calculate header position (centered)
	headerLine := strings.TrimSuffix(header, "\n\n")
	headerPadding := (width - len(headerLine)) / 2
	centeredHeader := strings.Repeat(" ", headerPadding) + headerLine

	// Add vertical padding to center the entire content block
	totalContentLines := 1 + len(selectionLines) + strings.Count(footer, "\n") + 1 // +1 for spacing
	topPadding := (height - totalContentLines) / 2

	// Add top padding
	for i := 0; i < topPadding; i++ {
		result = append(result, "")
	}

	// Add centered header
	result = append(result, centeredHeader)
	result = append(result, "") // Empty line after header

	// Add selection lines with 2-space offset from header position
	selectionPadding := headerPadding + 2
	for _, line := range selectionLines {
		if len(line) > 0 {
			result = append(result, strings.Repeat(" ", selectionPadding)+line)
		} else {
			result = append(result, "")
		}
	}

	// Add footer aligned with header
	footerLines := strings.Split(footer, "\n")
	for _, line := range footerLines {
		if len(line) > 0 {
			result = append(result, strings.Repeat(" ", headerPadding)+line)
		} else {
			result = append(result, "")
		}
	}

	return strings.Join(result, "\n")
}

// centerHorizontally centers text horizontally within the given width
func centerHorizontally(text string, width int) string {
	lines := strings.Split(text, "\n")
	var centeredLines []string

	for _, line := range lines {
		if len(line) >= width {
			centeredLines = append(centeredLines, line)
		} else {
			padding := (width - len(line)) / 2
			centeredLines = append(centeredLines, strings.Repeat(" ", padding)+line)
		}
	}

	return strings.Join(centeredLines, "\n")
}

// centerVertically adds vertical padding to center content
func centerVertically(text string, height int) string {
	lines := strings.Split(text, "\n")
	contentHeight := len(lines)

	if contentHeight >= height {
		return text
	}

	topPadding := (height - contentHeight) / 2
	var result []string

	// Add top padding
	for i := 0; i < topPadding; i++ {
		result = append(result, "")
	}

	// Add content
	result = append(result, lines...)

	return strings.Join(result, "\n")
}

func (m model) headerView() string {
	return fmt.Sprintf("❤️: %02d\tRemaining: %d\n\n", m.life, len(m.dungeon))
}

func (m model) footerView() string {
	s := ""

	if m.weapon.card.Rank != 0 {
		s += fmt.Sprintf("\n🗡  Power: %d", m.weapon.card.Rank)
	}

	if len(m.weapon.slain) > 0 {
		s += fmt.Sprintf(" (Last slain: %d)", attackStrength(m.weapon.slain[len(m.weapon.slain)-1]))
	}
	s += "\n\n\nPress q to quit."
	return s
}

func (m model) roomView() string {
	header := fmt.Sprintf("❤️: %02d\tRemaining: %d", m.life, len(m.dungeon))
	footer := m.footerView()

	var selectionLines []string

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
		selectionLines = append(selectionLines, fmt.Sprintf("%s %s%d", cursor, symbol, attackStrength(card)))
	}

	if m.skippable {
		cursor := " "
		if m.selection == len(m.room) {
			cursor = ">"
		}
		selectionLines = append(selectionLines, "")
		selectionLines = append(selectionLines, fmt.Sprintf("%s Skip this room", cursor))
	}

	return layoutView(header, selectionLines, footer, m.width, m.height)
}

func (m model) chooseAttackView() string {
	header := fmt.Sprintf("❤️: %02d\tRemaining: %d", m.life, len(m.dungeon))
	footer := m.footerView()

	cursor := map[bool]string{true: ">", false: " "}

	var selectionLines []string
	selectionLines = append(selectionLines, fmt.Sprintf("%s Fight with 👊", cursor[m.attackTypeSelection == 0]))
	selectionLines = append(selectionLines, fmt.Sprintf("%s Fight with 🗡️ %d", cursor[m.attackTypeSelection == 1], attackStrength(m.weapon.card)))
	selectionLines = append(selectionLines, "")
	selectionLines = append(selectionLines, fmt.Sprintf("%s Cancel", cursor[m.attackTypeSelection == 2]))

	return layoutView(header, selectionLines, footer, m.width, m.height)
}

func (m model) gameOverView() string {
	s := "💀 Game Over 💀\n\n"
	s += fmt.Sprintf("Score: %d\n\n", m.score())
	s += "Press enter to play again. Press q to quit."
	return s
}
