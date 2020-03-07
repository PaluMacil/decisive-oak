package analysis

import (
	"math"
)

func gain(entropySet float64, attrValues ...AttributeValue) float64 {
	thisGain := entropySet
	for _, av := range attrValues {
		thisGain = thisGain - (float64(av.Occurrences) * av.Entropy)
	}

	return thisGain
}

func entropy(occurrences []int) float64 {
	var entropy float64
	var total int
	for _, occ := range occurrences {
		// if any target has zero occurrences, entropy is 0
		if occ == 0 {
			return 0
		}
		total += occ
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
