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

var arr = [15][20]int{}
var alive bool = true
var loopSpeed int = 200 //ms
var key string = ""
var lastKey string = ""
var eaten bool = true
var food [2]int
var tail [][2]int
var pb int = 0

func listener() {
	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if string(char) == "w" || string(char) == "a" || string(char) == "s" || string(char) == "d" {
			key = string(char)
		}

		if string(char) == "q" && !alive {
			os.Exit(0)
		}
	}
}

func show() {
	// visualisation
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

func main() {

	go listener()
	rand.Seed(time.Now().UnixNano())

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

	arr[len(arr)/2][len(arr[0])/2] = 1
	for alive {

		for i := 0; i < 10; i++ {
			fmt.Println()
		}

		var head [2]int
		var found bool = false

		for row := 0; row < len(arr); row++ {
			for col := 0; col < len(arr[0]); col++ {
				if arr[row][col] == 1 && !found {
					head[0] = row
					head[1] = col

					if key == "w" && lastKey == "s" {
						key = "s"
					} else if key == "s" && lastKey == "w" {
						key = "w"
					} else if key == "a" && lastKey == "d" {
						key = "d"
					} else if key == "d" && lastKey == "a" {
						key = "a"
					}

					switch key {
					case "w":
						if !(row-1 == -1) && arr[row-1][col] != 2 {
							arr[row][col] = 0
							arr[row-1][col] = 1
							found = true
						} else {
							alive = false
						}
					case "a":
						if !(col-1 == -1) && arr[row][col-1] != 2 {
							arr[row][col] = 0
							arr[row][col-1] = 1
							found = true
						} else {
							alive = false
						}

					case "s":
						if !(row+1 == len(arr)) && arr[row+1][col] != 2 {
							arr[row][col] = 0
							arr[row+1][col] = 1
							found = true
						} else {
							alive = false
						}

					case "d":
						if !(col+1 == len(arr[0])) && arr[row][col+1] != 2 {
							arr[row][col] = 0
							arr[row][col+1] = 1
							found = true
						} else {
							alive = false
						}
					}
				}
			}
		}

		if head[0] == food[0] && head[1] == food[1] {
			eaten = true
			var newTail [2]int

			if len(tail) == 0 {
				newTail[0] = head[0]
				newTail[1] = head[1]
			} else {
				newTail[0] = tail[len(tail)-1][0]
				newTail[1] = tail[len(tail)-1][1]
			}

			tail = append(tail, newTail)
		}

		if len(tail) >= 2 {
			for i := len(tail) - 1; i > 0; i-- {
				if i == len(tail)-1 && !eaten {
					arr[tail[i][0]][tail[i][1]] = 0
				}

				tail[i][0] = tail[i-1][0]
				tail[i][1] = tail[i-1][1]
			}
		}

		if len(tail) >= 1 {
			if len(tail) == 1 {
				arr[tail[0][0]][tail[0][1]] = 0
			}

			tail[0][0] = head[0]
			tail[0][1] = head[1]
		}

		// for row := 0; row < len(arr); row++ {
		// 	for col := 0; col < len(arr[0]); col++ {
		// 		for i := 0; i < len(tail); i++ {
		// 			if tail[i][0] == row && tail[i][1] == col {
		// 				arr[row][col] = 2
		// 			}
		// 		}
		// 	}

		for i := 0; i < len(tail); i++ {
			arr[tail[i][0]][tail[i][1]] = 2
		}

		if eaten {

			var zeros [][2]int
			var newZero [2]int

			for row := 0; row < len(arr); row++ {
				for col := 0; col < len(arr[0]); col++ {
					if arr[row][col] == 0 {
						newZero[0] = row
						newZero[1] = col
						zeros = append(zeros, newZero)
					}
				}
			}

			var r int = rand.Intn(len(zeros))
			arr[zeros[r][0]][zeros[r][1]] = 3

			food[0] = zeros[r][0]
			food[1] = zeros[r][1]

			eaten = false
		}

		lastKey = key
		show()

		pbB, err := os.ReadFile(highscore.Name())
		if err != nil {
			log.Fatal(err)
		}

		pb, err = strconv.Atoi(string(pbB))
		if err != nil {
			log.Fatal(err)
		}

		if len(tail) > pb {
			os.Truncate(highscore.Name(), 0)
			_, err = highscore.WriteString(fmt.Sprint(len(tail)))
			if err != nil {
				log.Fatal(err)
			}

			pb = len(tail)
		}

		var scoreS string = ("SCORE: " + fmt.Sprint(len(tail)))
		var pbS string = ("HIGHSCORE: " + fmt.Sprint(pb))
		fmt.Print("  ", scoreS)

		for i := 1; i < (2*len(arr[0]))-(len(scoreS)+len(pbS)); i++ {
			fmt.Print(" ")
		}

		fmt.Println(pbS)
		time.Sleep(time.Millisecond * time.Duration(loopSpeed))
	}

	for i := 0; i < 10; i++ {
		fmt.Println()
	}

	fmt.Println("You lost! Your final score is", len(tail), "point(s).")
	fmt.Println("Press Q to exit...")
	fmt.Println()

	for {
		time.Sleep(time.Second)
	}
}
