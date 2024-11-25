package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type OpCode byte

const (
	OpConstant OpCode = iota
	OpPop
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpJump
	OpJumpNotTruthy
	OpNull
	OpGreaterThan
	OpEqual
	OpNotEqual
	OpLoad
	OpStore
	OpHalt
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpPop:           {"OpPop", []int{}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpLoad:          {"OpLoad", []int{2}},
	OpStore:         {"OpStore", []int{2}},
	OpHalt:          {"OpHalt", []int{}},
}

type Instructions []byte

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("unknown opcode 0x%x", op)
	}

	return def, nil
}

type Builder struct {
	ins Instructions
}

func NewBuilder() *Builder {
	return &Builder{ins: Instructions{}}
}

func (b *Builder) Add(op OpCode, operands ...int) *Builder {
	b.ins = append(b.ins, Make(op, operands...)...)
	return b
}

func (b *Builder) Build() Instructions {
	return b.ins
}

func Make(op OpCode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

func ReadOperands(definition *Definition, instructions Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0

	for i, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instructions[offset:]))
		case 1:
			operands[i] = int(ReadUint8(instructions[offset:]))
		}

		offset += width
	}
	return operands, offset
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			_, err := fmt.Fprintf(&out, "ERROR: %s\n", err)
			if err != nil {
				log.Fatalf("failed to write to buffer, %v", err)
			}
		}

		operands, offset := ReadOperands(def, ins[i+1:])

		_, err = fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		if err != nil {
			log.Fatalf("failed to write to buffer, %v", err)
		}

		i += 1 + offset
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(definition *Definition, operands []int) string {
	count := len(operands)

	if count != len(definition.OperandWidths) {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), count)
	}

	switch count {
	case 0:
		return definition.Name
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", definition.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
}
