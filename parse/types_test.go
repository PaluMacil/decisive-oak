package parse_test

import (
	"github.com/PaluMacil/decisive-oak/parse"
	"testing"
)

func Test_AttributeTypes_Index(t *testing.T) {
	i := attributeTypes.Index("astigmatism")
	const expected = 2
	if i != expected {
		t.Errorf("for 'astigmatism' expected index %d, got %d", expected, i)
	}
}

func Test_AttributeTypes_Delete(t *testing.T) {
	at := make(parse.AttributeTypes, len(attributeTypes))
	copy(at, attributeTypes)
	originalLength := len(at)
	at = at.Delete("age")
	if len(at) != originalLength-1 {
		t.Errorf("after deleting an element, expected length %d, got %d", originalLength-1, len(at))
	}
}

func Test_AttributeTypes_IsValid(t *testing.T) {
	at := attributeTypes

	r1 := at.IsValid("tear-rate")
	if !r1 {
		t.Errorf("expected true, got false for 'tear-rate'")
	}
	r2 := at.IsValid("frogs")
	if r2 {
		t.Errorf("expected false, got true for 'frogs'")
	}
}

func Test_AttributeType_IsValidValue(t *testing.T) {
	at := attributeTypes[0]

	r1 := at.IsValidValue("old")
	if r1 {
		t.Errorf("expected false, got true for 'old'")
	}
	r2 := at.IsValidValue("young")
	if !r2 {
		t.Errorf("expected true, got false for 'young'")
	}
}

func Test_AttributeOccurrenceLookup_AttributeValueTotal(t *testing.T) {
	sample, err := parse.FromFile("../new-treatment.data.txt")
	if err != nil {
		t.Errorf("failed parsing file new-treatment.data.txt: %v", err)
	}
	if checkTotal(sample, "age", "<25") != 1 {
		t.Errorf("total for age >25 should be 1")
	}
	if checkTotal(sample, "pulse", "normal") != 3 {
		t.Errorf("total for pulse normal should be 3")
	}
	if checkTotal(sample, "pulse", "rapid") != 2 {
		t.Errorf("total for pulse normal should be 2")
	}
	if checkTotal(sample, "bp", "normal") != 3 {
		t.Errorf("total for bp normal should be 3")
	}
}

func checkTotal(sample parse.Sample, attrTypeName, attrValue string) int {
	return sample.AttributeTypes[sample.AttributeTypes.
		Index(attrTypeName)].
		OccurrencesInTargets(sample).
		AttributeValueTotal(attrValue)
}

func Test_Example_DeleteValue(t *testing.T) {
	example := parse.Example{
		StringValues: []string{
			"young",
			"hypermetrope",
			"no",
			"normal",
		},
	}
	originalLength := len(example.StringValues)
	example = example.DeleteValue(2)
	if len(example.StringValues) != originalLength-1 {
		t.Errorf("after deleting an element, expected length %d, got %d",
			originalLength-1, len(example.StringValues))
	}
}

func Test_Sample_Filter(t *testing.T) {
	sample, err := parse.FromFile("../contact-lenses.data.txt")
	if err != nil {
		t.Errorf("failed parsing file contact-lenses.data.txt: %v", err)
	}
	ageSample := sample.Filter("age", "young")
	if ageSample.NumTargets != 3 {
		t.Errorf("filtering on age young: expected NumTargets %d, got %d",
			3, ageSample.NumTargets)
	}
	if len(ageSample.Targets) != 3 {
		t.Errorf("filtering on age young: expected length of Targets %d, got %d",
			3, len(ageSample.Targets))
	}
	if ageSample.NumAttributes != sample.NumAttributes-1 {
		t.Errorf("filtering on age young: expected NumAttributes to decrease by one")
	}
	if len(ageSample.AttributeTypes) != 3 {
		t.Errorf("filtering on age young: expected AttributeTypes %d, got %d",
			3, len(ageSample.AttributeTypes))
	}
	if ageSample.NumExamples != len(ageSample.Examples) {
		t.Errorf("filtering on age young: expected NumExamples to equal Examples length")
	}
	if len(ageSample.Examples) != 8 {
		t.Errorf("filtering on age young: expected 8 examples, got %d", len(ageSample.Examples))
	}

	astigmatismSample := ageSample.Filter("astigmatism", "no")
	if astigmatismSample.NumTargets != 2 {
		t.Errorf("filtering on astigmatism no: expected NumTargets %d, got %d",
			2, astigmatismSample.NumTargets)
	}
	if len(astigmatismSample.Targets) != 2 {
		t.Errorf("filtering on astigmatism no: expected length of Targets %d, got %d",
			2, len(astigmatismSample.Targets))
	}
	if astigmatismSample.NumAttributes != ageSample.NumAttributes-1 {
		t.Errorf("filtering on astigmatism no: expected NumAttributes to decrease by one")
	}
	if len(astigmatismSample.AttributeTypes) != 2 {
		t.Errorf("filtering on astigmatism no: expected AttributeTypes %d, got %d",
			2, len(astigmatismSample.AttributeTypes))
	}
	if astigmatismSample.NumExamples != len(astigmatismSample.Examples) {
		t.Errorf("filtering on astigmatism no: expected NumExamples to equal Examples length")
	}
	if len(astigmatismSample.Examples) != 4 {
		t.Errorf("filtering on astigmatism no: expected 4 examples, got %d", len(ageSample.Examples))
	}
}

var attributeTypes = parse.AttributeTypes{
	parse.AttributeType{
		Name:      "age",
		NumValues: 3,
		Values: []string{
			"young",
			"pre-presbyopic",
			"presbyopic",
		},
		Real: false,
	},
	parse.AttributeType{
		Name:      "prescription",
		NumValues: 2,
		Values: []string{
			"myope",
			"hypermetrope",
		},
		Real: false,
	},
	parse.AttributeType{
		Name:      "astigmatism",
		NumValues: 2,
		Values: []string{
			"no",
			"yes",
		},
		Real: false,
	},
	parse.AttributeType{
		Name:      "tear-rate",
		NumValues: 2,
		Values: []string{
			"reduced",
			"normal",
		},
		Real: false,
	},
}
