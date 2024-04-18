package main

import (
	"fmt"
	"os"
)

type Memoria []byte

type CPU struct {
	mem        Memoria
	tamanhoMem uint64

	// registradores
	rax, rbx, rcx, rdx int64
	rpc                uint64
	zf                 bool
    flag_pulando bool
}

func (c *CPU) dumpRegistradores() {
	fmt.Printf("RAX: %02x\nRBX: %02x\nRCX: %02x\nRDX: %02x\nRPC: %02x\nZF: %v\n\n", c.rax, c.rbx, c.rcx, c.rdx, c.rpc, c.zf)
}

func (c *CPU) getValorMemoria(pos uint64) byte {
	if pos < 0 || pos > c.tamanhoMem-1 {
		fmt.Printf("Tentativa de aceeso em posição inválida de memória (%x), saindo do programa\n", pos)
		os.Exit(1)
	}
	return c.mem[pos]
}

func (c *CPU) setValorMemoria(pos uint64, valor byte) {
	if pos < 0 || pos > c.tamanhoMem-1 {
		fmt.Println("Posição inválida para colocar na memória")
		os.Exit(1)
	}
	c.mem[pos] = valor
}

func (c *CPU) carregarPrograma(nomeArq string) {
	var err error
	c.mem, err = os.ReadFile(nomeArq)
	if err != nil {
		fmt.Println("Erro abrindo arquivo, saindo do programa")
		os.Exit(1)
	}
	c.tamanhoMem = uint64(len(c.mem))
    fmt.Printf("Arquivo carregado, tamanho da memória: %v\n", c.tamanhoMem)
}

func (c *CPU) escolherReg(reg byte) *int64 {
	switch reg {
	case 0x02:
		return &c.rax
	case 0x03:
		return &c.rbx
	case 0x04:
		return &c.rcx
	case 0x05:
		return &c.rdx
	}
	return nil
}

func (c *CPU) executarInstrucao(instrucao byte) {
	fmt.Printf("Intrução atual: %x\n", instrucao)

	var par1 byte
	var par2 byte

	switch instrucao {

    case 0x00: // ADD REG, BYTE
        par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) += int64(par2)
        c.rpc += 3
        fmt.Printf("Instrução ADD REG, BYTE | %x <- %x\n", par1, par2)
        break

    case 0x01: // ADD REG, REG
		par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) += *c.escolherReg(par2)
        c.rpc += 3
        fmt.Printf("Instrução ADD REG, REG | %x <- %x\n", par1, par2)
        break

    case 0x10: // INC REG
        par1 = c.getValorMemoria(c.rpc + 1)
        *c.escolherReg(par1)++
        c.rpc += 2
        fmt.Printf("Instrução INC REG | %x++\n", par1)
        break

    case 0x20: // DEC REG
        par1 = c.getValorMemoria(c.rpc + 1)
        *c.escolherReg(par1)--
        c.rpc += 2
        fmt.Printf("Instrução DEC REG | %x--\n", par1)
        break

    case 0x30: // SUB REG, BYTE
        par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) -= int64(par2)
        c.rpc += 3
        fmt.Printf("Instrução SUB REG, BYTE | %x <- %x\n", par1, par2)
        break

    case 0x31: // SUB REG, REG
        par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) -= *c.escolherReg(par2)
        c.rpc += 3
        fmt.Printf("Instrução SUB REG, SUB | %x <- %x\n", par1, par2)
        break

    case 0x40: // MOV REG, BYTE
		par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) = int64(par2)
		c.rpc += 3
		fmt.Printf("Instrução MOV REG, BYTE | %x <- %x\n", par1, par2)
		break

    case 0x41: // MOV REG, REG
		par1 = c.getValorMemoria(c.rpc + 1)
		par2 = c.getValorMemoria(c.rpc + 2)
        *c.escolherReg(par1) = *c.escolherReg(par2)
		c.rpc += 3
		fmt.Printf("Instrução MOV REG, REG | %x <- %x\n", par1, par2)
		break

    case 0x50: // JMP BYTE
        par1 = c.getValorMemoria(c.rpc + 1)
        c.rpc = uint64(par1)
        fmt.Printf("Instrução JMP BYTE | RPC <- %x\n", par1)
        break

    case 0x60: // CMP REG, BYTE
        par1 = c.getValorMemoria(c.rpc + 1)
        par2 = c.getValorMemoria(c.rpc + 2)
        c.zf = *c.escolherReg(par1) == int64(par2)
        c.rpc += 3
        fmt.Printf("Instrução CMP REG, BYTE | %x == %x\n", *c.escolherReg(par1), par2)
        break
        
    case 0x61: // CMP REG, REG
        par1 = c.getValorMemoria(c.rpc + 1)
        par2 = c.getValorMemoria(c.rpc + 2)
        c.zf = *c.escolherReg(par1) == *c.escolherReg(par2)
        c.rpc += 3
        fmt.Printf("Instrução CMP REG, REG | %x == %x\n", *c.escolherReg(par1), *c.escolherReg(par2))
        break

    case 0x70: // JZ BYTE
        par1 = c.getValorMemoria(c.rpc + 1)
        if c.zf {
            c.rpc = uint64(par1)
        } else {
            c.rpc += 2
        }
        fmt.Printf("Instrução JZ BYTE | ZF ? PC = %x : PC += 2\n", par1)

	default:
		fmt.Println("Instrução não suportada, saindo do programa")
		os.Exit(1)
	}

}

func (c *CPU) rodarPrograma() {
	rodando := true

	// loop principal
	for rodando {
		instrucao := c.getValorMemoria(c.rpc)
		c.executarInstrucao(instrucao)
		c.dumpRegistradores()
	}
}

func main() {
	fmt.Println("Começo de CPU")

	var cpu CPU
	cpu.carregarPrograma(os.Args[1])
	cpu.dumpRegistradores()
	cpu.rodarPrograma()
}
