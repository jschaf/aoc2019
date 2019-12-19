package main

import (
	"aoc2019/files"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func main() {
	lines := files.ReadAllLines("day14/input.txt")

	rs := parseReactions(lines)
	// numOre := part1(rs)
	numOre := part1a(rs)
	fmt.Printf("Part1: numOre: %d", numOre)

	bytes, _ := ioutil.ReadFile("day14/input.txt")
	fmt.Println("Part 1: " + part1c(string(bytes)))
	fmt.Println("Part 2: " + part2(string(bytes)))
}

type reaction struct {
	elem        string
	numProduced int
	inputs      map[string]int
}

func (r reaction) String() string {
	b := strings.Builder{}
	for elem, count := range r.inputs {
		b.WriteString(fmt.Sprintf("%d %s, ", count, elem))
	}
	b.WriteString(" => ")
	b.WriteString(fmt.Sprintf("%d %s", r.numProduced, r.elem))
	return b.String()
}

func parseElemCount(e string) (int, string) {
	split := strings.Split(e, " ")
	n, _ := strconv.Atoi(split[0])
	return n, split[1]
}

func parseReactions(lines []string) map[string]reaction {
	rs := make(map[string]reaction)
	for _, line := range lines {
		// 2 DLZS, 2 VCFX, 15 PDTP, 14 ZDWX, 35 NBZC, 20 JVMF, 1 BGWMS => 3 DWRH
		split := strings.Split(line, " => ")
		inputs := strings.Split(split[0], ", ")
		produced := split[1]
		r := reaction{}
		r.inputs = make(map[string]int)
		for _, input := range inputs {
			n, elem := parseElemCount(input)
			r.inputs[elem] = n
		}
		n, elem := parseElemCount(produced)
		r.elem = elem
		r.numProduced = n
		rs[elem] = r
	}
	return rs
}

func alchemize(elem string, quantity int, reactions map[string]reaction, extras map[string]int) int {
	if elem == "ORE" {
		return quantity
	}

	if extras[elem] >= quantity {
		return 0
	} else {
		quantity -= extras[elem]
		extras[elem] = 0
	}

	r := reactions[elem]
	multFactor := int(math.Ceil(float64(quantity) / float64(r.numProduced)))
	ore := 0
	for childElem, childQuantity := range r.inputs {
		ore += alchemize(childElem, childQuantity*multFactor, reactions, extras)
	}

	total := multFactor * r.numProduced
	extras[elem] += total - quantity
	return ore
}

func part1a(reactions map[string]reaction) int {
	return alchemize("FUEL", 1, reactions, make(map[string]int))
}

// COPIED FROM https://github.com/davidaayers/advent-of-code-2019/blob/992e57bdfa9848066a561a7f2a8f8621cc973afa/day14/day14.go#L28
// =========================
type ReactionRecipe struct {
	Ingredient
	inputs []Ingredient
}

type Ingredient struct {
	result string
	count  int
}

func Produce(desiredElement string, numDesired int, recipes map[string]ReactionRecipe, excess map[string]int) int {
	// we make ore for free!
	if desiredElement == "ORE" {
		return numDesired
	}

	// if we have enough excess already, consume it
	if excess[desiredElement] >= numDesired {
		excess[desiredElement] -= numDesired
		return 0
	}

	// if we don't have enough in excess, use what we have
	if excess[desiredElement] > 0 {
		numDesired -= excess[desiredElement]
		excess[desiredElement] = 0
	}

	// how many batches must we make?
	recipe := recipes[desiredElement]
	batches := int(math.Ceil(float64(numDesired) / float64(recipe.count)))

	// consume the necessary ingredients to produce this element
	ore := 0
	for _, input := range recipe.inputs {
		ore += Produce(input.result, input.count*batches, recipes, excess)
	}

	// produce, and store any excess for later use
	numProduced := batches * recipe.count
	excess[desiredElement] += numProduced - numDesired

	return ore
}
func parseReaction(reaction string) ReactionRecipe {
	parts := strings.Split(reaction, "=>")
	reactionResult, resultCnt := parseElement(parts[1])

	recipe := ReactionRecipe{
		Ingredient: Ingredient{reactionResult, resultCnt},
		inputs:     make([]Ingredient, 0),
	}

	for _, inputStr := range strings.Split(parts[0], ",") {
		input, inputCnt := parseElement(inputStr)
		recipe.inputs = append(recipe.inputs, Ingredient{
			result: input,
			count:  inputCnt,
		})
	}

	return recipe
}

func parseElement(element string) (elementName string, count int) {
	element = strings.TrimSpace(element)
	parts := strings.Split(element, " ")
	elementName = parts[1]
	count, _ = strconv.Atoi(parts[0])
	return
}

func DetermineRequiredOre(reactions []string, numFuelDesired int) int {
	recipes := make(map[string]ReactionRecipe, len(reactions))
	for _, reaction := range reactions {
		fmt.Println("reaction: " + reaction)
		recipe := parseReaction(reaction)
		recipes[recipe.result] = recipe
	}

	return Produce("FUEL", numFuelDesired, recipes, make(map[string]int))
}

func part1c(input string) string {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")
	ore := DetermineRequiredOre(lines, 1)
	return "Answer: " + strconv.Itoa(ore)

}

func DetermineMaxFuelForOre(reactions []string, ore int) int {
	start := 0
	end := ore
	guesses := 0
	lastGuess := 0
	fuelGuess := 0
	for {
		guesses++

		fuelGuess = (end-start)/2 + start

		requiredOre := DetermineRequiredOre(reactions, fuelGuess)
		if requiredOre == ore {
			break
		}

		if requiredOre > ore {
			end = fuelGuess
		} else {
			start = fuelGuess
		}

		// circuit breaker
		if guesses > 1000 || fuelGuess == lastGuess {
			break
		}

		lastGuess = fuelGuess
	}

	return fuelGuess
}

func part2(input string) string {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")
	fuel := DetermineMaxFuelForOre(lines, 1000000000000)
	return "Answer: " + strconv.Itoa(fuel)
}
