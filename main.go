package main

import (
	"encoding/binary"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
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

func NewKeyMap() map[int]byte {
	return map[int]byte{
		sdl.SCANCODE_1: 0x1,
		sdl.SCANCODE_2: 0x2,
		sdl.SCANCODE_3: 0x3,
		sdl.SCANCODE_4: 0xc,
		sdl.SCANCODE_Q: 0x4,
		sdl.SCANCODE_W: 0x5,
		sdl.SCANCODE_E: 0x6,
		sdl.SCANCODE_R: 0xd,
		sdl.SCANCODE_A: 0x7,
		sdl.SCANCODE_S: 0x8,
		sdl.SCANCODE_D: 0x9,
		sdl.SCANCODE_F: 0xe,
		sdl.SCANCODE_Z: 0xa,
		sdl.SCANCODE_X: 0x0,
		sdl.SCANCODE_C: 0xb,
		sdl.SCANCODE_V: 0xf,
	}
}

type Emulator struct {
	Opcode     uint16      // two bytes opcodes
	Memory     [4096]uint8 // 4K memory
	V          [16]uint8   // 15 8-bit registers for general purpose and one for "carry-flag"
	I          uint16      // index register
	Pc         uint16      // program counter
	Gfx        [2048]uint8 // 2048 black or white pixels
	delayTimer uint8       // Timer registor for general purpose
	soundTimer uint8       // Timer registor for sound
	Stack      [16]uint16  // to store current pc
	Sp         uint16      // stack pointer
	keys       [16]bool    // to store current stats of key
	keyMap     map[int]byte
	shouldDraw bool
	surface    *sdl.Surface
	window     *sdl.Window
}

// NewEmulator creates Emulator
func NewEmulator(fonts [80]uint8) *Emulator {
	var memory [4096]uint8
	for i, font := range fonts {
		memory[i] = font
	}
	keyMap := NewKeyMap()
	return &Emulator{
		Pc:     0x200,
		Opcode: 0,
		Memory: memory,
		I:      0,
		Sp:     0,
		keyMap: keyMap,
	}
}

func (e *Emulator) InitDisplay() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}
	//defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}
	//defer window.Destroy()

	window.Raise()
	e.window = window

	// window has been created, now need to get the window surface to draw on window
	surface, err := window.GetSurface()
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create surface: %s\n", err)
		os.Exit(2)
	}
	e.surface = surface
}

func (e *Emulator) DestroyDisplay() {
	sdl.Quit()
	e.window.Destroy()
}

func (e *Emulator) draw() {
	for i := range e.Gfx {
		x := int32(i % 64)
		y := int32(int(i / 64))
		rect := sdl.Rect{x * 10, y * 10, 10, 10}
		if e.Gfx[i] == 1 {
			e.surface.FillRect(&rect, sdl.MapRGB(e.surface.Format, 200, 200, 200))
		} else {
			e.surface.FillRect(&rect, sdl.MapRGB(e.surface.Format, 35, 35, 35))
		}
	}
	e.window.UpdateSurface()
}

func (e *Emulator) next() {
	e.Pc += 2
}

func (e *Emulator) skip() {
	e.Pc += 4
}

func (e *Emulator) jump(next uint16) {
	e.Pc = next
}

func (e *Emulator) Load(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("file size:%v\n", stat.Size())

	buf := make([]byte, stat.Size())
	err = binary.Read(file, binary.BigEndian, &buf)
	if err != nil {
		log.Fatalln(err)
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

func (e *Emulator) Exec(opcode uint16) {
	// https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Instruction-Set
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00E0:
			// CLS: Clear the screen
			e.Gfx = [2048]uint8{}
			e.shouldDraw = true
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x00EE:
			e.Sp--
			e.Pc = e.Stack[e.Sp]
			e.next()
		default:
			log.Fatalf("Unexpected opcode 0x%x\n", opcode)
		}
	case 0x1000:
		// Goto NNN: Jump to address NNN
		e.Pc = opcode & 0x0FFF
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x2000:
		// CALL: Call the subroutine at address NNN
		e.Stack[e.Sp] = e.Pc
		e.Sp++
		e.Pc = opcode & 0x0FFF
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x3000:
		// skips the next instruction if VX equals NN.
		// (Usually the next instruction is a jump to skip a code block)
		x := opcode & 0x0F00 >> 8
		log.Printf("VF: %v\n", e.V[x])
		log.Printf("NN: %v\n", opcode&0x00FF)
		if int(e.V[x]) == int(opcode&0x00FF) {
			e.skip()
		} else {
			e.next()
		}
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x4000:
		// skips the next instruction if VX doesn't equal NN.
		// (Usually the next instruction is a jump to skip a code block)
		x := opcode & 0x0F00 >> 8
		if int(e.V[x]) != int(opcode&0x00FF) {
			e.skip()
		} else {
			e.next()
		}
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x5000:
		// skips the next instruction if VX equals VY.
		// (Usually the next instruction is a jump to skip a code block)
		x := opcode & 0x0F00 >> 8
		y := opcode & 0x00F0 >> 4
		if e.V[x] == e.V[y] {
			e.skip()
		} else {
			e.next()
		}
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x6000:
		// Sets VX to NN.
		x := opcode & 0x0F00 >> 8
		if x == 0 && uint8(opcode&0x00FF) == 1 {
			log.Printf("VX %v", uint8(opcode&0x00FF))
		}
		e.V[x] = uint8(opcode & 0x00FF)
		e.next()
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x7000:
		// 	Adds NN to VX. (Carry flag is not changed)
		x := opcode & 0x0F00 >> 8
		e.V[x] += uint8(opcode & 0x00FF)
		e.next()
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0x8000:
		switch opcode & 0x000F {
		case 0:
			// Sets VX to the value of VY.
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			e.V[x] = e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 1:
			// 	Sets VX to VX or VY. (Bitwise OR operation)
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			e.V[x] = e.V[x] | e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 2:
			// Sets VX to VX and VY. (Bitwise AND operation)
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			e.V[x] = e.V[x] & e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 3:
			// Sets VX to VX xor VY.
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			e.V[x] = e.V[x] ^ e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 4:
			// Add the value of register VY to register VX
			// Set VF to 01 if a carry occurs
			// Set VF to 00 if a carry does not occur
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			if uint16(e.V[x])+uint16(e.V[y]) > 0xFF {
				e.V[0xF] = 0x1
			} else {
				e.V[0xF] = 0x0
			}
			e.V[x] += e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 5:
			// Subtract the value of register VY from register VX
			// Set VF to 00 if a borrow occurs
			// Set VF to 01 if a borrow does not occur
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			if e.V[x] < e.V[y] {
				e.V[0xF] = 0x0
			} else {
				e.V[0xF] = 0x1
			}
			e.V[x] -= e.V[y]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 6:
			// Store the value of register VY shifted right one bit in register VX¹
			// Set register VF to the least significant bit prior to the shift
			// VY is unchanged
			x := opcode & 0x0F00 >> 8
			if (e.V[x] & 0x01) == 1 {
				e.V[0xF] = 0x1
			} else {
				e.V[0xF] = 0x0
			}
			e.V[x] >>= 1
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 7:
			// Set register VX to the value of VY minus VX
			// Set VF to 00 if a borrow occurs
			// Set VF to 01 if a borrow does not occur
			x := opcode & 0x0F00 >> 8
			y := opcode & 0x00F0 >> 4
			if e.V[y]-e.V[x] < 0 {
				e.V[0xF] = 0x0
			} else {
				e.V[0xF] = 0x1
			}
			e.V[x] = e.V[y] - e.V[x]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0xE:
			// Store the value of register VY shifted left one bit in register VX¹
			// Set register VF to the most significant bit prior to the shift
			// VY is unchanged
			x := opcode & 0x0F00 >> 8
			if e.V[x]>>7 == 1 {
				e.V[0xF] = 0x1
			} else {
				e.V[0xF] = 0x0
			}
			e.V[x] <<= 1
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		default:
			log.Fatalf("Unexpected opcode 0x%x\n", opcode)
		}
	case 0x9000:
		x := opcode & 0x0F00 >> 8
		y := opcode & 0x00F0 >> 4
		if e.V[x] != e.V[y] {
			e.skip()
		} else {
			e.next()
		}
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0xA000:
		// LD: Sets I to the address NNN.
		e.I = opcode & 0x0FFF
		e.next()
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0xB000:
		e.Pc = opcode&0x0FFF + uint16(e.V[0])
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0xC000:
		x := opcode & 0x0F00 >> 8
		mask := opcode & 0x00FF
		e.V[x] = uint8(rand.Uint32() & uint32(mask))
		e.next()
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0xD000:
		vx := e.V[opcode&0x0F00>>8]
		vy := e.V[opcode&0x00F0>>4]
		height := opcode & 0x000F
		e.V[0xF] = 0
		for yi := 0; yi < int(height); yi++ {
			row := e.Memory[int(e.I)+yi]
			for xi := 0; xi < 8; xi++ {
				// 1000 0000 >> xi
				if row&(0x80>>uint8(xi)) != 0 {
					x := int(vx) + xi
					y := int(vy) + yi
					// allow for wrapping
					// https://www.reddit.com/r/EmuDev/comments/aar9nb/chip_8_emulator_collision_detection_not_working/
					if x >= 64 {
						x %= 64
					}
					if y >= 32 {
						y %= 32
					}
					if e.Gfx[x+y*64] == 1 {
						// when collision detected
						e.V[0xF] = 1
					} else {
						e.V[0xF] = 0
					}
					e.Gfx[x+y*64] ^= 1
				}
			}
		}
		e.shouldDraw = true
		e.next()
		log.Printf("Exec opcode 0x%x\n", opcode)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x9E:
			x := opcode & 0x0F00 >> 8
			key := byte(e.V[x])
			if e.pressed(key) {
				e.skip()
			} else {
				e.next()
			}
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0xA1:
			x := opcode & 0x0F00 >> 8
			key := byte(e.V[x])
			if !e.pressed(key) {
				e.skip()
			} else {
				e.next()
			}
			log.Printf("Exec opcode 0x%x\n", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x07:
			x := opcode & 0x0F00 >> 8
			e.V[x] = e.delayTimer
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x0A:
			pressed := false
			for i, v := range e.keys {
				if v {
					x := opcode & 0x0F00 >> 8
					e.V[x] = byte(i)
					pressed = true
				}
			}
			if pressed {
				e.next()
			}
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x15:
			x := opcode & 0x0F00 >> 8
			e.delayTimer = e.V[x]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x18:
			x := opcode & 0x0F00 >> 8
			e.soundTimer = e.V[x]
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x1E:
			x := opcode & 0x0F00 >> 8
			e.I += uint16(e.V[x])
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x29:
			// 0xFX29 Sets I to the location of the sprite for the character in VX.
			// Characters 0-F (in hexadecimal) are represented by a 4x5 font
			vx := e.V[opcode&0x0F00>>8]
			e.I = uint16(vx) * 5
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x33:
			x := opcode & 0x0F00 >> 8
			e.Memory[e.I] = e.V[x] / 100
			e.Memory[e.I+1] = (e.V[x] / 10) % 10
			e.Memory[e.I+2] = e.V[x] % 10
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x55:
			x := opcode & 0x0F00 >> 8
			for i := 0; i < int(x)+1; i++ {
				e.Memory[int(e.I)+i] = e.V[i]
			}
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		case 0x65:
			x := opcode & 0x0F00 >> 8
			for i := 0; i < int(x)+1; i++ {
				e.V[i] = e.Memory[int(e.I)+i]
			}
			e.next()
			log.Printf("Exec opcode 0x%x\n", opcode)
		default:
			log.Fatalf("Unexpected opcode 0x%x\n", opcode)
		}
	default:
		log.Fatalf("Unexpected opcode 0x%x\n", opcode)
	}
}

// Print Emulator status
func (e *Emulator) Print() {
	fmt.Printf("opcode:%v\n", e.Opcode)
	fmt.Printf("memory:%v\n", e.Memory)
	fmt.Printf("v:%v\n", e.V)
	fmt.Printf("i:%v\n", e.I)
	fmt.Printf("pc:%v\n", e.Pc)
	fmt.Printf("gfx:%v\n", e.Gfx)
	fmt.Printf("delayTimer:%v\n", e.delayTimer)
	fmt.Printf("soundTimer:%v\n", e.soundTimer)
	fmt.Printf("stack:%v\n", e.Stack)
	fmt.Printf("sp:%v\n", e.Sp)
	fmt.Printf("key:%v\n", e.keys)
}

func (e *Emulator) pressed(key byte) bool {
	return e.keys[key]
}

func (e *Emulator) pressedKey() byte {
	for i, v := range e.keys {
		if v {
			return byte(i)
		}
	}
	return 0xff
}

// https://github.com/veandco/go-sdl2-examples/blob/master/examples/keyboard-input/keyboard-input.go
func (e *Emulator) Run() (err error) {

	running := true
	for running {

		opcode := e.Fetch()
		e.Exec(opcode)
		if e.delayTimer > 0 {
			e.delayTimer--
		}
		if e.soundTimer > 0 {
			e.soundTimer--
		}
		if e.shouldDraw {
			e.draw()
		}
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			fmt.Println("event loop")
			switch et := ev.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				if et.Type == sdl.KEYUP {
					fmt.Printf("keyup: %v", byte(et.Keysym.Scancode))
					if v, ok := e.keyMap[int(et.Keysym.Scancode)]; ok {
						e.keys[v] = false
					}
				} else if et.Type == sdl.KEYDOWN {
					fmt.Printf("keydown: %v", byte(et.Keysym.Scancode))
					if v, ok := e.keyMap[int(et.Keysym.Scancode)]; ok {
						e.keys[v] = true
					}
				}
			}
		}
		// Chip8 cpu clock worked at frequency of 60Hz, so set delay to (1000/60)ms
		sdl.Delay(1000 / 60)

	}

	return
}

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("no ROM file")
	}
	filepath := os.Args[1]
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	emu.InitDisplay()
	defer emu.DestroyDisplay()
	emu.Load(filepath)
	if err := emu.Run(); err != nil {
		os.Exit(1)
	}

}
