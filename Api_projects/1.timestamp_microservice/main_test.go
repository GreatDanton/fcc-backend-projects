package main

import (
	"testing"
)

// testing formatDate function from the main package
func TestFormatDate(t *testing.T) {

	tests := []struct {
		input   string
		unix    string
		natural string
	}{
		{"1495124047", "1495124047", "May 18, 2017"},
		{"947462400", "947462400", "January 10, 2000"},
		{"Jan 10, 2015", "1420848000", "January 10, 2015"},
		{"10 Feb, 2010", "1265760000", "February 10, 2010"},
		{"February 10, 2011", "1297296000", "February 10, 2011"},
		{"1999 feb 10", "918604800", "February 10, 1999"},
		{"December 30, 2018", "1546128000", "December 30, 2018"},
		{"2015 January 10", "1420848000", "January 10, 2015"},
		{"10 2015 January", "1420848000", "January 10, 2015"},
		{"January 10 2015", "1420848000", "January 10, 2015"},
		{"10 1 2010", "1285891200", "October 1, 2010"},
		{"2010 10 1", "1285891200", "October 1, 2010"},
	}

	for _, test := range tests {
		type testWants struct {
			unix    string
			natural string
		}

		wanted := testWants{unix: test.unix, natural: test.natural}
		funcOut := formatDate(test.input)

		if funcOut.Unix != wanted.unix || funcOut.Natural != wanted.natural {
			t.Errorf("formatDate(%v) = %v; wanted: %v", test.input, funcOut, wanted)
		}

	}

}
