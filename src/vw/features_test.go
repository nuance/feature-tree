package vw

import (
	"testing"
)

type parseTest struct {
	example string
	data    Data
}

var tests []parseTest = []parseTest{
	{"1 1.0 |MetricFeatures:3.28 height:1.5 length:2.0 |Says black with white stripes |OtherFeatures NumberOfLegs:4.0 HasStripes",
		Data{Label: 1, Importance: 1, Tag: "",
			Features: map[string]map[string]float64{
				"MetricFeatures": {"height": 4.92, "length": 6.56},
				"Says":           {"black": 1, "with": 1, "white": 1, "stripes": 1},
				"OtherFeatures":  {"NumberOfLegs": 4, "HasStripes": 1}}}},
	{"0 1 a| a:2", Data{Label: 0, Importance: 1, Tag: "a",
		Features: map[string]map[string]float64{
			"": {"a": 2}}}},
	{"1 2 b | a:3", Data{Label: 1, Importance: 2, Tag: "b",
		Features: map[string]map[string]float64{
			"": {"a": 3}}}},
	{"0 | price:.23 sqft:.25 age:.05 2006", Data{Label: 0, Importance: 1, Tag: "",
		Features: map[string]map[string]float64{
			"": {"price": 0.23, "sqft": 0.25, "age": 0.05, "2006": 1}}}},
	{"1 2 'second_house | price:.18 sqft:.15 age:.35 1976", Data{Label: 1, Importance: 2, Tag: "'second_house",
		Features: map[string]map[string]float64{
			"": {"price": .18, "sqft": .15, "age": .35, "1976": 1}}}},
	{"0 0.5 'third_house | price:.53 sqft:.32 age:.87 1924", Data{Label: 0, Importance: 0.5, Tag: "'third_house",
		Features: map[string]map[string]float64{
			"": {"price": .53, "sqft": .32, "age": .87, "1924": 1}}}},
	{"1 | foo:2 bar | baz:1 foo:2", Data{Label: 1, Importance: 1, Tag: "",
		Features: map[string]map[string]float64{
			"": {"foo": 4, "bar": 1, "baz": 1}}}},
}

func compare(t *testing.T, expected, actual Data) {
	if expected.Label != actual.Label {
		t.Errorf("Expected label '%f', got '%f'", expected.Label, actual.Label)
	}

	if expected.Importance != actual.Importance {
		t.Errorf("Expected importance '%f', got '%f'", expected.Importance, actual.Importance)
	}

	if expected.Tag != actual.Tag {
		t.Errorf("Expected tag '%s', got '%s'", expected.Tag, actual.Tag)
	}

	for namespace := range actual.Features {
		if _, ok := expected.Features[namespace]; !ok {
			t.Errorf("Unexpected namespace '%s'", namespace)
		}
	}

	for namespace, features := range expected.Features {
		if _, ok := actual.Features[namespace]; !ok {
			t.Errorf("Missing expected namespace '%s'", namespace)
			continue
		}

		for feature, count := range features {
			if _, ok := actual.Features[namespace][feature]; !ok {
				t.Errorf("Missing expected feature '%s.%s'", namespace, feature)
			}

			if count != actual.Features[namespace][feature] {
				t.Errorf("Wrong value for feature '%s.%s', %f vs expected %f", namespace, feature, actual.Features[namespace][feature], count)
			}
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		t.Logf("Testing '%s'", test.example)

		data, err := Parse([]byte(test.example))
		if err != nil {
			t.Error(err)
		}

		compare(t, test.data, data)
	}
}
