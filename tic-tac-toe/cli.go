package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Represents an arbitrarily tic-tac-toe board.
type TicTacToeBoard struct {
	winLineLength int      // Minimum number of consecutive cells a player must have in a straight line to win. This includes horizontal, vertical, and diagonal lines.
	board         [][]rune // Board of cells for tic-tac-toe. ' ' = empty cell. Different runes represent different players. Indexed by board[y][x].
	freeSpots     int      // Number of free spots for players to put marks in. Game ends in a tie if all the spots are taken an no player has won.
	winner        rune     // Winner of the game, ' ' if there are no winners
	tie           bool     // Whether there's a tie or not.
}

// Creates a new TicTacToeBoard with a win line length, a board width, and a board height.
func NewTicTacToeBoard(winLineLength int, boardWidth int, boardHeight int) TicTacToeBoard {
	if winLineLength == 0 {
		panic("winLineLength must be > 0")
	}
	if boardWidth < 1 || boardHeight < 1 {
		panic("boardWidth must be >= 1 and boardHeight must be >= 1")
	}
	board := make([][]rune, boardHeight)
	for r := 0; r < boardHeight; r++ {
		board[r] = make([]rune, boardWidth)
	}
	return TicTacToeBoard{
		winLineLength,
		board,
		boardWidth * boardHeight,
		0,
		false,
	}
}

// Prints the board.
func (tb *TicTacToeBoard) Print() {
	fmt.Println()
	fmt.Println(tb.String())
	fmt.Println()
}

// Returns a string representation of the board.
func (tb *TicTacToeBoard) String() string {
	str := ""
	str += "     " + strings.Repeat("   ", len(tb.board[0])/2) + "board\n"
	for y := len(tb.board) - 1; y >= 0; y-- {
		row := "   "
		if y == len(tb.board)/2 {
			row = "  y"
		}
		row += fmt.Sprintf("%2v | ", y)
		for x := 0; x < len(tb.board[0]); x++ {
			cell_str := '.'
			if tb.board[y][x] != 0 {
				cell_str = tb.board[y][x]
			}
			row += fmt.Sprintf("%-2c ", cell_str)
		}
		str += row + "\n"
	}
	str += "      +-" + strings.Repeat("---", len(tb.board[0])) + "\n"
	str += "        "
	for x := 0; x < len(tb.board[0]); x++ {
		str += fmt.Sprintf("%-2v ", x)
	}
	str += "\n"
	str += "        " + strings.Repeat("   ", len(tb.board[0])/2) + "x"
	return str
}

// Checks if x, y are within the bounds of the board.
func (tb *TicTacToeBoard) InBounds(x int, y int) bool {
	return x >= 0 && x < len(tb.board[0]) && y >= 0 && y < len(tb.board)
}

// Returns the mark at point x, y.
func (tb *TicTacToeBoard) GetMark(x int, y int) rune {
	if !tb.InBounds(x, y) {
		return 0
	}
	return tb.board[y][x]
}

// Returns whether the point x, y is empty. A point x, y
// is empty if no player has placed a mark there.
func (tb *TicTacToeBoard) IsSpotEmpty(x int, y int) bool {
	return tb.GetMark(x, y) == 0
}

// Returns if the game is over.
// A game is over if there is a tie or if there is a specific winner.
func (tb *TicTacToeBoard) IsGameOver() bool {
	return tb.tie || tb.winner != 0
}

// Places a player's mark down on the board, and then updates the state of the game.
// If the game is won, then tb.winner is set to the winner's mark.
// If the game is a tie, then tb.tie is set to true.
func (tb *TicTacToeBoard) PlaceMove(x int, y int, player rune) {
	if tb.IsGameOver() || !tb.InBounds(x, y) || tb.board[y][x] != 0 || tb.freeSpots == 0 {
		return
	}
	tb.board[y][x] = player
	tb.freeSpots--

	getLineLength := func(x int, y int, x_off int, y_off int) int {
		length := 0
		x += x_off
		y += y_off
		for tb.InBounds(x, y) && tb.board[y][x] == player {
			length += 1
			x += x_off
			y += y_off
		}
		return length
	}

	DirOffsets := [...]struct {
		x int
		y int
	}{{0, 1}, {1, 1}, {1, 0}, {1, -1}}
	for _, dir := range DirOffsets {
		// find length by adding current point, the points in dir, and the points opposite of dir
		// |<------------*------------>|
		// -dir      curr_point       dir
		dirLength := getLineLength(x, y, dir.x, dir.y)
		oppDirLength := getLineLength(x, y, -dir.x, -dir.y)
		length := 1 + dirLength + oppDirLength
		if length >= tb.winLineLength {
			tb.winner = player
			break
		}
	}

	if tb.winner == 0 && tb.freeSpots == 0 {
		// tie if no winner exists and all spots are taken
		tb.tie = true
	}
}

func tryArgInt(i int, defaultValue int) int {
	if i < len(os.Args) {
		val, err := strconv.Atoi(os.Args[i])
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func tryArg(i int, defaultValue string) string {
	if i < len(os.Args) {
		return os.Args[i]
	}
	return defaultValue
}

// CLI for tic-tac-toe
// tic-tac-toe [board height] [board width] [win line length] [players]
func main() {
	boardHeight := max(tryArgInt(1, 3), 1)
	boardWidth := max(tryArgInt(2, 3), 1)
	winLineLength := max(tryArgInt(3, 3), 1)
	players := tryArg(4, "XO")
	if len(players) < 2 {
		players = "XO"
	}
	scanner := bufio.NewScanner(os.Stdin)

	tb := NewTicTacToeBoard(winLineLength, boardWidth, boardHeight)
	currPlayerIndex := 0
	fmt.Printf("\nTic-tac-toe\n  %v x %v board\n  %v marks in row to win\n  %v players = %v\n", boardHeight, boardWidth, winLineLength, len(players), players)
	for {
		tb.Print()
		fmt.Printf("'%c' turn. Enter your move as 'x y':\n", players[currPlayerIndex])
		if !scanner.Scan() {
			break
		}
		var x, y int
		n, err := fmt.Sscanf(scanner.Text(), "%d %d\n", &x, &y)
		if n != 2 || err != nil || x < 0 || x >= boardWidth || y < 0 || y >= boardHeight {
			fmt.Println("Invalid input. Please input two space separated integer coordinates 'x y'.")
			continue
		}
		if !tb.IsSpotEmpty(x, y) {
			fmt.Printf("Spot (%v %v) is taken, please choose another spot.\n", x, y)
			continue
		}
		tb.PlaceMove(x, y, rune(players[currPlayerIndex]))
		fmt.Printf("Placed '%c' at (%v, %v)\n", players[currPlayerIndex], x, y)

		if tb.IsGameOver() {
			tb.Print()
			if tb.winner != 0 {
				fmt.Printf("'%c' wins!\n", tb.winner)
				break
			} else {
				fmt.Println("Game tied!")
				break
			}
		}
		currPlayerIndex = (currPlayerIndex + 1) % len(players)
	}
}
