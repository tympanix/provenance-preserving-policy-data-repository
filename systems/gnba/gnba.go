package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/systems/ba"
)

// GNBA is a structure of a generalized non-deterministic Büchi automaton
type GNBA struct {
	States         []*ba.State
	StartingStates ba.StateSet
	FinalStates    []ba.StateSet
}

// NewGNBA return a new empty GNBA
func NewGNBA() *GNBA {
	return &GNBA{
		States:         make([]*ba.State, 0),
		StartingStates: ba.NewStateSet(),
		FinalStates:    make([]ba.StateSet, 0),
	}
}

// IsAcceptanceState return true if state is in any of the acceptance sets
func (g *GNBA) IsAcceptanceState(state *ba.State) (int, bool) {
	for i, s := range g.FinalStates {
		if s.Contains(state) {
			return i, true
		}
	}
	return -1, false
}

// IsStartingState returns true if state is a starting state for the GNBA
func (g *GNBA) IsStartingState(state *ba.State) bool {
	return g.StartingStates.Contains(state)
}

// HasState returns true if state is a part of the GNBA
func (g *GNBA) HasState(state *ba.State) bool {
	for _, s := range g.States {
		if s == state {
			return true
		}
	}
	return false
}

// FindStateIndex finds the index of the state in the GNBA structure
func (g *GNBA) FindStateIndex(state *ba.State) int {
	for i, s := range g.States {
		if s == state {
			return i
		}
	}
	return -1
}

// Copy creates a copy of the GNBA
func (g *GNBA) Copy() *GNBA {
	gnba := NewGNBA()

	var rt = make(ba.RenameTable)

	// Create a copy of each state and add to rename table
	for _, s := range g.States {
		copy := s.Copy()
		rt[s] = copy
		gnba.States = append(gnba.States, copy)
	}

	// Translate state transitions with renaming table
	for _, s := range gnba.States {
		s.Rename(rt)
	}

	// Copy and rename starting states
	for s := range g.StartingStates {
		gnba.StartingStates.Add(rt[s])
	}

	// Copy and rename acceptance set
	accSet := make([]ba.StateSet, 0)
	for _, s := range g.FinalStates {
		accSet = append(accSet, s.Copy(rt))
	}
	gnba.FinalStates = accSet

	return gnba
}

func (g GNBA) String() string {
	var sb strings.Builder
	for _, s := range g.States {
		var prefix string
		if g.IsStartingState(s) {
			prefix = ">"
		}

		var suffix string
		if i, ok := g.IsAcceptanceState(s); ok {
			suffix = fmt.Sprintf("{%d}", i)
		}

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, s.ElementarySet, suffix)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t-->\t%s\n", t.Label, t.State.ElementarySet)
		}
	}

	return sb.String()
}