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

func TestEmulator_Decode0x00E0(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x00
	data[1] = 0xE0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Gfx
	expected := [32][64]bool{}
	if actual != expected {
		t.Errorf("got: %v,but expected: %v", actual, expected)
	}
}

func TestEmulator_Decode0x00EE(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 4)
	data[0] = 0x20
	data[1] = 0xF0
	data[2] = 0x00
	data[3] = 0xEE
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
	emu.V[14] = 0xF0
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
	data[0] = 0x4E
	data[1] = 0xF0
	emu.Load(data)
	emu.V[14] = 0xE0
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
	data[1] = 0xE0
	emu.Load(data)
	emu.V[0] = 0x0F
	emu.V[14] = 0x0F
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
	data[0] = 0x6E
	data[1] = 0xF0
	emu.Load(data)
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := int(0xF0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x7XNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x7E
	data[1] = 0xF0
	emu.Load(data)
	emu.V[14] = 1
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 1 + int(0xF0)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY0(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD0
	emu.Load(data)
	emu.V[13] = 1
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 1
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY1(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD1
	emu.Load(data)
	emu.V[13] = 0xF0
	emu.V[14] = 0x0F
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0xFF
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY2(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD2
	emu.Load(data)
	emu.V[13] = 0x0F
	emu.V[14] = 0xFF
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x0F
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY3(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD3
	emu.Load(data)
	emu.V[13] = 0x0F
	emu.V[14] = 0xFF
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0xF0
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY4(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD4
	emu.Load(data)
	emu.V[13] = 0x0E
	emu.V[14] = 0x01
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x0F
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY5(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD5
	emu.Load(data)
	emu.V[13] = 0x01
	emu.V[14] = 0x0E
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x0D
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY6(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD6
	emu.Load(data)
	emu.V[14] = 0x02
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x01
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XY7(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xD7
	emu.Load(data)
	emu.V[13] = 0x0E
	emu.V[14] = 0x01
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x0D
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x8XYE(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x8E
	data[1] = 0xDE
	emu.Load(data)
	emu.V[14] = 0x01
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.V[14])
	expected := 0x02
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}

func TestEmulator_Decode0x9XY0(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0x9E
	data[1] = 0xD0
	emu.Load(data)
	emu.V[0xE] = 0x01
	emu.V[0xD] = 0x02
	expected := int(emu.Pc) + 4
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := int(emu.Pc)
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

func TestEmulator_Decode0xBNNN(t *testing.T) {
	fonts := NewFonts()
	emu := NewEmulator(fonts)
	data := make([]byte, 2)
	data[0] = 0xB0
	data[1] = 0x01
	emu.Load(data)
	emu.V[0] = 1
	opcode := emu.Fetch()
	emu.Decode(opcode)
	actual := emu.Pc
	expected := uint16(2)
	if actual != expected {
		t.Errorf("got: 0x%x,but expected: 0x%x", actual, expected)
	}
}
