package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

var key string = ""
var lastKey string = ""
var alive bool = true
var pb int = 0

var arr = [13][13]int{}

func readInput() {
	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if string(char) == "w" || string(char) == "a" || string(char) == "s" || string(char) == "d" {
			key = string(char)
		} else if string(char) == "q" && !alive {
			os.Exit(0)
		}
	}
}

func show() {

	// visualisation

	for i := 0; i < 10; i++ {
		fmt.Println()
	}

	fmt.Print(" _")
	for i := 0; i < len(arr[0])-1; i++ {
		fmt.Print("_ ")
	}
	fmt.Println("__  ")

	for row := 0; row < len(arr); row++ {
		fmt.Print("| ")

		for col := 0; col < len(arr[0]); col++ {
			switch arr[row][col] {
			case 0:
				fmt.Print("  ")
			case 1:
				fmt.Print("0 ")
			case 2:
				fmt.Print("o ")
			case 3:
				fmt.Print(". ")
			}

		}

		fmt.Print("| ")
		fmt.Println()
	}

	fmt.Print(" ¯")
	for i := 0; i < len(arr[0])-1; i++ {
		fmt.Print("¯ ")
	}
	fmt.Println("¯¯  ")
}

func locateHead() [2]int {
	for row := 0; row < len(arr); row++ {
		for col := 0; col < len(arr[row]); col++ {
			if arr[row][col] == 1 {
				return [2]int{row, col}
			}
		}
	}
	return [2]int{0, 0}
}

func moveHead(head [2]int) {

	switch key {
	case "w":
		if lastKey == "s" {
			key = "s"
		}
	case "a":
		if lastKey == "d" {
			key = "d"
		}
	case "s":
		if lastKey == "w" {
			key = "w"
		}
	case "d":
		if lastKey == "a" {
			key = "a"
		}
	}

	lastKey = key

	for row := 0; row < len(arr); row++ {
		for col := 0; col < len(arr[row]); col++ {
			if arr[row][col] == 1 {
				head[0] = row
				head[1] = col
			}
		}
	}

	switch key {
	case "w":
		if head[0]-1 < 0 || arr[head[0]-1][head[1]] == 2 {
			die()
		} else {

			arr[head[0]][head[1]] = 0
			arr[head[0]-1][head[1]] = 1
		}
	case "a":
		if head[1]-1 < 0 || arr[head[0]][head[1]-1] == 2 {
			die()
		} else {

			arr[head[0]][head[1]] = 0
			arr[head[0]][head[1]-1] = 1
		}
	case "s":
		if head[0]+1 >= len(arr) || arr[head[0]+1][head[1]] == 2 {
			die()
		} else {

			arr[head[0]][head[1]] = 0
			arr[head[0]+1][head[1]] = 1
		}
	case "d":
		if head[1]+1 >= len(arr[head[0]]) || arr[head[0]][head[1]+1] == 2 {
			die()
		} else {

			arr[head[0]][head[1]] = 0
			arr[head[0]][head[1]+1] = 1
		}
	}
}

func spawnFood() [2]int {
	var free = [][2]int{}

	for row := 0; row < len(arr); row++ {
		for col := 0; col < len(arr[row]); col++ {
			if arr[row][col] == 0 {
				free = append(free, [2]int{row, col})
			}
		}
	}

	var r int = rand.Intn(len(free) - 1)

	if len(free) > 0 {
		arr[free[r][0]][free[r][1]] = 3
		return [2]int{free[r][0], free[r][1]}
	} else {
		die()
		return [2]int{0, 0}
	}
}

func eatFood(head [2]int, food [2]int) bool {
	if head[0] == food[0] && head[1] == food[1] {
		return true
	} else {
		return false
	}
}

func spawnTail(tail [][2]int, head [2]int) [][2]int {

	if len(tail) >= 1 {
		tail = append(tail, [2]int{tail[len(tail)-1][len(tail[0])-1]})
	} else {
		tail = append(tail, [2]int{head[0], head[1]})
	}

	if len(tail) >= 2 {
		for i := len(tail) - 1; i > 0; i-- {
			tail[i][0] = tail[i-1][0]
			tail[i][1] = tail[i-1][1]
		}
	}

	tail[0][0] = head[0]
	tail[0][1] = head[1]

	return tail
}

func moveTail(tail [][2]int, head [2]int) [][2]int {
	if len(tail) > 1 {
		arr[tail[len(tail)-1][0]][tail[len(tail)-1][1]] = 0
		for i := len(tail) - 1; i > 0; i-- {
			tail[i][0] = tail[i-1][0]
			tail[i][1] = tail[i-1][1]
		}
	} else if len(tail) == 1 {
		arr[tail[0][0]][tail[0][1]] = 0
	}

	if len(tail) > 0 {
		tail[0][0] = head[0]
		tail[0][1] = head[1]
	}

	return tail
}

func printTail(tail [][2]int) {
	for i := 0; i < len(tail); i++ {
		arr[tail[i][0]][tail[i][1]] = 2
	}
}

func die() {
	alive = false

	for i := 0; i < 9; i++ {
		fmt.Println()
	}

	fmt.Println("You lost! Your final score is", pb, "point(s).")
	fmt.Println("Press Q to exit.")

	for key != "q" {
		time.Sleep(time.Millisecond)
	}

	os.Exit(0)
}

func scoreCounter(tail [][2]int, pb int) {
	var scoreS string = ("SCORE: " + fmt.Sprint(len(tail)))
	var pbS string = ("HIGHSCORE: " + fmt.Sprint(pb))
	fmt.Print("  ", scoreS)

	for i := 1; i < (2*len(arr[0]))-(len(scoreS)+len(pbS)); i++ {
		fmt.Print(" ")
	}

	fmt.Println(pbS)
}

func readHighscore(highscore os.File) int {
	pbB, err := os.ReadFile(highscore.Name())
	if err != nil {
		log.Fatal(err)
	}

	pb, err := strconv.Atoi(string(pbB))
	if err != nil {
		log.Fatal(err)
	}

	return pb
}

func writeHighscore(highscore os.File, tail [][2]int) {
	os.Truncate(highscore.Name(), 0)
	_, err := highscore.WriteString(fmt.Sprint(len(tail)))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go readInput()

	arr[len(arr)/2][len(arr[0])/2] = 1

	var head [2]int = locateHead()
	var tail = [][2]int{}
	var eat bool = false
	var food [2]int = spawnFood()

	if _, err := os.Stat("./.highscore"); err != nil {
		f, err := os.Create("./.highscore")
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString("0")
		if err != nil {
			log.Fatal(err)
		}

		f.Close()
	}

	highscore, err := os.OpenFile("./.highscore", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer highscore.Close()

	pb = readHighscore(*highscore)

	for alive {
		head = locateHead()
		moveHead(head)

		eat = eatFood(head, food)
		if eat {
			food = spawnFood()
			tail = spawnTail(tail, head)
		} else {
			tail = moveTail(tail, head)
		}

		printTail(tail)
		show()

		if len(tail) > pb {
			pb = len(tail)
			writeHighscore(*highscore, tail)
		}
		scoreCounter(tail, pb)

		time.Sleep(time.Millisecond * 250)
	}
}
