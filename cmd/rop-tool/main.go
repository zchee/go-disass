package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bnagy/gapstone"
)

func main() {

	disassembler := NewDisassembler()
	disassembler.Open(os.Args[1])
	defer disassembler.f.Close()

	err := disassembler.StartEngineX86_64()
	if err != nil {
		log.Fatal(err)
	}
	defer disassembler.e.Close()

	textSection := disassembler.f.Section(".text")
	if textSection == nil {
		log.Fatal("No text section")
	}

	textData, err := textSection.Data()
	if err != nil {
		log.Fatal(err)
	}

	// Collect the location of every c3 in .text
	c3Locations := []int{}
	for i := 0; i < len(textData); i++ {
		if textData[i] == byte(195) {
			c3Locations = append(c3Locations, i)
		}
	}

	for _, location := range c3Locations {
		for instructionLength := 1; instructionLength < 9; instructionLength++ {
			instructions, err := disassembler.e.Disasm(textData[location-instructionLength:location], 0, 0x0)
			if err == nil {

				instructionAddress := textSection.Addr + uint64(location-instructionLength)

				fmt.Printf("\n\nINSTRUCTIONS FOUND AT LOCATION 0x%x, SIZE %d\n", instructionAddress, instructionLength)
				fmt.Println("________________________________________________")
				for _, ins := range instructions {
					printInstruction(instructionAddress, ins)
				}
				fmt.Println("________________________________________________")
			}
		}
	}
}

func printInstruction(startingAddress uint64, ins gapstone.Instruction) {
	fmt.Printf("0x%x:\t%s\t\t%s\n", startingAddress+uint64(ins.Address), ins.Mnemonic, ins.OpStr)
}
