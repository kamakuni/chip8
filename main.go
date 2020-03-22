package main

import (
	"encoding/binary"
	"fmt"
	sdl "github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
)

// NewFonts creates fonts array
func NewFonts() [80]uint8 {
	//
	// https://github.com/pierreyoda/rust-chip8/blob/master/src/display.rs
	//
	// Chip8 font set.
	// Each number or character is 4x5 pixels and is stored as 5 bytes.
	// In each byte, only the first nibble (the first 4 bytes) is used.
	// For instance, with the number 3 :
	//  hex    bin     ==> drawn pixels
	// 0xF0  1111 0000        ****
	// 0X10  0001 0000           *
	// 0xF0  1111 0000        ****
	// 0x10  0001 0000           *
	// 0xF0  1111 0000        ****
	return [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
}

type Emulator struct {
	Opcode     uint16       // two bytes opcodes
	Memory     [4096]uint8  // 4K memory
	v          [16]uint8    // 15 8-bit registers for general purpose and one for "carry-flag"
	I          uint16       // index register
	Pc         uint16       // program counter
	Gfx        [32][64]bool // 2048 black or white pixels
	delayTimer uint8        // Timer registor for general purpose
	soundTimer uint8        // Timer registor for sound
	Stack      [16]uint16   // to store current pc
	Sp         uint16       // stack pointer
	key        [16]uint8    // to store current stats of key
	ShouldDraw bool
}

// NewEmulator creates Emulator
func NewEmulator(fonts [80]uint8) *Emulator {
	var memory [4096]uint8
	for i, font := range fonts {
		memory[i] = font
	}
	return &Emulator{
		Pc:     0x200,
		Opcode: 0,
		Memory: memory,
		I:      0,
		Sp:     0,
	}
}

// Load loads data to memory
func (e *Emulator) Load(data []byte) {
	for i, b := range data {
		e.Memory[int(e.Pc)+i] = b
		fmt.Printf("load byte:0x%x", e.Memory[int(e.Pc)+i])
	}
}

func (e *Emulator) load(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		// TODO:logging
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		// TODO:logging
		panic(err)
	}
	fmt.Printf("file size:%v\n", stat.Size())

	buf := make([]byte, stat.Size())
	err = binary.Read(file, binary.BigEndian, &buf)
	if err != nil {
		// TODO:logging
		panic(err)
	}

	for i, b := range buf {
		e.Memory[int(e.Pc)+i] = b
		fmt.Printf("%x", e.Memory[int(e.Pc)+i])
	}
}

func (e *Emulator) Fetch() uint16 {
	op1 := uint16(e.Memory[int(e.Pc)])
	op2 := uint16(e.Memory[int(e.Pc)+1])
	return op1<<8 | op2
}

func (e *Emulator) Decode(opcode uint16) {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x000:
			e.Gfx = [32][64]bool{}
			e.ShouldDraw = true
			e.Pc += 2
		}
		break
	case 0xA000:
		e.I = opcode & 0x0FFF
		e.Pc += 2
		break
	default:
		log.Fatalf("Unexpected opcode 0x%x", opcode)
	}
}

// Print Emulator status
func (e *Emulator) Print() {
	fmt.Printf("opcode:%v\n", e.Opcode)
	fmt.Printf("memory:%v\n", e.Memory)
	fmt.Printf("v:%v\n", e.v)
	fmt.Printf("i:%v\n", e.I)
	fmt.Printf("pc:%v\n", e.Pc)
	fmt.Printf("gfx:%v\n", e.Gfx)
	fmt.Printf("delayTimer:%v\n", e.delayTimer)
	fmt.Printf("soundTimer:%v\n", e.soundTimer)
	fmt.Printf("stack:%v\n", e.Stack)
	fmt.Printf("sp:%v\n", e.Sp)
	fmt.Printf("key:%v\n", e.key)
}

func main() {

	fmt.Println("start")
	if len(os.Args) != 2 {
		// TODO:logging
		panic(nil)
	}
	filepath := os.Args[1]
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	emu.load(filepath)
	emu.Print()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		// TODO:logging
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP-8 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		64, 32, sdl.WINDOW_SHOWN)
	if err != nil {
		// TODO:logging
		panic(err)
	}
	defer window.Destroy()

	opcode := emu.Fetch()
	emu.Decode(opcode)

}
