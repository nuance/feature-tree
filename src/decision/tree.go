package decision

import (
	"math/rand"
)

type node interface {
	Label(map[string]float64) bool
	IsLeaf() bool
}

type item struct {
	features map[string]float64
	label    bool
}

type leafNode struct {
	label  bool
	buffer []item
}

func (l leafNode) Label(features map[string]float64) bool {
	return l.label
}

func (l *leafNode) Learn(features map[string]float64, label bool) {
	l.buffer = append(l.buffer, item{features, label})
}

func (l leafNode) Split(features []string) node {
	return MakeDecider(features, l.buffer)
}

func (l leafNode) Size() int {
	return len(l.buffer)
}

func (l leafNode) IsLeaf() bool {
	return true
}

type Tree struct {
	rng            *rand.Rand
	splitThreshold int
	numFeatures    int

	nodes       []node
	features    map[string]bool
	allFeatures []string
}

func MakeTree(rng *rand.Rand, splitThreshold int, numFeatures int) *Tree {
	return &Tree{
		rng:            rng,
		splitThreshold: splitThreshold,
		numFeatures:    numFeatures,
		nodes:          []node{leafNode{label: true}},
		features:       map[string]bool{},
	}
}

func (t Tree) find(features map[string]float64) int {
	idx := 0

	for !t.nodes[idx].IsLeaf() {
		if t.nodes[idx].Label(features) {
			// go right
			idx = (idx + 1) * 2
		} else {
			// go left
			idx = (idx+1)*2 - 1
		}
	}

	return idx
}

func (t *Tree) Label(features map[string]float64) bool {
	return t.nodes[t.find(features)].Label(features)
}

func (t *Tree) Learn(features map[string]float64, label bool) {
	for key := range features {
		if !t.features[key] {
			t.allFeatures = append(t.allFeatures, key)
			t.features[key] = true
		}
	}

	idx := t.find(features)
	leaf := t.nodes[idx].(*leafNode)

	leaf.Learn(features, label)
	if leaf.Size() > t.splitThreshold {
		features := make([]string, t.numFeatures)
		for i := 0; i < t.numFeatures; i++ {
			features[i] = t.allFeatures[t.rng.Int()%len(t.allFeatures)]
		}

		t.nodes[idx] = leaf.Split(features)
		t.nodes = append(t.nodes, &leafNode{label: false}, &leafNode{label: true})
	}
}
