package instructions

import (
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/cp"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/stack"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/vars"
)

const (
	// NOP (No Operation)
	NOP int = iota
	// LDC (Load Constant)
	LDC
	// ADD (Addition)
	ADD
	// CALL (Call Function)
	CALL
	// STV (Store Value)
	STV
	// LDV (Load Value)
	LDV
)

// Instruction has the common data of all instructions.
type Instruction struct {
	Opcode int
}

// InstructionImplementation has the minimal interface to be a valid instruction.
type InstructionImplementation interface {
	FetchOperands(fetch FetchOperandsImplementation) error
	Execute(cp *cp.CP, vars *vars.Vars, st *stack.Stack, stdin StdinInterface, stdout StdoutInterface) error
}

// FetchOperandsImplementation has the minimal interface to be a valid BCE.
type FetchOperandsImplementation interface {
	Next() (code int, err error)
}
