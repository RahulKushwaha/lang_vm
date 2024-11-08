package vm

import (
	"fmt"
	"lang_vm/code"
	"lang_vm/object"
)

const (
	maxStackSize = 2048
)

type VM struct {
	ins   code.Instructions
	stack []object.Object
	sp    int
	ip    int
}

func New(ins code.Instructions) *VM {
	return &VM{
		ins:   ins,
		stack: make([]object.Object, maxStackSize),
		sp:    0,
		ip:    0,
	}
}

func (vm *VM) Run() error {
	var err error
	running := true

	for running {
		opcode := code.OpCode(vm.ins[vm.ip])
		vm.ip++
		switch opcode {

		case code.OpConstant:
			val := code.ReadUint16(vm.ins[vm.ip:])
			vm.ip += 2
			err = vm.push(&object.Integer{Value: int64(val)})
			if err != nil {
				running = false
			}

		case code.OpAdd:
			if err = vm.add(); err != nil {
				running = false
			}

		case code.OpHalt:
			running = false

		default:
			err = fmt.Errorf("unkown opcode: %d", opcode)
			running = false
		}
	}

	return err
}

func (vm *VM) push(o object.Object) error {
	vm.stack[vm.sp] = o
	if vm.sp-1 >= maxStackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.sp++
	return nil
}

func (vm *VM) Pop() object.Object {
	vm.sp--
	o := vm.stack[vm.sp]

	return o
}

func (vm *VM) add() error {
	a := vm.Pop()
	b := vm.Pop()

	temp := object.Integer{Value: 0}

	if a.Type() != temp.Type() || b.Type() != temp.Type() {
		return fmt.Errorf("%v or %v is not of integer type", a.Inspect(), b.Inspect())
	}

	c := a.(*object.Integer).Value + b.(*object.Integer).Value
	o := &object.Integer{Value: c}

	if err := vm.push(o); err != nil {
		return err
	}

	return nil
}
