package blance

import (
	"reflect"
	"strings"
	"testing"
)

func TestFindStateChanges(t *testing.T) {
	tests := []struct {
		begStateIdx     int
		endStateIdx     int
		state           string
		states          []string
		begNodesByState map[string][]string
		endNodesByState map[string][]string
		expected        []string
	}{
		{0, 0, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c"},
			},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c"},
			},
			nil,
		},
		{1, 2, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c"},
			},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c"},
			},
			nil,
		},
		{0, 0, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{},
				"replica": []string{"a"},
			},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{},
			},
			nil,
		},
		{1, 2, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{},
				"replica": []string{"a"},
			},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{},
			},
			[]string{"a"},
		},
		{0, 1, "replica",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{},
			},
			map[string][]string{
				"master":  []string{},
				"replica": []string{"a"},
			},
			[]string{"a"},
		},
		{1, 2, "replica",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{},
			},
			map[string][]string{
				"master":  []string{},
				"replica": []string{"a"},
			},
			nil,
		},
		{1, 2, "replica",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{},
				"replica": []string{"a"},
			},
			map[string][]string{
				"master":  []string{},
				"replica": []string{},
			},
			nil,
		},
		{1, 2, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c", "d"},
			},
			map[string][]string{
				"master":  []string{"b"},
				"replica": []string{"a", "c", "d"},
			},
			[]string{"b"},
		},
		{1, 2, "master",
			[]string{"master", "replica"},
			map[string][]string{
				"master":  []string{"a"},
				"replica": []string{"b", "c", "d"},
			},
			map[string][]string{
				"master":  []string{"x"},
				"replica": []string{"a", "c", "d"},
			},
			nil,
		},
	}

	for i, test := range tests {
		got := findStateChanges(test.begStateIdx, test.endStateIdx,
			test.state, test.states,
			test.begNodesByState,
			test.endNodesByState)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("i: %d, got: %#v, expected: %#v, test: %#v",
				i, got, test.expected, test)
		}
	}
}

func TestCalcPartitionMoves(t *testing.T) {
	states := []string{"master", "replica"}

	tests := []struct {
		before string
		moves  string
		after  string
	}{
		//  master | replica
		//  -------|--------
		{
			" a",
			"",
			" a",
		},
		{
			"      | a",
			"",
			"      | a",
		},
		{
			" a    | b",
			"",
			" a    | b",
		},
		{
			" a",
			` a +b |
             -a  b |`,
			"    b",
		},
		{
			" a    | b  c",
			` a +b |-b  c
             -a  b |    c
                 b |    c +d`,
			"    b |    c  d",
		},
		{
			" a    |    b",
			` a +b |   -b
             -a  b |+a`,
			"    b | a",
		},
		{
			" a    |    b",
			` a +c |    b
             -a  c |+a  b
                 c | a -b`,
			"    c | a",
		},
		{
			" a    | b",
			` a +c | b
             -a  c | b
                 c | b +d
                 c |-b  d`,
			"    c |    d",
		},
	}

	negate := map[string]string{
		"+": "-",
		"-": "+",
	}

	for testi, test := range tests {
		before := convertLineToNodesByState(test.before, states)
		after := convertLineToNodesByState(test.after, states)

		var movesExp []map[string][]string

		if test.moves != "" {
			moveLines := strings.Split(test.moves, "\n")
			for _, moveLine := range moveLines {
				moveExp := convertLineToNodesByState(moveLine, states)
				movesExp = append(movesExp, moveExp)
			}
		}

		movesGot := CalcPartitionMoves(states, before, after)

		if len(movesGot) != len(movesExp) {
			t.Errorf("testi: %d, mismatch lengths,"+
				" before: %#v, after: %#v,"+
				" movesExp: %#v, movesGot: %#v, test: %#v",
				testi, before, after, movesExp, movesGot, test)

			continue
		}

		for moveExpi, moveExp := range movesExp {
			moveGot := movesGot[moveExpi]

			found := false

			for statei, state := range states {
				if found {
					continue
				}

				for _, move := range moveExp[state] {
					if found {
						continue
					}

					op := move[0:1]
					if op == "+" || op == "-" {
						found = true

						if moveGot.Node != move[1:] {
							t.Errorf("testi: %d, wrong node,"+
								" before: %#v, after: %#v,"+
								" movesExp: %#v, movesGot: %#v,"+
								" test: %#v",
								testi, before, after,
								movesExp, movesGot, test)
						}

						flipSideFound := ""
						flipSideState := ""
						flipSide := negate[op] + move[1:]
						for j := statei + 1; j < len(states); j++ {
							for _, x := range moveExp[states[j]] {
								if x == flipSide {
									flipSideFound = flipSide
									flipSideState = states[j]
								}
							}
						}

						stateExp := state
						if flipSideFound != "" {
							if op == "-" {
								stateExp = flipSideState
							}
						} else {
							if op == "-" {
								stateExp = ""
							}
						}

						if moveGot.State != stateExp {
							t.Errorf("testi: %d, not stateExp: %q,"+
								" before: %#v, after: %#v,"+
								" movesExp: %#v, movesGot: %#v,"+
								" test: %#v, move: %s,"+
								" flipSideFound: %q, flipSideState: %q",
								testi, stateExp, before, after,
								movesExp, movesGot, test, move,
								flipSideFound, flipSideState)
						}
					}
				}
			}
		}
	}
}

// Converts an input line string like " a b | +c -d", with input
// states of ["master", "replica"] to something like {"master": ["a",
// "b"], "replica": ["+c", "-d"]}.
func convertLineToNodesByState(
	line string, states []string) map[string][]string {
	nodesByState := map[string][]string{}

	line = strings.Trim(line, " ")
	for {
		linex := strings.Replace(line, "  ", " ", -1)
		if linex == line {
			break
		}
		line = linex
	}

	parts := strings.Split(line, "|")
	for i, state := range states {
		if i >= len(parts) {
			break
		}
		part := strings.Trim(parts[i], " ")
		if part != "" {
			nodes := strings.Split(part, " ")
			nodesByState[state] = append(nodesByState[state], nodes...)
		}
	}

	return nodesByState
}
