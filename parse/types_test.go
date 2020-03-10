package parse_test

import (
	"github.com/PaluMacil/decisive-oak/parse"
	"testing"
)

func Test_AttributeTypes_Index(t *testing.T) {
	i, _ := attributeTypes.Index("astigmatism")
	const expected = 2
	if i != expected {
		t.Errorf("for 'astigmatism' expected index %d, got %d", expected, i)
	}
}

func Test_AttributeTypes_Delete(t *testing.T) {
	at := make(parse.AttributeTypes, len(attributeTypes))
	copy(at, attributeTypes)
	originalLength := len(at)
	at, err := at.Delete("age")
	if err != nil {
		t.Error(err.Error())
	}
	if at == nil {
		t.Errorf("attribute types were nil after delete")
	}
	if len(at) != originalLength-1 {
		t.Errorf("after deleting an element, expected length %d, got %d", originalLength-1, len(at))
	}
	if _, err = at.Delete("fakeNotExist"); err == nil {
		t.Error("expected error when deleting attribute type that doesn't exist")
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
	sample, err := parse.FromFile("../data/new-treatment.data.txt")
	if err != nil {
		t.Errorf("failed parsing file new-treatment.data.txt: %v", err)
	}
	if checkTotal(sample, "age", "<25", t) != 1 {
		t.Errorf("total for age >25 should be 1")
	}
	if checkTotal(sample, "pulse", "normal", t) != 3 {
		t.Errorf("total for pulse normal should be 3")
	}
	if checkTotal(sample, "pulse", "rapid", t) != 2 {
		t.Errorf("total for pulse normal should be 2")
	}
	if checkTotal(sample, "bp", "normal", t) != 3 {
		t.Errorf("total for bp normal should be 3")
	}
}

func checkTotal(sample parse.Sample, attrTypeName, attrValue string, t *testing.T) int {
	idx, err := sample.AttributeTypes.Index(attrTypeName)
	if err != nil {
		t.Errorf("checking total: %s", err.Error())
	}
	lookup, err := sample.AttributeTypes[idx].OccurrencesInTargets(sample)
	if err != nil {
		t.Errorf("getting occurrences in targets: %s", err.Error())
	}
	if lookup == nil {
		t.Error("lookup was nil")
	}

	return lookup.AttributeValueTotal(attrValue)
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
	sample, err := parse.FromFile("../data/contact-lenses.data.txt")
	if err != nil {
		t.Errorf("failed parsing file contact-lenses.data.txt: %v", err)
	}
	ageSample, err := sample.Filter("age", "young")
	if err != nil {
		t.Errorf("filtering on age, young: %s", err.Error())
	}
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

	astigmatismNoSample, err := ageSample.Filter("astigmatism", "no")
	if err != nil {
		t.Errorf("filtering on astigmatism, no: %s", err.Error())
	}
	if astigmatismNoSample.NumTargets != 2 {
		t.Errorf("filtering on astigmatism no: expected NumTargets %d, got %d",
			2, astigmatismNoSample.NumTargets)
	}
	if len(astigmatismNoSample.Targets) != 2 {
		t.Errorf("filtering on astigmatism no: expected length of Targets %d, got %d",
			2, len(astigmatismNoSample.Targets))
	}
	if astigmatismNoSample.NumAttributes != ageSample.NumAttributes-1 {
		t.Errorf("filtering on astigmatism no: expected NumAttributes to decrease by one")
	}
	if len(astigmatismNoSample.AttributeTypes) != 2 {
		t.Errorf("filtering on astigmatism no: expected AttributeTypes %d, got %d",
			2, len(astigmatismNoSample.AttributeTypes))
	}
	if astigmatismNoSample.NumExamples != len(astigmatismNoSample.Examples) {
		t.Errorf("filtering on astigmatism no: expected NumExamples to equal Examples length")
	}
	if len(astigmatismNoSample.Examples) != 4 {
		t.Errorf("filtering on astigmatism no: expected 4 examples, got %d", len(ageSample.Examples))
	}

	astigmatismYesSample, err := ageSample.Filter("astigmatism", "yes")
	if err != nil {
		t.Errorf("filtering on astigmatism, yes: %s", err.Error())
	}
	if len(astigmatismYesSample.Examples) != 4 {
		t.Errorf("filtering on astigmatism yes: expected 4 examples, got %d", len(astigmatismYesSample.Examples))
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
