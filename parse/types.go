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
func (s Sample) Filter(attrName, attrValue string) (Sample, error) {
	fmt.Printf("Filtering parse sample on attribute type %s, value %s\n", attrName, attrValue)
	fmt.Printf("pre-filter attribute type names and values:\n%s", s.AttributeTypes.TerminalSummary())
	var filteredExamples Examples
	remainingTargetSet := make(map[string]bool)
	attrIndex, err := s.AttributeTypes.Index(attrName)
	if err != nil {
		return s, fmt.Errorf("filtering by %s, value %s: %w",
			attrName, attrValue, err)
	}
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
	at, err := s.AttributeTypes.Delete(attrName)
	if err != nil {
		return s, fmt.Errorf("filtering by %s, value %s: %w",
			attrName, attrValue, err)
	}
	s.AttributeTypes = at
	s.NumAttributes = s.NumAttributes - 1
	s.Examples = filteredExamples
	s.NumExamples = len(s.Examples)
	fmt.Printf("post-filter attribute type names and values:\n%s", s.AttributeTypes.TerminalSummary())

	return s, nil
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

func (t Targets) Index(target string) (int, error) {
	for i, thisTarget := range t {
		if thisTarget == target {

			return i, nil
		}
	}

	return -1, ErrIndexNotFound{For: target}
}

type ErrIndexNotFound struct {
	For string
}

func (e ErrIndexNotFound) Error() string {
	return fmt.Sprintf("could not find index for %s",
		e.For)
}

type AttributeTypes []AttributeType

func (at AttributeTypes) TerminalSummary() string {
	sb := strings.Builder{}
	for _, attrType := range at {
		attrTypeString := fmt.Sprintf("\t... %s: ", attrType.Name)
		sb.WriteString(attrTypeString)
		for i, attrValue := range attrType.Values {
			sb.WriteString(attrValue)
			// print values in a comma-separated list on one line
			if i != len(attrType.Values)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (at AttributeTypes) Index(name string) (int, error) {
	for i, attr := range at {
		if attr.Name == name {

			return i, nil
		}
	}

	return -1, ErrIndexNotFound{For: name}
}

func (at AttributeTypes) Delete(name string) (AttributeTypes, error) {
	attributeTypes := []AttributeType(at)
	i, err := at.Index(name)
	if err != nil {
		return attributeTypes,
			fmt.Errorf("deleting attribute type: %w", err)
	}
	attributeTypes = append(attributeTypes[:i], attributeTypes[i+1:]...)

	return attributeTypes, nil
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
func (at AttributeType) OccurrencesInTargets(s Sample) (AttributeOccurrenceLookup, error) {
	lookup := make(AttributeOccurrenceLookup)
	for _, value := range at.Values {
		lookup[value] = make([]int, s.NumTargets)
	}

	attrIndex, err := s.AttributeTypes.Index(at.Name)
	if err != nil {
		return lookup, fmt.Errorf("finding attribute type %s: %w",
			at.Name, err)
	}
	for _, eg := range s.Examples {
		attrValue := eg.StringValues[attrIndex]
		targetIndex, err := s.Targets.Index(eg.Target)
		if err != nil {
			return lookup, fmt.Errorf("finding target %s: %w",
				eg.Target, err)
		}
		lookup[attrValue][targetIndex] += 1
	}

	return lookup, nil
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
