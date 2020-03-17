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

func main() {
	fmt.Println("start")
}
