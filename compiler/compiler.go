package compiler

import (
	"fmt"
	"lang_vm/ast"
	"lang_vm/code"
	"lang_vm/object"
)

type Compiler struct {
	ins         code.Instructions
	scopes      []CompilationScope
	scopeIndex  int
	symbolTable *SymbolTable
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type CompilationScope struct {
	ins code.Instructions
}

func NewCompiler() *Compiler {
	c := &Compiler{
		ins:         code.Instructions{},
		scopes:      make([]CompilationScope, 0),
		scopeIndex:  0,
		symbolTable: NewSymbolTable(),
	}

	c.scopes = append(c.scopes, CompilationScope{ins: code.Instructions{}})

	return c
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

	case *ast.Identifier:
		if err := c.compileIdentifier(n); err != nil {
			return err
		}

	case *ast.LetStatement:
		if err := c.compileLetStatement(n); err != nil {
			return err
		}

	case *ast.ExpressionStatement:
		if err := c.Compile(n.Expression); err != nil {
			return err
		}

	case *ast.IfExpression:
		if err := c.compileIfExpression(*n); err != nil {
			return err
		}

	case *ast.BinaryExpression:
		if err := c.compileBinaryExpression(*n); err != nil {
			return err
		}

	case *ast.IntegerLiteral:
		c.emit(code.OpConstant, int(n.Value))

	default:
		return fmt.Errorf("unknown node type: %T", n)
	}

	return nil
}

func (c *Compiler) compileIdentifier(n *ast.Identifier) error {
	symbol, ok := c.symbolTable.Resolve(n.Value)
	if !ok {
		return fmt.Errorf("identifier could not be resolved: %s", n.Value)
	}

	c.emit(code.OpLoad, symbol.Index)

	return nil
}

func (c *Compiler) compileLetStatement(n *ast.LetStatement) error {
	symbol := c.symbolTable.Define(n.Name.String())
	if err := c.Compile(n.Value); err != nil {
		return err
	}

	c.emit(code.OpStore, symbol.Index)

	return nil
}

func (c *Compiler) compileBinaryExpression(n ast.BinaryExpression) error {
	if err := c.Compile(n.Left); err != nil {
		return err
	}

	if err := c.Compile(n.Right); err != nil {
		return err
	}

	switch n.Operator {
	case "+":
		c.emit(code.OpAdd)
	case "-":
		c.emit(code.OpSub)
	case "*":
		c.emit(code.OpMul)
	case "/":
		c.emit(code.OpDiv)
	case "==":
		c.emit(code.OpEqual)
	case "!=":
		c.emit(code.OpNotEqual)
	case ">":
		c.emit(code.OpGreaterThan)

	default:
		return fmt.Errorf("unsupported operator: %s", n.Operator)
	}
	return nil
}

func (c *Compiler) compileIfExpression(n ast.IfExpression) error {
	if err := c.Compile(n.Condition); err != nil {
		return err
	}

	jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

	if err := c.Compile(n.Consequence); err != nil {
		return err
	}

	jumpPos := c.emit(code.OpJump, 9999)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	if n.Alternative == nil {
		c.emit(code.OpNull)
	} else {
		if err := c.Compile(n.Alternative); err != nil {
			return err
		}
	}

	afterAlternative := len(c.currentInstructions())
	c.changeOperand(jumpPos, afterAlternative)
	return nil
}

func (c *Compiler) emit(op code.OpCode, operands ...int) int {
	ins := code.Make(op, operands...)
	return c.addInstruction(ins)
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())

	updatedInstructions := append(c.currentInstructions(), ins...)
	c.scopes[c.scopeIndex].ins = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].ins
}

func (c *Compiler) replaceInstruction(pos int, newIns []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newIns); i++ {
		ins[pos+i] = newIns[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.OpCode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) ByteCode() *ByteCode {
	b := &ByteCode{
		Instructions: c.currentInstructions(),
	}

	return b
}
