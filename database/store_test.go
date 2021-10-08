package database

import (
	"testing"
)

func TestComparison(t *testing.T) {
	type TestStr struct {
		ID     string
		Number int32
	}

	idValue := "foo"
	strEqual := Filter{
		"ID": {
			Operator: Equal,
			Value:    idValue,
		},
	}

	elements := []TestStr{
		{ID: idValue, Number: 33},
		{ID: "nothing", Number: 42},
	}

	testCases := []struct {
		str         TestStr
		filters     []Filter
		expectError bool
		result      bool
	}{
		{elements[0], []Filter{strEqual}, false, true},
		{elements[1], []Filter{strEqual}, false, false},
		{elements[0], []Filter{{"ID": {Operator: Equal, Value: int32(8)}}}, true, false},
		{elements[1], []Filter{{"Number": {Operator: Equal, Value: elements[1].Number}}}, false, true},
		{elements[1], []Filter{{"Number": {Operator: LessThan, Value: elements[1].Number}}}, false, false},
		{elements[1], []Filter{
			{"Number": {Operator: GreaterThan, Value: elements[1].Number - 1}},
			{"Number": {Operator: LessThan, Value: elements[1].Number + 1}},
		}, false, false},
		{elements[1], []Filter{{"Number": {Operator: LessThanEq, Value: elements[1].Number}}}, false, true},
		{elements[1], []Filter{{"Number": {Operator: GreaterThanEq, Value: elements[1].Number}}}, false, true},
	}

	for _, testCase := range testCases {
		e := testCase.str
		result, err := checkFilters(e, testCase.filters)
		if testCase.expectError {
			if err == nil {
				t.Fail()
				t.Log(testCase)
				t.Logf("expected error on this testcase, but no error was triggered")
			}
		} else {
			if err != nil {
				t.Fail()
				t.Log(testCase)
				t.Logf("unexpected error: %s", err)
			}
			if result != testCase.result {
				t.Fail()
				t.Log(testCase)
				t.Logf("unexpected test result: %t", result)
			}
		}
	}
}
