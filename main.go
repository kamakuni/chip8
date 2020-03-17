package main

import (
	"fmt"
)

type Emulator struct {
	opcode     uint16        // two bytes opcodes
	memory     [4096]uint8   // 4K memory
	v          [16]uint8     // 15 8-bit registers for general purpose and one for "carry-flag"
	i          uint16        // index register
	pc         uint16        // program counter
	gfx        [64 * 32]bool // 2048 black or white pixels
	delayTimer uint8         // Timer registor for general purpose
	soundTimer uint8         // Timer registor for sound
	stack      [16]uint16    // to store current pc
	sp         uint16        // stack pointer
	key        [16]uint8     // to store current stats of key
}

// Create Emulator
func NewEmulator() *Emulator {
	return &Emulator{
		pc:     0x200,
		opcode: 0,
		i:      0,
		sp:     0,
	}
}

// Print Emulator status
func (e Emulator) Print() {
	fmt.Printf("opcode:%v\n", e.opcode)
	fmt.Printf("memory:%v\n", e.memory)
	fmt.Printf("v:%v\n", e.v)
	fmt.Printf("i:%v\n", e.i)
	fmt.Printf("pc:%v\n", e.pc)
	fmt.Printf("gfx:%v\n", e.gfx)
	fmt.Printf("delayTimer:%v\n", e.delayTimer)
	fmt.Printf("soundTimer:%v\n", e.soundTimer)
	fmt.Printf("stack:%v\n", e.stack)
	fmt.Printf("sp:%v\n", e.sp)
	fmt.Printf("key:%v\n", e.key)
}

func main() {
	fmt.Println("start")
	emu := NewEmulator()
	emu.Print()
}
