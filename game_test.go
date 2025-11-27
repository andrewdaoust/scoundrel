package main

import (
	"testing"

	"github.com/andrewdaoust/scoundrel/deck"
)

func TestNewDungeon(t *testing.T) {
	d := newDungeon()
	assertExpectedDungeonLength(t, len(d), 52-8)

	for _, c := range d {
		if (c.Suit == deck.Heart || c.Suit == deck.Diamond) && (c.Rank == deck.Ace || c.Rank == deck.King || c.Rank == deck.Queen || c.Rank == deck.Jack) {
			t.Errorf("dungeon contains invalid card: %s", c.String())
		}
	}
}

func TestDrawToRoom(t *testing.T) {
	tests := []struct {
		m                  model
		n                  int
		expectedDungeonLen int
		expectedRoomLen    int
	}{
		{
			m: model{
				dungeon: testDungeon(),
				room:    []deck.Card{},
			},
			n:                  4,
			expectedDungeonLen: len(testDungeon()) - 4,
			expectedRoomLen:    4,
		},
		{
			m: model{
				dungeon: testDungeon(),
				room:    []deck.Card{{Suit: deck.Heart, Rank: 5}},
			},
			n:                  3,
			expectedDungeonLen: len(testDungeon()) - 3,
			expectedRoomLen:    4,
		},
		{
			m: model{
				dungeon: []deck.Card{
					{Suit: deck.Spade, Rank: 2},
					{Suit: deck.Spade, Rank: 3},
				},
				room: []deck.Card{{Suit: deck.Heart, Rank: 5}},
			},
			n:                  3,
			expectedDungeonLen: 0,
			expectedRoomLen:    3,
		},
	}

	for _, tt := range tests {
		tt.m.drawToRoom(tt.n)
		assertExpectedRoomLength(t, len(tt.m.room), tt.expectedRoomLen)
		assertExpectedDungeonLength(t, len(tt.m.dungeon), tt.expectedDungeonLen)
	}
}

func TestUsePotion(t *testing.T) {
	tests := []struct {
		initialLife int
		potion      deck.Card
		expected    int
	}{
		{10, deck.Card{Suit: deck.Heart, Rank: 5}, 15},
		{18, deck.Card{Suit: deck.Heart, Rank: 5}, 20},
		{20, deck.Card{Suit: deck.Heart, Rank: 5}, 20},
	}

	for _, test := range tests {
		m := model{life: test.initialLife}
		m.usePotion(test.potion)
		assertExpectedLife(t, m.life, test.expected)
	}
}

func TestEquipWeapon(t *testing.T) {
	tests := []struct {
		m          model
		weaponCard deck.Card
	}{
		{
			m:          model{},
			weaponCard: deck.Card{Suit: deck.Spade, Rank: 10},
		},
		{
			m: model{
				weapon: weapon{
					card:  deck.Card{Suit: deck.Heart, Rank: 3},
					slain: []deck.Card{{Suit: deck.Diamond, Rank: 5}},
				},
			},
			weaponCard: deck.Card{Suit: deck.Club, Rank: 7},
		},
	}

	for _, test := range tests {
		test.m.equipWeapon(test.weaponCard)
		if test.m.weapon.card != test.weaponCard {
			t.Errorf("expected weapon card to be %s, got %s", test.weaponCard.String(), test.m.weapon.card.String())
		}
		if len(test.m.weapon.slain) != 0 {
			t.Errorf("expected slain list to be empty, got length %d", len(test.m.weapon.slain))
		}
	}
}

func TestAttackStrength(t *testing.T) {
	tests := []struct {
		card     deck.Card
		expected int
	}{
		{deck.Card{Rank: 2}, 2},
		{deck.Card{Rank: 3}, 3},
		{deck.Card{Rank: 4}, 4},
		{deck.Card{Rank: 5}, 5},
		{deck.Card{Rank: 6}, 6},
		{deck.Card{Rank: 7}, 7},
		{deck.Card{Rank: 8}, 8},
		{deck.Card{Rank: 9}, 9},
		{deck.Card{Rank: 10}, 10},
		{deck.Card{Rank: deck.Jack}, 11},
		{deck.Card{Rank: deck.Queen}, 12},
		{deck.Card{Rank: deck.King}, 13},
		{deck.Card{Rank: deck.Ace}, 14},
	}

	for _, test := range tests {
		result := attackStrength(test.card)
		if result != test.expected {
			t.Errorf("expected attack strength of %s to be %d, got %d", test.card.String(), test.expected, result)
		}
	}
}

func TestAttackWithFists(t *testing.T) {
	tests := []struct {
		m            model
		c            deck.Card
		expectedLife int
	}{
		{model{life: 20}, deck.Card{Rank: 5}, 15},
		{model{life: 10}, deck.Card{Rank: 3}, 7},
		{model{life: 4}, deck.Card{Rank: 10}, 0},
		{model{life: 15}, deck.Card{Rank: deck.Ace}, 1},
	}

	for _, test := range tests {
		test.m.attackWithFists(test.c)
		assertExpectedLife(t, test.m.life, test.expectedLife)
	}
}

func TestCanUseWeapon(t *testing.T) {
	tests := []struct {
		m        model
		c        deck.Card
		expected bool
	}{
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 0}, slain: []deck.Card{}}},
			c:        deck.Card{Rank: 5},
			expected: false,
		},
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{}}},
			c:        deck.Card{Rank: 5},
			expected: true,
		},
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{}}},
			c:        deck.Card{Rank: deck.Queen},
			expected: true,
		},
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{{Rank: 3}}}},
			c:        deck.Card{Rank: 5},
			expected: false,
		},
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{{Rank: 3}}}},
			c:        deck.Card{Rank: 2},
			expected: true,
		},
		{
			m:        model{weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{{Rank: 3}}}},
			c:        deck.Card{Rank: deck.Ace},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.m.canUseWeapon(test.c)
		if result != test.expected {
			t.Errorf("expected canUseWeapon with card %s to be %t, got %t", test.c.String(), test.expected, result)
		}
	}
}

func TestAttackWithWeapon(t *testing.T) {
	tests := []struct {
		m                model
		c                deck.Card
		expectedLife     int
		expectedSlainLen int
	}{
		{
			m:                model{life: 20, weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{}}},
			c:                deck.Card{Rank: 5},
			expectedLife:     20,
			expectedSlainLen: 1,
		},
		{
			m:                model{life: 20, weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{}}},
			c:                deck.Card{Rank: deck.Queen},
			expectedLife:     18,
			expectedSlainLen: 1,
		},
		{
			m:                model{life: 20, weapon: weapon{card: deck.Card{Rank: 10}, slain: []deck.Card{}}},
			c:                deck.Card{Rank: 10},
			expectedLife:     20,
			expectedSlainLen: 1,
		},
		{
			m:                model{life: 2, weapon: weapon{card: deck.Card{Rank: 5}, slain: []deck.Card{}}},
			c:                deck.Card{Rank: 10},
			expectedLife:     0,
			expectedSlainLen: 1,
		},
		{
			m:                model{life: 10, weapon: weapon{card: deck.Card{Rank: 5}, slain: []deck.Card{{Rank: 4}}}},
			c:                deck.Card{Rank: 3},
			expectedLife:     10,
			expectedSlainLen: 2,
		},
	}

	for _, test := range tests {
		test.m.attackWithWeapon(test.c)
		assertExpectedLife(t, test.m.life, test.expectedLife)
		if len(test.m.weapon.slain) != test.expectedSlainLen {
			t.Errorf("expected slain length to be %d after attack with %s, got %d", test.expectedSlainLen, test.c.String(), len(test.m.weapon.slain))
		}
		if test.m.weapon.slain[len(test.m.weapon.slain)-1] != test.c {
			t.Errorf("expected last slain card to be %s, got %s", test.c.String(), test.m.weapon.slain[len(test.m.weapon.slain)-1].String())
		}

	}
}

func TestSkipRoom(t *testing.T) {
	tests := []model{
		{
			dungeon:   testDungeon(),
			room:      testRoom(4),
			selection: 1,
			skippable: true,
		},
		{
			dungeon:   testDungeon(),
			room:      testRoom(4),
			selection: 2,
			skippable: true,
		},
		{
			dungeon:   testDungeon(),
			room:      testRoom(4),
			selection: 1,
			skippable: true,
		},
	}

	for _, test := range tests {
		prevDungeonLen := len(test.dungeon)
		test.skipRoom()
		if len(test.room) != 4 {
			t.Errorf("expected room length to be 4 after skip, got %d", len(test.room))
		}
		if len(test.dungeon) != prevDungeonLen {
			t.Errorf("expected dungeon length to be %d after skip, got %d", prevDungeonLen, len(test.dungeon))
		}
		if test.selection != 0 {
			t.Errorf("expected selection to be 0 after skip, got %d", test.selection)
		}
		assertExpectedSkippable(t, test.skippable, false)
	}
}

func TestDiscard(t *testing.T) {
	tests := []struct {
		m                  model
		expectedDungeonLen int
		expectedRoomLen    int
		expectedSkippable  bool
		expectedViewState  viewState
		expectedLastCard   deck.Card
	}{
		{
			m: model{
				life:      20,
				dungeon:   testDungeon(),
				room:      testRoom(4),
				selection: 1,
				skippable: true,
				viewState: viewStateAttack,
			},
			expectedDungeonLen: len(testDungeon()),
			expectedRoomLen:    3,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
			expectedLastCard:   testRoom(4)[1],
		},
		{
			m: model{
				life:      20,
				dungeon:   testDungeon(),
				room:      testRoom(3),
				selection: 2,
				skippable: false,
				viewState: viewStateRoom,
			},
			expectedDungeonLen: len(testDungeon()),
			expectedRoomLen:    2,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
			expectedLastCard:   testRoom(3)[2],
		},
		{
			m: model{
				life:      20,
				dungeon:   testDungeon(),
				room:      testRoom(2),
				selection: 1,
				skippable: false,
				viewState: viewStateAttack,
			},
			expectedDungeonLen: len(testDungeon()) - 3,
			expectedRoomLen:    4,
			expectedSkippable:  true,
			expectedViewState:  viewStateRoom,
			expectedLastCard:   testRoom(2)[1],
		},
	}

	for _, tt := range tests {
		selectedCard := tt.m.room[tt.m.selection]
		tt.m.discard()
		if len(tt.m.room) != tt.expectedRoomLen {
			t.Errorf("expected room length to be %d after discard, got %d", tt.expectedRoomLen, len(tt.m.room))
		}
		if len(tt.m.dungeon) != tt.expectedDungeonLen {
			t.Errorf("expected dungeon length to be %d after discard, got %d", tt.expectedDungeonLen, len(tt.m.dungeon))
		}
		if tt.m.selection != 0 {
			t.Errorf("expected selection to be 0 after discard, got %d", tt.m.selection)
		}
		assertExpectedSkippable(t, tt.m.skippable, tt.expectedSkippable)
		if tt.m.viewState != tt.expectedViewState {
			t.Errorf("expected viewState to be %s after discard, got %s", tt.expectedViewState, tt.m.viewState)
		}
		assertLastCard(t, tt.m.lastCard, tt.expectedLastCard)
		for _, c := range tt.m.room {
			if c.Rank == selectedCard.Rank && c.Suit == selectedCard.Suit {
				t.Errorf("expected discarded card %s to not be in room after discard", selectedCard.String())
			}
		}
	}
}

func TestChooseAttack(t *testing.T) {
	tests := []struct {
		m                  model
		expectedViewState  viewState
		expectedRoomLen    int
		expectedDungeonLen int
		expectedLife       int
		expectedSelection  int
		expectedSkippable  bool
	}{
		{
			m: model{
				life:    15,
				dungeon: testDungeon(),
				room: []deck.Card{
					{Suit: deck.Heart, Rank: 5},
					{Suit: deck.Diamond, Rank: 7},
					{Suit: deck.Club, Rank: 7},
					{Suit: deck.Club, Rank: 6},
				},
				selection: 2,
				skippable: true,
				viewState: viewStateRoom,
			},
			expectedViewState:  viewStateRoom,
			expectedRoomLen:    3,
			expectedDungeonLen: len(testDungeon()),
			expectedLife:       8,
			expectedSelection:  0,
			expectedSkippable:  false,
		},
	}

	for _, tt := range tests {
		tt.m.chooseAttack()
		if tt.m.viewState != tt.expectedViewState {
			t.Errorf("expected viewState to be %s after chooseAttack, got %s", tt.expectedViewState, tt.m.viewState)
		}
		if len(tt.m.room) != tt.expectedRoomLen {
			t.Errorf("expected room length to be %d after chooseAttack, got %d", tt.expectedRoomLen, len(tt.m.room))
		}
		if len(tt.m.dungeon) != tt.expectedDungeonLen {
			t.Errorf("expected dungeon length to be %d after chooseAttack, got %d", tt.expectedDungeonLen, len(tt.m.dungeon))
		}
		assertExpectedLife(t, tt.m.life, tt.expectedLife)
		if tt.m.selection != tt.expectedSelection {
			t.Errorf("expected selection to be %d after chooseAttack, got %d", tt.expectedSelection, tt.m.selection)
		}
		assertExpectedSkippable(t, tt.m.skippable, tt.expectedSkippable)
	}
}

func TestPlayAttack(t *testing.T) {}

func TestPlayRoom(t *testing.T) {
	tests := []struct {
		m                  model
		expectedLife       int
		expectedRoomLen    int
		expectedDungeonLen int
		expectedWeaponRank deck.Rank
		expectedSkippable  bool
		expectedViewState  viewState
	}{
		{ // Test using a potion
			m: model{
				life:    15,
				dungeon: testDungeon(),
				room: []deck.Card{
					{Suit: deck.Heart, Rank: 5},
					{Suit: deck.Diamond, Rank: 7},
					{Suit: deck.Club, Rank: 7},
					{Suit: deck.Club, Rank: 6},
				},
				selection: 0,
				skippable: true,
				viewState: viewStateRoom,
			},
			expectedLife:       20,
			expectedRoomLen:    3,
			expectedDungeonLen: len(testDungeon()),
			expectedWeaponRank: 0,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
		},
		{ // Test equipping a weapon
			m: model{
				life:    15,
				dungeon: testDungeon(),
				room: []deck.Card{
					{Suit: deck.Heart, Rank: 5},
					{Suit: deck.Diamond, Rank: 7},
					{Suit: deck.Club, Rank: 7},
					{Suit: deck.Club, Rank: 6},
				},
				selection: 1,
				skippable: true,
				viewState: viewStateRoom,
			},
			expectedLife:       15,
			expectedRoomLen:    3,
			expectedDungeonLen: len(testDungeon()),
			expectedWeaponRank: 7,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
		},
		{ // Test skipping a room
			m: model{
				life:    15,
				dungeon: testDungeon(),
				room: []deck.Card{
					{Suit: deck.Heart, Rank: 5},
					{Suit: deck.Diamond, Rank: 7},
					{Suit: deck.Club, Rank: 7},
					{Suit: deck.Club, Rank: 6},
				},
				selection: 4,
				skippable: true,
				viewState: viewStateRoom,
			},
			expectedLife:       15,
			expectedRoomLen:    4,
			expectedDungeonLen: len(testDungeon()),
			expectedWeaponRank: 0,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
		},
		{ // Test attack with no weapon equipped
			m: model{
				life:    15,
				dungeon: testDungeon(),
				room: []deck.Card{
					{Suit: deck.Heart, Rank: 5},
					{Suit: deck.Diamond, Rank: 7},
					{Suit: deck.Club, Rank: 7},
					{Suit: deck.Club, Rank: 6},
				},
				selection: 2,
				skippable: true,
				viewState: viewStateRoom,
			},
			expectedLife:       8,
			expectedRoomLen:    3,
			expectedDungeonLen: len(testDungeon()),
			expectedWeaponRank: 0,
			expectedSkippable:  false,
			expectedViewState:  viewStateRoom,
		},
	}

	for _, tt := range tests {
		tt.m.playRoom()
		assertExpectedLife(t, tt.m.life, tt.expectedLife)
		if len(tt.m.room) != tt.expectedRoomLen {
			t.Errorf("expected room length to be %d after playRoom, got %d", tt.expectedRoomLen, len(tt.m.room))
		}
		if len(tt.m.dungeon) != tt.expectedDungeonLen {
			t.Errorf("expected dungeon length to be %d after playRoom, got %d", tt.expectedDungeonLen, len(tt.m.dungeon))
		}
		if tt.m.weapon.card.Rank != tt.expectedWeaponRank {
			t.Errorf("expected weapon rank to be %d after playRoom, got %d", tt.expectedWeaponRank, tt.m.weapon.card.Rank)
		}
		assertExpectedSkippable(t, tt.m.skippable, tt.expectedSkippable)
		if tt.m.viewState != tt.expectedViewState {
			t.Errorf("expected viewState to be %s after playRoom, got %s", tt.expectedViewState, tt.m.viewState)
		}
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		m             model
		expectedScore int
	}{
		{
			m: model{
				life: 0,
				dungeon: []deck.Card{
					{Suit: deck.Spade, Rank: 5},
					{Suit: deck.Club, Rank: 7},
				},
				lastCard: deck.Card{Suit: deck.Spade, Rank: 9},
			},
			expectedScore: -12,
		},
		{
			m: model{
				life:     3,
				dungeon:  []deck.Card{},
				lastCard: deck.Card{Suit: deck.Spade, Rank: 9},
			},
			expectedScore: 3,
		},
		{
			m: model{
				life:     3,
				dungeon:  []deck.Card{},
				lastCard: deck.Card{Suit: deck.Heart, Rank: 9},
			},
			expectedScore: 12,
		},
	}

	for _, tt := range tests {
		score := tt.m.score()
		if score != tt.expectedScore {
			t.Errorf("expected score to be %d, got %d", tt.expectedScore, score)
		}
	}
}

func assertExpectedLife(t testing.TB, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expected life to be %d, got %d", expected, got)
	}
}

func assertExpectedSkippable(t testing.TB, got, expected bool) {
	t.Helper()
	if got != expected {
		t.Errorf("expected skippable to be %t, got %t", expected, got)
	}
}

func assertExpectedRoomLength(t testing.TB, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expected room length to be %d, got %d", expected, got)
	}
}

func assertExpectedDungeonLength(t testing.TB, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expected dungeon length to be %d, got %d", expected, got)
	}
}

func assertLastCard(t testing.TB, got, expected deck.Card) {
	t.Helper()
	if got.Rank != expected.Rank || got.Suit != expected.Suit {
		t.Errorf("expected last card to be %s, got %s", expected.String(), got.String())
	}
}

func testDungeon() []deck.Card {
	return []deck.Card{
		{Suit: deck.Spade, Rank: 2},
		{Suit: deck.Spade, Rank: 3},
		{Suit: deck.Spade, Rank: 4},
		{Suit: deck.Spade, Rank: 5},
		{Suit: deck.Spade, Rank: 6},
		{Suit: deck.Spade, Rank: 7},
		{Suit: deck.Spade, Rank: 8},
		{Suit: deck.Spade, Rank: 9},
	}
}

func testRoom(n int) []deck.Card {
	room := []deck.Card{
		{Suit: deck.Club, Rank: 6},
		{Suit: deck.Club, Rank: 7},
		{Suit: deck.Heart, Rank: 8},
		{Suit: deck.Diamond, Rank: 9},
	}
	return room[:n]
}
