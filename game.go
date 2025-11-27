package main

import (
	"slices"

	"github.com/andrewdaoust/scoundrel/deck"
)

func newDungeon() []deck.Card {
	d := deck.New(
		deck.Filter(func(c deck.Card) bool {
			if c.Suit == deck.Heart || c.Suit == deck.Diamond {
				if c.Rank == deck.Ace || c.Rank == deck.King || c.Rank == deck.Queen || c.Rank == deck.Jack {
					return true
				}
			}
			return false
		}),
		deck.Shuffle,
	)
	return d
}

func (m *model) drawToRoom(n int) {
	if len(m.dungeon) >= n {
		m.room = append(m.room, m.dungeon[:n]...)
		m.dungeon = m.dungeon[n:]
	} else {
		m.room = append(m.room, m.dungeon...)
		m.dungeon = []deck.Card{}
	}
}

func (m *model) usePotion(c deck.Card) {
	m.life = min(20, m.life+int(c.Rank))
}

func (m *model) equipWeapon(c deck.Card) {
	m.weapon = weapon{
		card:  c,
		slain: []deck.Card{},
	}
}

type attackType int

const (
	withFists attackType = iota
	withWeapon
)

func attackStrength(c deck.Card) int {
	rank := int(c.Rank)
	// Scale Ace to 14 for attacks
	if rank == 1 {
		rank = 14
	}
	return rank
}

func (m *model) attackWithFists(c deck.Card) {
	m.life = max(0, m.life-attackStrength(c))
}

func (m *model) canUseWeapon(c deck.Card) bool {
	// No weapon equipped
	if m.weapon.card.Rank == 0 {
		return false
	}

	// Weapon unused, can use any card
	if len(m.weapon.slain) == 0 {
		return true
	}

	// Check if card rank is less than or equal to last slain
	last := m.weapon.slain[len(m.weapon.slain)-1]
	return attackStrength(c) <= attackStrength(last)
}

func (m *model) attackWithWeapon(c deck.Card) {
	attack := max(0, attackStrength(c)-int(m.weapon.card.Rank))
	m.life = max(0, m.life-attack)
	m.weapon.slain = append(m.weapon.slain, c)
}

func (m *model) skipRoom() {
	m.selection = 0
	m.dungeon = append(m.dungeon, m.room...)
	m.room = []deck.Card{}
	m.skippable = false
	m.drawToRoom(4)
}

func (m *model) discard() {
	m.lastCard = m.room[m.selection]
	m.viewState = viewStateRoom
	m.room = slices.Delete(m.room, m.selection, m.selection+1)
	m.selection = 0
	m.attackTypeSelection = 1
	m.skippable = false

	if len(m.room) == 1 {
		m.drawToRoom(3)
		m.skippable = true
	}

	m.gameOverCheck()
}

func (m *model) chooseAttack() {
	c := m.room[m.selection]
	if m.canUseWeapon(c) {
		m.viewState = viewStateAttack
	} else {
		m.attackWithFists(c)
		m.discard()
	}
}

func (m *model) playAttack() {
	c := m.room[m.selection]
	switch m.attackTypeSelection {
	case int(withFists):
		m.attackWithFists(c)
	case int(withWeapon):
		m.attackWithWeapon(c)
	default:
		m.viewState = viewStateRoom
		return
	}
	m.discard()
}

func (m *model) playRoom() {
	if m.selection == len(m.room) {
		m.skipRoom()
		return
	}

	c := m.room[m.selection]
	switch c.Suit {
	case deck.Heart:
		m.usePotion(c)
	case deck.Diamond:
		m.equipWeapon(c)
	case deck.Spade, deck.Club:
		m.chooseAttack()
		return
	}
	m.discard()
}

func (m *model) gameOverCheck() {
	if m.life <= 0 || (len(m.dungeon) == 0 && len(m.room) == 0) {
		m.viewState = viewStateGameOver
	}
}

func (m model) score() int {
	score := m.life
	if len(m.dungeon) > 0 {
		for _, c := range append(m.dungeon, m.room...) {
			score -= attackStrength(c)
		}
		return score
	}

	if m.lastCard.Suit == deck.Heart {
		score += int(m.lastCard.Rank)
	}

	return score
}

// func (m *model) newGame() {
// 	newModel := initModel()
// 	*m = newModel
// }

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *model) up() {
	switch m.viewState {
	case viewStateRoom:
		maxSelections := len(m.room)
		if m.skippable {
			maxSelections += 1
		}
		m.selection = abs(m.selection - 1 + maxSelections) % maxSelections
	case viewStateAttack:
		m.attackTypeSelection = abs(m.attackTypeSelection - 1 + 3) % 3
	}
}

func (m *model) down() {
	switch m.viewState {
	case viewStateRoom:
		maxSelections := len(m.room)
		if m.skippable {
			maxSelections += 1
		}
		m.selection = abs(m.selection + 1) % maxSelections
	case viewStateAttack:
		m.attackTypeSelection = abs(m.attackTypeSelection + 1) % 3
	}
}
