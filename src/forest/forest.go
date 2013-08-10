package forest

import (
	"decision"
	"math/rand"
)

type item struct {
	features map[string]float64
	label    bool
}

type Forest struct {
	rng *rand.Rand

	holdOutPct float64
	holdOut    [][]item

	trees []*decision.Tree
}

func MakeForest(rng *rand.Rand, trees, split, features int, holdOut float64) *Forest {
	f := &Forest{
		rng:        rng,
		trees:      make([]*decision.Tree, trees),
		holdOutPct: holdOut,
	}

	for idx := range f.trees {
		f.trees[idx] = decision.MakeTree(rng, split, features)
	}

	return f
}

func (f *Forest) Learn(features map[string]float64, label bool) {
	for idx, t := range f.trees {
		if f.rng.Float64() < f.holdOutPct {
			f.holdOut[idx] = append(f.holdOut[idx], item{features, label})
		} else {
			t.Learn(features, label)
		}
	}
}

func (f *Forest) Label(features map[string]float64) float64 {
	score := 0
	for _, t := range f.trees {
		if t.Label(features) {
			score++
		}
	}

	return float64(score) / float64(len(f.trees))
}

func (f *Forest) Evaluate() float64 {
	score := 0
	count := 0

	for idx, t := range f.trees {
		for _, example := range f.holdOut[idx] {
			count++
			if example.label == t.Label(example.features) {
				score++
			}
		}
	}

	return float64(score) / float64(count)
}
