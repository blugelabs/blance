package blance

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func TestflattenNodesByState(t *testing.T) {
	tests := []struct {
		a   map[string][]string
		exp []string
	}{
		{map[string][]string{},
			[]string{}},
		{map[string][]string{"master": []string{}},
			[]string{}},
		{map[string][]string{"master": []string{"a"}},
			[]string{"a"}},
		{map[string][]string{"master": []string{"a", "b"}},
			[]string{"a", "b"}},
		{map[string][]string{
			"master": []string{"a", "b"},
			"slave":  []string{"c"},
		}, []string{"a", "b", "c"}},
		{map[string][]string{
			"master": []string{"a", "b"},
			"slave":  []string{},
		}, []string{"a", "b"}},
	}
	for i, c := range tests {
		r := flattenNodesByState(c.a)
		if !reflect.DeepEqual(r, c.exp) {
			t.Errorf("i: %d, a: %#v, exp: %#v, got: %#v",
				i, c.a, c.exp, r)
		}
	}
}

func TestRemoveNodesFromNodesByState(t *testing.T) {
	tests := []struct {
		nodesByState map[string][]string
		removeNodes  []string
		exp          map[string][]string
	}{
		{map[string][]string{"master": []string{"a", "b"}},
			[]string{"a", "b"},
			map[string][]string{"master": []string{}},
		},
		{map[string][]string{"master": []string{"a", "b"}},
			[]string{"b", "c"},
			map[string][]string{"master": []string{"a"}},
		},
		{map[string][]string{"master": []string{"a", "b"}},
			[]string{"a", "c"},
			map[string][]string{"master": []string{"b"}},
		},
		{map[string][]string{"master": []string{"a", "b"}},
			[]string{},
			map[string][]string{"master": []string{"a", "b"}},
		},
		{
			map[string][]string{
				"master": []string{"a", "b"},
				"slave":  []string{"c"},
			},
			[]string{},
			map[string][]string{
				"master": []string{"a", "b"},
				"slave":  []string{"c"},
			},
		},
		{
			map[string][]string{
				"master": []string{"a", "b"},
				"slave":  []string{"c"},
			},
			[]string{"a"},
			map[string][]string{
				"master": []string{"b"},
				"slave":  []string{"c"},
			},
		},
		{
			map[string][]string{
				"master": []string{"a", "b"},
				"slave":  []string{"c"},
			},
			[]string{"a", "c"},
			map[string][]string{
				"master": []string{"b"},
				"slave":  []string{},
			},
		},
	}
	for i, c := range tests {
		r := removeNodesFromNodesByState(c.nodesByState, c.removeNodes, nil)
		if !reflect.DeepEqual(r, c.exp) {
			t.Errorf("i: %d, nodesByState: %#v, removeNodes: %#v, exp: %#v, got: %#v",
				i, c.nodesByState, c.removeNodes, c.exp, r)
		}
	}
}

func TestStateNameSorter(t *testing.T) {
	tests := []struct {
		m   PartitionModel
		s   []string
		exp []string
	}{
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{},
			[]string{},
		},
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{"master", "slave"},
			[]string{"master", "slave"},
		},
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{"slave", "master"},
			[]string{"master", "slave"},
		},
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{"a", "master"},
			[]string{"a", "master"},
		},
		{
			PartitionModel{
				"master": &PartitionModelState{Priority: 0},
				"slave":  &PartitionModelState{Priority: 1},
			},
			[]string{"master", "a"},
			[]string{"a", "master"},
		},
	}
	for i, c := range tests {
		sort.Sort(&stateNameSorter{m: c.m, s: c.s})
		if !reflect.DeepEqual(c.s, c.exp) {
			t.Errorf("i: %d, m: %#v, s: %#v, exp: %#v",
				i, c.m, c.s, c.exp)
		}
	}
}

func TestCountStateNodes(t *testing.T) {
	tests := []struct {
		m   PartitionMap
		w   map[string]int
		exp map[string]map[string]int
	}{
		{
			PartitionMap{
				"0": &Partition{NodesByState: map[string][]string{
					"master": []string{"a"},
					"slave":  []string{"b", "c"},
				}},
				"1": &Partition{NodesByState: map[string][]string{
					"master": []string{"b"},
					"slave":  []string{"c"},
				}},
			},
			nil,
			map[string]map[string]int{
				"master": map[string]int{
					"a": 1,
					"b": 1,
				},
				"slave": map[string]int{
					"b": 1,
					"c": 2,
				},
			},
		},
		{
			PartitionMap{
				"0": &Partition{NodesByState: map[string][]string{
					"slave": []string{"b", "c"},
				}},
				"1": &Partition{NodesByState: map[string][]string{
					"master": []string{"b"},
					"slave":  []string{"c"},
				}},
			},
			nil,
			map[string]map[string]int{
				"master": map[string]int{
					"b": 1,
				},
				"slave": map[string]int{
					"b": 1,
					"c": 2,
				},
			},
		},
	}
	for i, c := range tests {
		r := countStateNodes(c.m, c.w)
		if !reflect.DeepEqual(r, c.exp) {
			t.Errorf("i: %d, m: %#v, w: %#v, exp: %#v",
				i, c.m, c.w, c.exp)
		}
	}
}

func TestPartitionMapToArrayCopy(t *testing.T) {
	tests := []struct {
		m   PartitionMap
		exp []*Partition
	}{
		{
			PartitionMap{
				"0": &Partition{
					Name: "0",
					NodesByState: map[string][]string{
						"master": []string{"a"},
						"slave":  []string{"b", "c"},
					},
				},
				"1": &Partition{
					Name: "1",
					NodesByState: map[string][]string{
						"master": []string{"b"},
						"slave":  []string{"c"},
					},
				},
			},
			[]*Partition{
				&Partition{
					Name: "0",
					NodesByState: map[string][]string{
						"master": []string{"a"},
						"slave":  []string{"b", "c"},
					},
				},
				&Partition{
					Name: "1",
					NodesByState: map[string][]string{
						"master": []string{"b"},
						"slave":  []string{"c"},
					},
				},
			},
		},
	}
	for _, c := range tests {
		r := c.m.toArrayCopy()
		testSubset := func(a, b []*Partition) {
			if len(a) != len(b) {
				t.Errorf("expected same lengths")
			}
			for _, ap := range a {
				found := false
				for _, bp := range b {
					if reflect.DeepEqual(ap, bp) {
						found = true
					}
				}
				if !found {
					t.Errorf("couldn't find a entry in b")
				}
			}
		}
		testSubset(r, c.exp)
		testSubset(c.exp, r)
	}
}

func TestPlanNextMap(t *testing.T) {
	tests := []struct {
		PrevMap               PartitionMap
		Nodes                 []string
		NodesToRemove         []string
		NodesToAdd            []string
		Model                 PartitionModel
		ModelStateConstraints map[string]int
		PartitionWeights      map[string]int
		StateStickiness       map[string]int
		NodeWeights           map[string]int
		exp                   PartitionMap
		expNumWarnings        int
	}{
		{
			PrevMap: PartitionMap{
				"0": &Partition{
					Name:         "0",
					NodesByState: map[string][]string{},
				},
				"1": &Partition{
					Name:         "1",
					NodesByState: map[string][]string{},
				},
			},
			Nodes:         []string{"a"},
			NodesToRemove: []string{},
			NodesToAdd:    []string{"a"},
			Model: PartitionModel{
				"master": &PartitionModelState{
					Priority: 0, Constraints: 0,
				},
				"slave": &PartitionModelState{
					Priority: 1, Constraints: 0,
				},
			},
			ModelStateConstraints: nil,
			PartitionWeights:      nil,
			StateStickiness:       nil,
			NodeWeights:           nil,
			exp: PartitionMap{
				"0": &Partition{
					Name: "0",
					NodesByState: map[string][]string{
						"master": []string{"a"},
					},
				},
				"1": &Partition{
					Name: "1",
					NodesByState: map[string][]string{
						"master": []string{"a"},
					},
				},
			},
			expNumWarnings: 0,
		},
	}
	for i, c := range tests {
		rMap, rWarnings := planNextMap(
			c.PrevMap,
			c.Nodes,
			c.NodesToRemove,
			c.NodesToAdd,
			c.Model,
			c.ModelStateConstraints,
			c.PartitionWeights,
			c.StateStickiness,
			c.NodeWeights)
		if !reflect.DeepEqual(rMap, c.exp) {
			jc, _ := json.Marshal(c)
			jrMap, _ := json.Marshal(rMap)
			jexp, _ := json.Marshal(c.exp)
			t.Errorf("i: %d, planNextMap, c: %s, rMap: %s, exp: %s",
				i, jc, jrMap, jexp)
		}
		if c.expNumWarnings != len(rWarnings) {
			t.Errorf("i: %d, planNextMap.warnings, c: %#v, rWarnings: %d, expNumWarnings: %d",
				i, c, rWarnings, c.expNumWarnings)
		}
	}
}
