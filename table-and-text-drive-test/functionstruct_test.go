package functionstruct

import (
	"testing"
)

var tests = []struct { // Test table
	arg1   int
	arg2   int
	result int
}{
	{1, 1, 2},
	{2, 2, 5},
	{3, 3, 7},
}

func TestFunction(t *testing.T) {
	for i, tt := range tests {
		s := FAdd(tt.arg1, tt.arg2)
		if s != tt.result {
			// fmt.Println("get a error")
			t.Errorf("%d. FAdd(%d,%d) => %d, wanted: %d, %s", i, tt.arg1, tt.arg2, s, tt.result, "this is error msg")
		}
	}
}
