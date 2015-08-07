package main

import (
	"flag"
	"fmt"
)

type connection struct {
	start, end int
	rel        string
}

type jump struct {
	start, middle, end int
}

type moveNode struct {
	filled   []bool
	parent   *moveNode
	children []*moveNode
}

func main() {
	var boardSize = flag.Int("size", 5, "Number of rows in the board triangle")
	var target = flag.Int("target", 1, "Desired number of pegs remaining")
	var blankHole = flag.Int("start", 0, "Position initially left empty")
	flag.Parse()

	fmt.Printf("Holes: %d, connections: %d, jumps: %d\n", rightVal(*boardSize), numConnections(*boardSize), numJumps(*boardSize))
	cons := generateConnections(*boardSize)
	jumps := generateJumps(cons, *boardSize)
	//fmt.Println(jumps)
	beginTree(*blankHole, *target, *boardSize, jumps)
}

func beginTree(emptyHole, target, boardSize int, jumps *[]jump) {
	startFilled := make([]bool, rightVal(boardSize)+1)
	for i := range startFilled {
		startFilled[i] = true
	}
	startFilled[emptyHole] = false
	findPossibleMoves(&moveNode{filled: startFilled, parent: nil, children: make([]*moveNode, 0)}, jumps, target, boardSize)
	//fmt.Println(startFilled)
	//    root := moveNode{filled: startFilled,
}

func findPossibleMoves(node *moveNode, jumps *[]jump, target, boardSize int) bool {
	for _, jumpCand := range *jumps {
		if node.filled[jumpCand.start] && node.filled[jumpCand.middle] && !node.filled[jumpCand.end] {
			newFilled := make([]bool, rightVal(boardSize)+1)
			copy(newFilled, node.filled)
			newFilled[jumpCand.start] = false
			newFilled[jumpCand.middle] = false
			newFilled[jumpCand.end] = true
			node.children = append(node.children, &moveNode{filled: newFilled, parent: node, children: make([]*moveNode, 0)})
		}
	}
	if len(node.children) == 0 {
		if checkConditions(node, target) {
			return true
		} else {
			//fmt.Println("Failure: ")
		}
	}
	for _, childNode := range node.children {
		if findPossibleMoves(childNode, jumps, target, boardSize) {
			return true
		}
	}
	return false
}

func checkConditions(node *moveNode, target int) bool {
	sum := 0
	for _, val := range node.filled {
		if val {
			sum += 1
		}
	}
	if target == sum {
		fmt.Println("Success!")
		printBoard(node.filled)
		for node.parent != nil {
			node = node.parent
			printBoard(node.filled)
		}
		return true
	}
	return false
}

func printBoard(arr []bool) {
	for i := 1; i > 0; i++ {
		for j := 0; j < i; j++ {
			if leftVal(i)+j == len(arr) {
				fmt.Print("\n\n\n")
				return
			}
			fmt.Print(arr[leftVal(i)+j], ", ")
		}
		fmt.Print("\n")
	}
}

func rightVal(row int) int {
	return (((row+1)*(row+1) - (row + 1)) / 2) - 1
}

func leftVal(row int) int {
	return rightVal(row-1) + 1
}

func numConnections(boardSize int) int {
	return 3 * (rightVal(boardSize) - boardSize)
}

func numJumps(boardSize int) int {
	return 6 * (rightVal(boardSize) - 3 - 2*(boardSize-2))
}

func findMatchingConnection(con connection, conMap *map[int][]connection) int {
	target := 0
	for _, jumpCand := range (*conMap)[con.end] {
		if jumpCand.rel == con.rel {
			target = jumpCand.end
		}
	}
	return target
}

func generateJumps(cons *map[int][]connection, boardSize int) *[]jump {
	jumps := make([]jump, 0, numJumps(boardSize))
	for pos := range *cons {
		for _, directCon := range (*cons)[pos] {
			cand := findMatchingConnection(directCon, cons)
			if cand > 0 {
				jumps = append(jumps, jump{start: pos, middle: directCon.end, end: cand})
				jumps = append(jumps, jump{start: cand, middle: directCon.end, end: pos})
			}
		}
	}
	return &jumps
}

func generateConnections(boardSize int) *map[int][]connection {
	cons := make(map[int][]connection)
	for row := 1; row <= boardSize; row += 1 {
		rowFirst := leftVal(row)
		for pos := rowFirst; pos <= rightVal(row); pos += 1 {
			if row != boardSize {
				rowPos := pos - rowFirst
				cons[pos] = append(cons[pos], connection{pos, leftVal(row+1) + rowPos, "cl"})
				cons[pos] = append(cons[pos], connection{pos, leftVal(row+1) + rowPos + 1, "cr"})
			}
			if pos != rightVal(row) {
				cons[pos] = append(cons[pos], connection{pos, pos + 1, "sr"})
			}
		}
	}
	return &cons
}
