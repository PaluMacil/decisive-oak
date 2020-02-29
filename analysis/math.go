package analysis

import (
	"fmt"
	"math"
)

type AttributeValue struct {
	Entropy  float64
	PofValue float64
}

func gain(entropySet float64, attrValues ...AttributeValue) float64 {
	thisGain := entropySet
	for _, av := range attrValues {
		thisGain = thisGain - (av.PofValue * av.Entropy)
	}
	fmt.Println(thisGain)
	return thisGain
}

func entropy(pos, neg float64, attr string) float64 {
	var entropy float64
	if pos == 0 || neg == 0 {
		entropy = 0
	} else {
		total := pos + neg
		PofPos := pos / total
		PofNeg := neg / total

		entropy = (-1 * PofPos * math.Log2(PofPos)) - (PofNeg * math.Log2(PofNeg))
	}

	fmt.Println("entropy for", attr, entropy)
	return entropy
}
