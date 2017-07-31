package parse

import (
	"reflect"
	"testing"
)

// testing parseImageContainer function
func Test_ParseImageContainers(t *testing.T) {
	tests := []struct {
		input string
		wants []string
	}{
		{
			`
			<td><p> Some text 1 </p></td>
			<td style="width:25%"><p> Some text </p></td>
			<td style="width:50%"><a href="http://www.google.com">Google</a></td>
			<td style="width:25%"><img src="http://www.somelink.com"/></td>
			`,
			[]string{
				`<td style="width:25%"><p> Some text </p></td>`,
				`<td style="width:25%"><img src="http://www.somelink.com"/></td>`,
			},
		},
	}

	for i, test := range tests {
		out := parseImageContainers(test.input, []string{})
		if !reflect.DeepEqual(out, test.wants) {
			t.Errorf("Test Case(%d)", i)
			t.Errorf("parseImageContainers(input) = %s", out)
			t.Errorf("Wants: %v", test.wants)
		}
	}
}
