package parser

import (
	"testing"

	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/bytecode"
	"github.com/djalmahenry/go-interpreter/pkg/domain/entities/cp"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/instructions"
	"github.com/djalmahenry/go-interpreter/pkg/domain/usecases/lexer"
	"github.com/stretchr/testify/assert"
)

func TestValidEmptyAlgorithm(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
fim`
	l := lexer.New(alg)
	p := New(l)
	err := p.Parser()
	assert.Nil(t, err)
}

func TestValidHelloWorldAlgorithm(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima("Olá mundo!");
fim`
	l := lexer.New(alg)
	p := New(l)
	err := p.Parser()
	assert.Nil(t, err)
}

func TestValidHelloWorldWithTwoSentences(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima("Olá...");
	imprima("Mundo!");
fim`
	l := lexer.New(alg)
	p := New(l)
	err := p.Parser()
	assert.Nil(t, err)
}

func TestBytecodeEmptyAlgorithm(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
fim`
	l := lexer.New(alg)
	p := New(l)
	bc := bytecode.New()
	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, bc, p.GetBC())
}

func TestBytecodeFunctionCallWithoutArguments(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima();
fim`
	// CP:
	//    0: STR "io.println"
	//    1: INT 0
	// CODE:
	//    LDC  1 (0)
	//    CALL 0 (io.println)

	expectedCp := cp.New()
	printlnIndex := expectedCp.Add("io.println")
	argsCountIndex := expectedCp.Add(0)

	l := lexer.New(alg)
	p := New(l)
	expectedBc := bytecode.New()
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)

	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, expectedCp, p.GetCP())
	assert.Equal(t, expectedBc, p.GetBC())
}

func TestBytecodeHelloWorldAlgorithm(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima("Olá mundo!");
fim`
	// CP:
	//    0: STR "io.println"
	//    1: STR "Olá mundo!"
	//    2: INT 1
	// CODE:
	//    LDC 1 (Olá mundo!)
	//    LDC 2 (1)
	//    CALL 0 (io.println)

	expectedCp := cp.New()
	printlnIndex := expectedCp.Add("io.println")
	messageIndex := expectedCp.Add("Olá mundo!")
	argsCountIndex := expectedCp.Add(1)

	l := lexer.New(alg)
	p := New(l)
	expectedBc := bytecode.New()
	expectedBc.Add(instructions.LDC, messageIndex)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)

	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, expectedCp, p.GetCP())
	assert.Equal(t, expectedBc, p.GetBC())
}

func TestBytecodeHelloWorldWithTwoWrites(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima("Olá...");
	imprima("mundo!");
fim`
	// CP:
	//    0: STR "io.println"
	//    1: STR "Olá..."
	//    2: INT 1
	//    3: STR "mundo!"
	// CODE:
	//    LDC 1 (Olá...)
	//    LDC 2 (1)
	//    CALL 0 (io.println)
	//    LDC 3 (mundo!)
	//    LDC 2 (1)
	//    CALL 0 (io.println)

	expectedCp := cp.New()
	printlnIndex := expectedCp.Add("io.println")
	messageIndex1 := expectedCp.Add("Olá...")
	argsCountIndex := expectedCp.Add(1)
	messageIndex2 := expectedCp.Add("mundo!")

	l := lexer.New(alg)
	p := New(l)
	expectedBc := bytecode.New()
	expectedBc.Add(instructions.LDC, messageIndex1)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)
	expectedBc.Add(instructions.LDC, messageIndex2)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)

	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, expectedCp, p.GetCP())
	assert.Equal(t, expectedBc, p.GetBC())
}

func TestInvalidCompleteAlgorithmDeclarationWithoutId(t *testing.T) {
	alg :=
		`algoritmo ;
início
	imprima("Olá...");
fim`
	l := lexer.New(alg)
	p := New(l)

	err := p.Parser()
	assert.EqualError(t, err, "Expected IDENT")
}

func TestInvalidCompleteAlgorithmDeclarationWithoutSemicolon(t *testing.T) {
	alg :=
		`algoritmo ola
início
	imprima("Olá...");
fim`
	l := lexer.New(alg)
	p := New(l)

	err := p.Parser()
	assert.EqualError(t, err, "Expected SEMICOLON")
}

func TestBytecodeHelloWorldWithInput(t *testing.T) {
	alg :=
		`algoritmo qual_o_seu_nome;

		variáveis
			nome: literal;
		fim-variáveis
		
		início
			imprima("Qual o seu nome?");
			nome := leia();
			imprima("Olá ");
			imprima(nome);
		fim
		`
	// CP:
	//    0: STR "io.println"
	//    1: STR "Qual o seu nome?"
	//    2: INT 1
	//    3: STR "io.readln"
	//    4: STR "Olá "
	// VAR:
	//    0: STR "nome"
	// CODE:
	//    LDC 1 (Qual o seu nome?)
	//    LDC 2 (1)
	//    CALL 0 (io.println)
	//    CALL 3 (io.readln)
	//    STV 0 (nome)
	//    LDC 4 (Olá )
	//    LDC 2 (1)
	//    CALL 0 (io.println)
	//    LDV 0 (nome)
	//    LDC 3 (1)
	//    CALL 0 (io.println)

	expectedCp := cp.New()
	printlnIndex := expectedCp.Add("io.println")
	messageIndex1 := expectedCp.Add("Qual o seu nome?")
	argsCountIndex := expectedCp.Add(1)
	readlnIndex := expectedCp.Add("io.readln")
	messageIndex2 := expectedCp.Add("Olá ")

	l := lexer.New(alg)
	p := New(l)
	expectedBc := bytecode.New()
	expectedBc.Add(instructions.LDC, messageIndex1)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)
	expectedBc.Add(instructions.CALL, readlnIndex)
	expectedBc.Add(instructions.STV, 0)
	expectedBc.Add(instructions.LDC, messageIndex2)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)
	expectedBc.Add(instructions.LDV, 0)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)

	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, expectedCp, p.GetCP())
	assert.Equal(t, expectedBc, p.GetBC())
}

func TestBytecodeHelloWorldWithTwoArgs(t *testing.T) {
	alg :=
		`algoritmo olá_mundo;
início
	imprima("Olá...", "mundo!");
fim`
	// CP:
	//    0: STR "io.println"
	//    1: STR "Olá..."
	//    2: STR "mundo!"
	//    3: INT 2
	// CODE:
	//    LDC 1 (Olá...)
	//    LDC 2 (mundo!)
	//    LDC 3 (2)
	//    CALL 0 (io.println)

	expectedCp := cp.New()
	printlnIndex := expectedCp.Add("io.println")
	messageIndex1 := expectedCp.Add("Olá...")
	messageIndex2 := expectedCp.Add("mundo!")
	argsCountIndex := expectedCp.Add(2)

	l := lexer.New(alg)
	p := New(l)
	expectedBc := bytecode.New()
	expectedBc.Add(instructions.LDC, messageIndex1)
	expectedBc.Add(instructions.LDC, messageIndex2)
	expectedBc.Add(instructions.LDC, argsCountIndex)
	expectedBc.Add(instructions.CALL, printlnIndex)

	err := p.Parser()
	assert.Nil(t, err)
	assert.Equal(t, expectedCp, p.GetCP())
	assert.Equal(t, expectedBc, p.GetBC())
}
