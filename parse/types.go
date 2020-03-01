package parse

import (
	"fmt"
	"strings"
)

type Sample struct {
	NumTargets int
	Targets Targets
	NumAttributes int
	AttributeTypes AttributeTypes
	NumExamples int
	Examples Examples
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

type AttributeTypes []AttributeType

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
	Name string
	NumValues int
	Values []string
	Real bool
}

func (at AttributeType) IsValidValue(value string) bool {
	for _, v := range at.Values {
		if v == value {
			return true
		}
	}

	return false
}

type Examples []Example

type Example struct {
	StringValues []string
	RealValues []float64
	Target string
}
