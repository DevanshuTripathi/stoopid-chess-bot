package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	gs := NewGameState()

	// Infinite loop waiting for Python's instructions
	for {
		// Read command from Python
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text = strings.TrimSpace(text)

		// 1. Handle Shutdown
		if text == "quit" {
			break
		}

		// 2. Handle Game State Sync
		// Expected format: "position moves e2e4 e7e5 g1f3..."
		if strings.HasPrefix(text, "position moves") {
			gs = NewGameState() // Start from scratch to sync move history
			parts := strings.Split(text, " ")
			if len(parts) > 2 {
				moveStrings := parts[2:]
				for _, mStr := range moveStrings {
					validMoves := gs.GetValidMoves()
					for _, vm := range validMoves {
						if moveToNotation(vm) == mStr {
							gs.MakeMove(vm)
							break
						}
					}
				}
			}
			fmt.Println("ready") // Signal back to Python that sync is done

			// 3. Handle AI Move Request
		} else if text == "go" {
			validMoves := gs.GetValidMoves()
			if len(validMoves) == 0 {
				fmt.Println("bestmove none")
			} else {
				bestMove := FindBestMove(gs, validMoves)
				fmt.Printf("bestmove %s\n", moveToNotation(bestMove))
			}
		}
	}
}

// moveToNotation converts Go's internal Move struct to a string like "e2e4" or "a7a8q"
func moveToNotation(m Move) string {
	cols := "abcdefgh"
	rows := "87654321"

	// Standard coordinate notation
	start := string(cols[m.StartCol]) + string(rows[m.StartRow])
	end := string(cols[m.EndCol]) + string(rows[m.EndRow])

	// Handle pawn promotion suffix
	promotion := ""
	if m.IsPawnPromotion {
		promotion = "q"
	}

	return start + end + promotion
}
