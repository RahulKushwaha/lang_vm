package vm

import (
	"github.com/stretchr/testify/assert"
	"lang_vm/code"
	"lang_vm/object"
	"testing"
)

func TestVmAdd(t *testing.T) {
	type testCase struct {
		ins code.Instructions
		out object.Object
	}

	testCases := map[string]testCase{
		"add_two_nums": {
			ins: code.NewBuilder().
				Add(code.OpConstant, 1).
				Add(code.OpConstant, 3).
				Add(code.OpAdd).
				Add(code.OpHalt).
				Build(),
			out: &object.Integer{Value: 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm := New(tc.ins)

			err := vm.Run()
			assert.NoError(t, err, name)

			pop := vm.Pop()
			assert.Equal(t, pop, tc.out)
		})
	}
}
