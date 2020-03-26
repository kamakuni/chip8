package main

import "testing"

func TestEmulator_Fetch(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0xA2
	data[1] = 0xF0
	emu.Load(data)
	actual := emu.Fetch()
	expected := uint16(0xA2F0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x0000(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x00
	data[1] = 0x00
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Gfx
	expected := [32][64]bool{}
	if actual != expected {
		t.Errorf("got: %v,but expected: %v", actual, expected)
	}
}

func TestEmulator_Decode0x000E(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x20
	data[1] = 0xF0
	data[2] = 0x00
	data[3] = 0x0E
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x00F0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x1NNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x10
	data[1] = 0xF0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x00F0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x2NNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x20
	data[1] = 0xF0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x00F0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x3XNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x3F
	data[1] = 0xF0
	emu.Load(data)
	emu.V[15] = 0xF0
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x200) + 2
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x4XNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x4F
	data[1] = 0xF0
	emu.Load(data)
	emu.V[15] = 0xE0
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x200) + 2
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x5XYN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x50
	data[1] = 0xF0
	emu.Load(data)
	emu.V[0] = 0x0F
	emu.V[15] = 0x0F
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(0x200) + 2
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x6XNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x6F
	data[1] = 0xF0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := int(0xF0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x7XNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x7F
	data[1] = 0xF0
	emu.Load(data)
	emu.V[15] = 1
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 1 + int(0xF0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY0(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE0
	emu.Load(data)
	emu.V[14] = 1
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 1
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY1(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE1
	emu.Load(data)
	emu.V[14] = 0xF0
	emu.V[15] = 0x0F
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 0xFF
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY2(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE2
	emu.Load(data)
	emu.V[14] = 0x0F
	emu.V[15] = 0xFF
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 0x0F
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY3(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE3
	emu.Load(data)
	emu.V[14] = 0x0F
	emu.V[15] = 0xFF
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 0xF0
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY4(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE4
	emu.Load(data)
	emu.V[14] = 0x0E
	emu.V[15] = 0x01
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 0x0F
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY5(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8F
	data[1] = 0xE5
	emu.Load(data)
	emu.V[14] = 0x01
	emu.V[15] = 0x0E
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[15])
	expected := 0x0D
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0xANNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0xA2
	data[1] = 0xF0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.I
	expected := uint16(0x02F0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}
