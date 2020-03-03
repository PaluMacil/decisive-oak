package analysis

import (
	"math"
)

func gain(entropySet float64, attrValues ...Attribute) float64 {
	thisGain := entropySet
	for _, av := range attrValues {
		thisGain = thisGain - (av.PofAttribute * av.Entropy)
	}

	return thisGain
}

func entropy(occurrences []int) float64 {
	var entropy float64
	var total int
	for _, o := range occurrences {
		// if any target has zero occurrences, entropy is 0
		if o == 0 {
			return 0
		}
		total = total + o
	}
	occurrenceRatios := make([]float64, len(occurrences))
	for i := range occurrences {
		occurrenceRatios[i] = float64(occurrences[i]) / float64(total)
	}

	for _, pOfTarget := range occurrenceRatios {
		entropy = entropy + (-1 * pOfTarget * math.Log2(pOfTarget))
	}

	return entropy
}
