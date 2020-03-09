package parse

import (
	"fmt"
	"strings"
)

type Sample struct {
	NumTargets     int
	Targets        Targets
	NumAttributes  int
	AttributeTypes AttributeTypes
	NumExamples    int
	Examples       Examples
}

// Filter takes the given attribute name and value and filters the Sample to reflect only this,
// returning a reduced Sample copy
func (s Sample) Filter(attrName, attrValue string) Sample {
	var filteredExamples Examples
	remainingTargetSet := make(map[string]bool)
	attrIndex := s.AttributeTypes.Index(attrName)
	for _, eg := range s.Examples {
		match := eg.StringValues[attrIndex] == attrValue
		eg = eg.DeleteValue(attrIndex)
		if match {
			filteredExamples = append(filteredExamples, eg)
			remainingTargetSet[eg.Target] = true
		}
	}
	// reset targets list
	s.Targets = make(Targets, 0)
	for target := range remainingTargetSet {
		s.Targets = append(s.Targets, target)
	}
	s.NumTargets = len(s.Targets)
	s.AttributeTypes = s.AttributeTypes.Delete(attrName)
	s.NumAttributes = s.NumAttributes - 1
	s.Examples = filteredExamples
	s.NumExamples = len(s.Examples)

	return s
}

type Targets []string

func (t Targets) IsValid(target string) bool {
	for _, validTarget := range t {
		if validTarget == target {
			return true
		}
	}
	return false
}

func (t Targets) Index(target string) int {
	for i, thisTarget := range t {
		if thisTarget == target {

			return i
		}
	}

	return -1
}

type AttributeTypes []AttributeType

func (at AttributeTypes) Index(name string) int {
	for i, attr := range at {
		if attr.Name == name {

			return i
		}
	}

	return -1
}

func (at AttributeTypes) Delete(name string) AttributeTypes {
	attributeTypes := []AttributeType(at)
	i := at.Index(name)
	if i != -1 {
		attributeTypes = append(attributeTypes[:i], attributeTypes[i+1:]...)
	}

	return attributeTypes
}

func (at AttributeTypes) String() string {
	var names strings.Builder
	for _, t := range at {
		names.WriteString(t.Name)
		names.WriteString(", ")
	}

	return fmt.Sprintf("[%s]", names.String())
}

func (at AttributeTypes) IsValid(attributeTypeName string) bool {
	for _, t := range at {
		if t.Name == attributeTypeName {
			return true
		}
	}

	return false
}

type AttributeType struct {
	Name      string
	NumValues int
	Values    []string
	Real      bool
}

func (at AttributeType) IsValidValue(value string) bool {
	for _, v := range at.Values {
		if v == value {
			return true
		}
	}

	return false
}

// OccurrencesInTargets returns an AttributeOccurrenceLookup with method
// AttributeValueTotal(attrValue string) int
func (at AttributeType) OccurrencesInTargets(s Sample) AttributeOccurrenceLookup {
	lookup := make(AttributeOccurrenceLookup)
	for _, value := range at.Values {
		lookup[value] = make([]int, s.NumTargets)
	}

	attrIndex := s.AttributeTypes.Index(at.Name)
	for _, eg := range s.Examples {
		attrValue := eg.StringValues[attrIndex]
		targetIndex := s.Targets.Index(eg.Target)
		lookup[attrValue][targetIndex] += 1
	}

	return lookup
}

// AttributeOccurrenceLookup is a map of attribute value to slice of int.
// slice of target occurrences
//  {
//    '<25': [1, 0],
//    '25-40': [1, 0],
//    '>40': [1, 2]
//  }
type AttributeOccurrenceLookup map[string][]int

func (aom AttributeOccurrenceLookup) AttributeValueTotal(attrValue string) int {
	var total int
	targetOccurrences := aom[attrValue]

	for _, value := range targetOccurrences {
		total += value
	}

	return total
}

type Examples []Example

type Example struct {
	StringValues []string
	RealValues   []float64
	Target       string
}

func (eg Example) DeleteValue(index int) Example {
	if len(eg.StringValues) > 0 {
		eg.StringValues = append(eg.StringValues[:index], eg.StringValues[index+1:]...)
	} else {
		eg.RealValues = append(eg.RealValues[:index], eg.RealValues[index+1:]...)
	}

	return eg
}
