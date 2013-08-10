package main

import (
	"bufio"
	"flag"
	"forest"
	"math/rand"
	"os"
	"vw"
)

var seed = flag.Int64("seed", 0, "random seed")
var trees = flag.Int("trees", 0, "num trees")
var split = flag.Int("split", 0, "num examples before splitting")
var features = flag.Int("features", 3, "num features per node")
var holdOut = flag.Float64("holdout", 0.05, "hold-out pct")

func main() {
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))
	f := forest.MakeForest(rng, *trees, *split, *features, *holdOut)

	br := bufio.NewReader(os.Stdin)
	for line, err := br.ReadSlice('\n'); err == nil; line, err = br.ReadSlice('\n') {
		data, err := vw.Parse(line)
		if err != nil {
			panic(err)
		}

		f.Learn(data.CollapseFeatures(), data.Label == 1.0)
		println(f.Evaluate())
	}
}
