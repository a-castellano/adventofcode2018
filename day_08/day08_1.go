// Álvaro Castellano Vela 2018/12/11

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	Metadata []int
	Childs   []*Tree
}

func readLineFromFile(filename string) []int {

	var line []string

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	line = strings.Split(scanner.Text(), " ")

	puzzle := make([]int, len(line))
	for i, number := range line {
		puzzle[i], _ = strconv.Atoi(number)
	}

	return puzzle
}

func showTree(tree *Tree, depth int) {
	tabs := strings.Repeat("\t", depth)
	fmt.Printf("%s%s\n", tabs, (*tree).Metadata)
	if len((*tree).Childs) > 0 {
		depth++
		for _, child := range (*tree).Childs {
			showTree(child, depth)
		}
	}
}

func calculateLicense(tree *Tree) int {

	var license int = 0
	for _, data := range (*tree).Metadata {
		license += data
	}
	if len((*tree).Childs) > 0 {
		for _, child := range (*tree).Childs {
			license += calculateLicense(child)
		}
	}

	return license
}

func makeTree(input *[]int, startIndex int) (Tree, int) {
	var numberOfChildNodes, metadataEntries int = (*input)[startIndex], (*input)[startIndex+1]
	var metadataIndex int = startIndex + 2
	var tree Tree
	var metadata []int
	var offset int = metadataIndex

	for i := 0; i < numberOfChildNodes; i++ {

		var child Tree

		child, offset = makeTree(input, offset)
		metadataIndex += offset
		tree.Childs = append(tree.Childs, &child)
	}

	for i := 0; i < metadataEntries; i++ {
		metadata = append(metadata, (*input)[offset])
		offset++
	}
	tree.Metadata = metadata

	return tree, offset
}

func main() {

	var input []int
	var tree Tree

	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	input = readLineFromFile(filename)
	tree, _ = makeTree(&input, 0)

	fmt.Printf("License: %d\n", calculateLicense(&tree))
}
