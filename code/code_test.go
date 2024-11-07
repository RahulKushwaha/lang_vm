package code

import (
	"fmt"
	"testing"
)

func TestOpPrint(t *testing.T) {
	ins := Make(OpConstant, 1)
	fmt.Println(ins)
}
