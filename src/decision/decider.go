package decision

import (
	"math"
)

type Decider struct {
	// prob of true
	baseProb float64
	// gaussian params of (positive, negative)
	gaussians map[string][2]gaussian
}

type meanVariance struct {
	count, mean, mean2 float64
}

func (mv *meanVariance) Add(count float64) {
	delta := count - mv.mean

	mv.count++
	mv.mean += delta / mv.count
	mv.mean2 += delta * (count - mv.mean)
}

func (mv *meanVariance) Train() gaussian {
	return gaussian{mean: mv.mean, variance: mv.mean2 / (mv.count - 1)}
}

type gaussian struct {
	mean, variance float64
}

func (g gaussian) pdf(count float64) float64 {
	return 1.0 / math.Sqrt(2*math.Pi*g.variance) * math.Exp(-math.Exp2(count-g.mean)/(2*g.variance))
}

func MakeDecider(features []string, items []item) Decider {
	positives := 0.0
	for _, item := range items {
		if item.label {
			positives++
		}
	}

	d := Decider{
		baseProb:  positives / float64(len(items)),
		gaussians: make(map[string][2]gaussian, len(features)),
	}

	for _, feature := range features {
		positiveDist, negativeDist := new(meanVariance), new(meanVariance)

		for _, item := range items {
			if item.label {
				positiveDist.Add(item.features[feature])
			} else {
				negativeDist.Add(item.features[feature])
			}
		}

		d.gaussians[feature] = [2]gaussian{positiveDist.Train(), negativeDist.Train()}
	}

	return d
}

func (d Decider) Label(features map[string]float64) bool {
	positive := d.baseProb
	negative := 1.0 - d.baseProb

	for feature, dists := range d.gaussians {
		positive *= dists[0].pdf(features[feature])
		negative *= dists[1].pdf(features[feature])
	}

	return positive > negative
}

func (d Decider) IsLeaf() bool {
	return false
}
