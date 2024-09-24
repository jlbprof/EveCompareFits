/*

MIT License

Copyright (c) 2024 jlbprof

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// The cargo item structure is designed to parse the following cargo items or drones that often appear in a fit file
// Example: Phased Plasma M x1600

type cargoItem struct {
	name   string
	amount int
}

// This is the representation of the fit after parsing.  Note I am using maps of strings where the string is the name of the
// item in that slot, and the integer is the number of those items
// This makes it trivial to compare 2 fits, I go through the keys of each slow type, on each ship and can tell what is missing
// or if there is a different number of those items between the ships

type ship struct {
	shipName string
	shipType string

	lowSlots  map[string]int
	midSlots  map[string]int
	highSlots map[string]int
	rigs      map[string]int

	drones map[string]int
	cargo  map[string]int
}

var version = "0.1.0"

// just encapsulate the logic to determine if a key is in a map
func doesKeyExist(mymap map[string]int, key string) bool {
	_, ok := mymap[key]
	if ok {
		return true
	}
	return false
}

// add `num` of this `name` to the map
func updateMap(mymap map[string]int, name string, num int) {
	if doesKeyExist(mymap, name) {
		mymap[name] = mymap[name] + num
	} else {
		mymap[name] = num
	}
}

// This parses a cargo item line like:
// Phased Plasma M x1600
func parseCargoItem(line string) cargoItem {
	var c cargoItem
	var xi int

	// I use a regular expression to determine the name and the number
	re := regexp.MustCompile(`^(.+) x(\d+)$`)
	res := re.FindAllStringSubmatch(line, -1)

	// if the regular expression fails, it is exactly 1 of that item
	if len(res) != 1 {
		c.name = line
		c.amount = 1
	} else {
		c.name = res[0][1]
		xi, _ = strconv.Atoi(res[0][2])
		c.amount = xi
	}

	return c
}

// Read a line from the file
func myReadLine(s *bufio.Scanner) (string, bool) {
	b := s.Scan()
	if !b {
		fmt.Println("ERROR myReadLine", b)
		return "", b
	}

	return s.Text(), b
}

// Print out all the items from this map (usually a slot type)
func printSlots(items map[string]int, label string) {
	fmt.Println("")
	fmt.Println(label)

	for k, v := range items {
		fmt.Printf("%8dx %s\n", v, k)
	}
}

// print the entire fit
func printShip(thisShip ship) {
	fmt.Println("Ship:", thisShip.shipName, "Type:", thisShip.shipType)

	printSlots(thisShip.lowSlots, "Low Slots")
	printSlots(thisShip.midSlots, "Mid Slots")
	printSlots(thisShip.highSlots, "High Slots")
	printSlots(thisShip.rigs, "Rigs")
	printSlots(thisShip.drones, "Drones")
	printSlots(thisShip.cargo, "Cargo")

}

// The first line of the fit is the most complex to parse, I am using Go's implementation
// of regular expressions to parse out the complex parts
func parseShipTypeAndName(line string) (string, string) {
	re := regexp.MustCompile(`\[([^,]+), ([^\]]+)\]`)
	res := re.FindAllStringSubmatch(line, -1)

	if len(res) != 1 {
		fmt.Println("Could not parse ShipType and ShipName", len(res))
		os.Exit(1)
	}

	shipType := res[0][1]
	shipName := res[0][2]

	return shipType, shipName
}

// Read in the fit, parsing all the items in proper order and hopefully correctly
func readInFit(fName string) ship {
	file, err := os.Open(fName)
	if err != nil {
		fmt.Println("Failed to open", fName)
		os.Exit(1)
	}

	defer file.Close()

	var thisShip ship
	var line string
	var b bool

	thisShip.lowSlots = make(map[string]int)
	thisShip.midSlots = make(map[string]int)
	thisShip.highSlots = make(map[string]int)
	thisShip.rigs = make(map[string]int)
	thisShip.drones = make(map[string]int)
	thisShip.cargo = make(map[string]int)

	scanner := bufio.NewScanner(file)
	if scanner == nil {
		fmt.Println("SCANNER IS NIL")
		os.Exit(1)
	}

	line, b = myReadLine(scanner)
	if !b {
		fmt.Println("LINE", line)
		fmt.Println("Error Reading Ship Type and Name")
		os.Exit(1)
	}

	thisShip.shipType, thisShip.shipName = parseShipTypeAndName(line)

	// immediately next we get the low slots
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Low Slots")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		updateMap(thisShip.lowSlots, line, 1)
	}

	// next we get the mid slots
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Mid Slots")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		updateMap(thisShip.midSlots, line, 1)
	}

	// next we get the high slots
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading High Slots")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		updateMap(thisShip.highSlots, line, 1)
	}

	// next we get the rigs
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Rigs")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		updateMap(thisShip.rigs, line, 1)
	}

	// read 2 blank lines
	for _ = range 2 {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Drones")
			os.Exit(1)
		}
		if len(line) != 0 {
			fmt.Println("Error Reading Drones")
			os.Exit(1)
		}
	}

	// next we get the drones
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Drones")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		c := parseCargoItem(line)
		updateMap(thisShip.drones, c.name, c.amount)
	}

	// next we get the cargo
	for {
		line, b = myReadLine(scanner)
		if !b {
			fmt.Println("Error Reading Cargo")
			os.Exit(1)
		}
		if len(line) == 0 {
			break
		}

		c := parseCargoItem(line)
		updateMap(thisShip.cargo, c.name, c.amount)
	}

	return thisShip
}

// What is missing in the 2nd ship, that is on the first ship, list them so they can be removed
func compareSlotsRemoval(ship1 map[string]int, ship2 map[string]int, label string) {
	bFlag := false
	for k, v := range ship1 {
		if !doesKeyExist(ship2, k) {
			if !bFlag {
				fmt.Printf("Remove from %s\n", label)
				bFlag = true
			}

			fmt.Printf("%10dx %s\n", v, k)
		} else {
			x := ship1[k] - ship2[k]
			if x > 0 {
				if !bFlag {
					fmt.Printf("Remove from %s\n", label)
					bFlag = true
				}

				fmt.Printf("%10dx %s\n", x, k)
			}
		}
	}
}

// What is on the 2nd ship, that is not on the first ship, so that we can get a buy list
func compareSlotsAdditions(ship1 map[string]int, ship2 map[string]int, label string) {
	bFlag := false
	for k, v := range ship2 {
		if !doesKeyExist(ship1, k) {
			if !bFlag {
				fmt.Printf("Add To %s\n", label)
				bFlag = true
			}

			fmt.Printf("%10dx %s\n", v, k)
		} else {
			x := ship2[k] - ship1[k]
			if x > 0 {
				if !bFlag {
					fmt.Printf("Add To %s\n", label)
					bFlag = true
				}

				fmt.Printf("%10dx %s\n", x, k)
			}
		}
	}
}

// compare all the items and list removals and adds
func compareTwoShips(ship1 ship, ship2 ship) {
	// Ship types must be the same for compare

	if ship1.shipType != ship2.shipType {
		fmt.Println("Ships are not of the same time, no comparison given.")
		return
	}

	// Process Removals First

	compareSlotsRemoval(ship1.lowSlots, ship2.lowSlots, "Low Slots")
	compareSlotsRemoval(ship1.midSlots, ship2.midSlots, "Mid Slots")
	compareSlotsRemoval(ship1.highSlots, ship2.highSlots, "High Slots")
	compareSlotsRemoval(ship1.rigs, ship2.rigs, "Rigs")
	compareSlotsRemoval(ship1.drones, ship2.drones, "Drones")
	compareSlotsRemoval(ship1.cargo, ship2.cargo, "Cargo")

	compareSlotsAdditions(ship1.lowSlots, ship2.lowSlots, "Low Slots")
	compareSlotsAdditions(ship1.midSlots, ship2.midSlots, "Mid Slots")
	compareSlotsAdditions(ship1.highSlots, ship2.highSlots, "High Slots")
	compareSlotsAdditions(ship1.rigs, ship2.rigs, "Rigs")
	compareSlotsAdditions(ship1.drones, ship2.drones, "Drones")
	compareSlotsAdditions(ship1.cargo, ship2.cargo, "Cargo")

}

func main() {
	bVersion := false   // Are they requesting only the version?
	bJustParse := false // should we only parse

	for _, arg := range os.Args {
		if strings.ToLower(arg) == "--version" {
			bVersion = true
		}

		if strings.ToLower(arg) == "--justparse" {
			bJustParse = true
		}
	}

	if bVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Identify the program
	fmt.Printf("Eve Compare Fits: Version %s\n\n", version)

	if bJustParse {
		for i, arg := range os.Args {
			if i == 0 {
				continue
			}
			if strings.ToLower(arg) == "--justparse" {
				continue
			}
			ship := readInFit(arg)
			fmt.Println("Ship File", arg)
			printShip(ship)
			fmt.Println("")
		}
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide 2 Fit files for comparison.")
		os.Exit(1)
	}

	firstFileName := os.Args[1]
	secondFileName := os.Args[2]

	ship1 := readInFit(firstFileName)
	ship2 := readInFit(secondFileName)

	compareTwoShips(ship1, ship2)
}
