package bce

import (
	"fmt"

	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/bytecode"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/cp"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/stack"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/vars"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions/calli"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions/ldci"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions/ldvi"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions/nopi"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions/stvi"
)

// BytecodeExecutor is responsible for execute the bytecode.
type BytecodeExecutor struct {
	bc           *bytecode.Bytecode
	instructions map[int]instructions.InstructionImplementation
	ip           int
}

// New creates a new BytecodeExecutor.
func New(bc *bytecode.Bytecode) *BytecodeExecutor {
	bce := &BytecodeExecutor{
		ip:           0,
		bc:           bc,
		instructions: make(map[int]instructions.InstructionImplementation),
	}

	bce.instructions[instructions.NOP] = nopi.New()
	bce.instructions[instructions.LDC] = ldci.New()
	bce.instructions[instructions.CALL] = calli.New()
	bce.instructions[instructions.LDV] = ldvi.New()
	bce.instructions[instructions.STV] = stvi.New()

	return bce
}

// Run receives the context (constant pool, vars, stack, stdin and stdout) and runs the bytecode.
func (bce *BytecodeExecutor) Run(cp *cp.CP, vars *vars.Vars, st *stack.Stack, stdin instructions.StdinInterface, stdout instructions.StdoutInterface) error {
	for {
		opcode, err := bce.Next()
		if err != nil {
			return nil
		}

		instruction, exist := bce.instructions[opcode]
		if !exist {
			return fmt.Errorf("Invalid opcode %d", opcode)
		}

		err = instruction.FetchOperands(bce)
		if err != nil {
			return err
		}

		err = instruction.Execute(cp, vars, st, stdin, stdout)
		if err != nil {
			return err
		}
	}
}

// Next returns the next bytecode.
func (bce *BytecodeExecutor) Next() (code int, err error) {
	code, err = bce.bc.Get(bce.ip)
	bce.ip++

	return
}
