package parse

type Sample struct {
	NumTargets int
	Targets Targets
	NumAttributes int
	AttributeTypes AttributeTypes
	NumExamples int
	Examples Examples
}

type Targets []string

type AttributeTypes []AttributeType

type AttributeType struct {
	Name string
	NumValues int
	Values []string
}

type Examples []Example

type Example struct {
	Values []string
	Target string
}
