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
