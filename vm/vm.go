package vm

import (
	"fmt"
	"lang_vm/code"
	"lang_vm/object"
)

type VM struct {
	ins   code.Instructions
	stack []object.Object
	sp    int
	ip    int
}

func (vm *VM) Run() error {
	var err error
	running := true
	
	for running {
		opcode := code.OpCode(vm.ins[vm.ip])
		switch opcode {

		case code.OpConstant:
			val := code.ReadUint16(vm.ins[vm.ip+1:])
			vm.sp += 2
			err = vm.push(&object.Integer{Value: int64(val)})
			if err != nil {
				running = false
			}

		default:
			err = fmt.Errorf("unkown opcode: %d", opcode)
			running = false
		}
	}

	return err
}

func (vm *VM) push(o object.Object) error {
	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}
