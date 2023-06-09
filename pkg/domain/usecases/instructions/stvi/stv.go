package stvi

import (
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/cp"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/stack"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/vars"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions"
)

var _ instructions.InstructionImplementation = &STVInst{}

// STVInst is responsible for pop a value from the stack and put into a var.
type STVInst struct {
	instructions.Instruction
	VarIndex int
}

// New creates a new STVInst.
func New() *STVInst {
	return &STVInst{
		instructions.Instruction{
			Opcode: instructions.STV,
		},
		0,
	}
}

// FetchOperands gets the opcode operands.
func (i *STVInst) FetchOperands(fetch instructions.FetchOperandsImplementation) error {
	var err error
	i.VarIndex, err = fetch.Next()

	return err
}

// Execute receives the context and runs the opcode.
func (i *STVInst) Execute(cp *cp.CP, vars *vars.Vars, st *stack.Stack, stdin instructions.StdinInterface, stdout instructions.StdoutInterface) error {
	value, err := st.Pop()
	if err != nil {
		return err
	}

	vars.Set(i.VarIndex, value)

	return nil
}
