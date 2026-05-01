package main

// We use bytes for the board.
// Uppercase = White ('P', 'R', 'N', 'B', 'Q', 'K')
// Lowercase = Black ('p', 'r', 'n', 'b', 'q', 'k')
// '-' = Empty Square

type GameState struct {
	Board               [8][8]byte
	WhiteToMove         bool
	MoveLog             []Move
	WhiteKingLoc        [2]int
	BlackKingLoc        [2]int
	Checkmate           bool
	Stalemate           bool
	EnPassantPossible   [2]int
	EnPassantLog        [][2]int
	CurrentCastleRights CastleRights
	CastleRightsLog     []CastleRights
}

type Move struct {
	StartRow, StartCol int
	EndRow, EndCol     int
	PieceMoved         byte
	PieceCaptured      byte
	IsPawnPromotion    bool
	IsEnPassantMove    bool
	IsCastleMove       bool
	Score              int
}

type CastleRights struct {
	Wks, Bks, Wqs, Bqs bool
}

// NewGameState initializes a fresh chess board
func NewGameState() *GameState {
	gs := &GameState{
		Board: [8][8]byte{
			{'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r'},
			{'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p'},
			{'-', '-', '-', '-', '-', '-', '-', '-'},
			{'-', '-', '-', '-', '-', '-', '-', '-'},
			{'-', '-', '-', '-', '-', '-', '-', '-'},
			{'-', '-', '-', '-', '-', '-', '-', '-'},
			{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'},
			{'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R'},
		},
		WhiteToMove:       true,
		WhiteKingLoc:      [2]int{7, 4},
		BlackKingLoc:      [2]int{0, 4},
		EnPassantPossible: [2]int{-1, -1}, // -1, -1 means no en passant possible
	}
	gs.CurrentCastleRights = CastleRights{Wks: true, Bks: true, Wqs: true, Bqs: true}
	gs.CastleRightsLog = append(gs.CastleRightsLog, gs.CurrentCastleRights)
	gs.EnPassantLog = append(gs.EnPassantLog, gs.EnPassantPossible)
	return gs
}

// MakeMove executes a move and handles all special rules
func (gs *GameState) MakeMove(move Move) {
	gs.Board[move.StartRow][move.StartCol] = '-'
	gs.Board[move.EndRow][move.EndCol] = move.PieceMoved
	gs.MoveLog = append(gs.MoveLog, move)
	gs.WhiteToMove = !gs.WhiteToMove

	// 1. Update King Location
	if move.PieceMoved == 'K' {
		gs.WhiteKingLoc = [2]int{move.EndRow, move.EndCol}
	} else if move.PieceMoved == 'k' {
		gs.BlackKingLoc = [2]int{move.EndRow, move.EndCol}
	}

	// 2. Pawn Promotion (Auto Queen for now)
	if move.IsPawnPromotion {
		if move.PieceMoved == 'P' {
			gs.Board[move.EndRow][move.EndCol] = 'Q'
		} else {
			gs.Board[move.EndRow][move.EndCol] = 'q'
		}
	}

	// 3. En Passant Capture
	if move.IsEnPassantMove {
		gs.Board[move.StartRow][move.EndCol] = '-'
	}

	// Update En Passant Possible square
	if (move.PieceMoved == 'P' || move.PieceMoved == 'p') && abs(move.StartRow-move.EndRow) == 2 {
		gs.EnPassantPossible = [2]int{(move.StartRow + move.EndRow) / 2, move.StartCol}
	} else {
		gs.EnPassantPossible = [2]int{-1, -1}
	}
	gs.EnPassantLog = append(gs.EnPassantLog, gs.EnPassantPossible)

	// 4. Castling Rights & Moves
	gs.updateCastleRights(move)
	gs.CastleRightsLog = append(gs.CastleRightsLog, gs.CurrentCastleRights)

	if move.IsCastleMove {
		if move.EndCol-move.StartCol == 2 { // Kingside
			gs.Board[move.EndRow][move.EndCol-1] = gs.Board[move.EndRow][move.EndCol+1]
			gs.Board[move.EndRow][move.EndCol+1] = '-'
		} else { // Queenside
			gs.Board[move.EndRow][move.EndCol+1] = gs.Board[move.EndRow][move.EndCol-2]
			gs.Board[move.EndRow][move.EndCol-2] = '-'
		}
	}
}

// UndoMove reverses the last move made
func (gs *GameState) UndoMove() {
	if len(gs.MoveLog) == 0 {
		return
	}

	// Pop the last move
	move := gs.MoveLog[len(gs.MoveLog)-1]
	gs.MoveLog = gs.MoveLog[:len(gs.MoveLog)-1]

	gs.Board[move.StartRow][move.StartCol] = move.PieceMoved
	gs.Board[move.EndRow][move.EndCol] = move.PieceCaptured
	gs.WhiteToMove = !gs.WhiteToMove

	// Undo King Location
	if move.PieceMoved == 'K' {
		gs.WhiteKingLoc = [2]int{move.StartRow, move.StartCol}
	} else if move.PieceMoved == 'k' {
		gs.BlackKingLoc = [2]int{move.StartRow, move.StartCol}
	}

	// Undo En Passant
	if move.IsEnPassantMove {
		gs.Board[move.EndRow][move.EndCol] = '-'
		gs.Board[move.StartRow][move.EndCol] = move.PieceCaptured
	}

	// Undo En Passant Log
	gs.EnPassantLog = gs.EnPassantLog[:len(gs.EnPassantLog)-1]
	gs.EnPassantPossible = gs.EnPassantLog[len(gs.EnPassantLog)-1]

	// Undo Castle Rights Log
	gs.CastleRightsLog = gs.CastleRightsLog[:len(gs.CastleRightsLog)-1]
	gs.CurrentCastleRights = gs.CastleRightsLog[len(gs.CastleRightsLog)-1]

	// Undo Castle Move
	if move.IsCastleMove {
		if move.EndCol-move.StartCol == 2 { // Kingside
			gs.Board[move.EndRow][move.EndCol+1] = gs.Board[move.EndRow][move.EndCol-1]
			gs.Board[move.EndRow][move.EndCol-1] = '-'
		} else { // Queenside
			gs.Board[move.EndRow][move.EndCol-2] = gs.Board[move.EndRow][move.EndCol+1]
			gs.Board[move.EndRow][move.EndCol+1] = '-'
		}
	}
}

// Helper to update castling rights based on piece movement
func (gs *GameState) updateCastleRights(move Move) {
	if move.PieceMoved == 'K' {
		gs.CurrentCastleRights.Wks = false
		gs.CurrentCastleRights.Wqs = false
	} else if move.PieceMoved == 'k' {
		gs.CurrentCastleRights.Bks = false
		gs.CurrentCastleRights.Bqs = false
	} else if move.PieceMoved == 'R' {
		if move.StartRow == 7 {
			if move.StartCol == 0 {
				gs.CurrentCastleRights.Wqs = false
			} else if move.StartCol == 7 {
				gs.CurrentCastleRights.Wks = false
			}
		}
	} else if move.PieceMoved == 'r' {
		if move.StartRow == 0 {
			if move.StartCol == 0 {
				gs.CurrentCastleRights.Bqs = false
			} else if move.StartCol == 7 {
				gs.CurrentCastleRights.Bks = false
			}
		}
	}
}

// Simple integer absolute value helper
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isWhite(p byte) bool { return p >= 'A' && p <= 'Z' }
func isBlack(p byte) bool { return p >= 'a' && p <= 'z' }
func isEnemy(gs *GameState, p byte) bool {
	if gs.WhiteToMove {
		return isBlack(p)
	}
	return isWhite(p)
}

func (gs *GameState) GetValidMoves() []Move {
	moves := gs.GetAllPossibleMoves()
	var validMoves []Move

	for _, move := range moves {
		gs.MakeMove(move)
		gs.WhiteToMove = !gs.WhiteToMove // Switch back to check our own King
		if !gs.InCheck() {
			validMoves = append(validMoves, move)
		}
		gs.WhiteToMove = !gs.WhiteToMove
		gs.UndoMove()
	}

	if len(validMoves) == 0 {
		if gs.InCheck() {
			gs.Checkmate = true
		} else {
			gs.Stalemate = true
		}
	} else {
		gs.Checkmate = false
		gs.Stalemate = false
	}

	// Add Castling at the very end to prevent infinite recursion!
	if gs.WhiteToMove {
		gs.GetCastleMoves(gs.WhiteKingLoc[0], gs.WhiteKingLoc[1], &validMoves)
	} else {
		gs.GetCastleMoves(gs.BlackKingLoc[0], gs.BlackKingLoc[1], &validMoves)
	}

	return validMoves
}

func (gs *GameState) InCheck() bool {
	if gs.WhiteToMove {
		return gs.SquareUnderAttack(gs.WhiteKingLoc[0], gs.WhiteKingLoc[1])
	}
	return gs.SquareUnderAttack(gs.BlackKingLoc[0], gs.BlackKingLoc[1])
}

func (gs *GameState) SquareUnderAttack(r, c int) bool {
	gs.WhiteToMove = !gs.WhiteToMove // Switch to enemy point of view
	oppMoves := gs.GetAllPossibleMoves()
	gs.WhiteToMove = !gs.WhiteToMove

	for _, move := range oppMoves {
		if move.EndRow == r && move.EndCol == c {
			return true
		}
	}
	return false
}

func (gs *GameState) GetAllPossibleMoves() []Move {
	var moves []Move // We start with an empty slice

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := gs.Board[r][c]
			if piece == '-' {
				continue
			}
			if (isWhite(piece) && gs.WhiteToMove) || (isBlack(piece) && !gs.WhiteToMove) {
				switch piece {
				case 'P', 'p':
					gs.GetPawnMoves(r, c, &moves)
				case 'R', 'r':
					gs.GetRookMoves(r, c, &moves)
				case 'N', 'n':
					gs.GetKnightMoves(r, c, &moves)
				case 'B', 'b':
					gs.GetBishopMoves(r, c, &moves)
				case 'Q', 'q':
					gs.GetQueenMoves(r, c, &moves) // Queen just uses Rook + Bishop
				case 'K', 'k':
					gs.GetKingMoves(r, c, &moves)
				}
			}
		}
	}
	return moves
}

// --- INDIVIDUAL PIECE MOVES ---

func (gs *GameState) GetPawnMoves(r, c int, moves *[]Move) {
	piece := gs.Board[r][c]

	if gs.WhiteToMove {
		// 1 square advance
		if r-1 >= 0 && gs.Board[r-1][c] == '-' {
			*moves = append(*moves, Move{r, c, r - 1, c, piece, '-', r-1 == 0, false, false, 0})
			// 2 square advance
			if r == 6 && gs.Board[r-2][c] == '-' {
				*moves = append(*moves, Move{r, c, r - 2, c, piece, '-', false, false, false, 0})
			}
		}
		// Captures
		if r-1 >= 0 && c-1 >= 0 && isBlack(gs.Board[r-1][c-1]) {
			*moves = append(*moves, Move{r, c, r - 1, c - 1, piece, gs.Board[r-1][c-1], r-1 == 0, false, false, 0})
		}
		if r-1 >= 0 && c+1 < 8 && isBlack(gs.Board[r-1][c+1]) {
			*moves = append(*moves, Move{r, c, r - 1, c + 1, piece, gs.Board[r-1][c+1], r-1 == 0, false, false, 0})
		}
		// En Passant
		if r-1 == gs.EnPassantPossible[0] && c-1 == gs.EnPassantPossible[1] {
			*moves = append(*moves, Move{r, c, r - 1, c - 1, piece, 'p', false, true, false, 0})
		}
		if r-1 == gs.EnPassantPossible[0] && c+1 == gs.EnPassantPossible[1] {
			*moves = append(*moves, Move{r, c, r - 1, c + 1, piece, 'p', false, true, false, 0})
		}
	} else {
		// BLACK PAWNS
		if r+1 < 8 && gs.Board[r+1][c] == '-' {
			*moves = append(*moves, Move{r, c, r + 1, c, piece, '-', r+1 == 7, false, false, 0})
			if r == 1 && gs.Board[r+2][c] == '-' {
				*moves = append(*moves, Move{r, c, r + 2, c, piece, '-', false, false, false, 0})
			}
		}
		if r+1 < 8 && c-1 >= 0 && isWhite(gs.Board[r+1][c-1]) {
			*moves = append(*moves, Move{r, c, r + 1, c - 1, piece, gs.Board[r+1][c-1], r+1 == 7, false, false, 0})
		}
		if r+1 < 8 && c+1 < 8 && isWhite(gs.Board[r+1][c+1]) {
			*moves = append(*moves, Move{r, c, r + 1, c + 1, piece, gs.Board[r+1][c+1], r+1 == 7, false, false, 0})
		}
		if r+1 == gs.EnPassantPossible[0] && c-1 == gs.EnPassantPossible[1] {
			*moves = append(*moves, Move{r, c, r + 1, c - 1, piece, 'P', false, true, false, 0})
		}
		if r+1 == gs.EnPassantPossible[0] && c+1 == gs.EnPassantPossible[1] {
			*moves = append(*moves, Move{r, c, r + 1, c + 1, piece, 'P', false, true, false, 0})
		}
	}
}

func (gs *GameState) GetRookMoves(r, c int, moves *[]Move) {
	directions := [4][2]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
	gs.getSlidingMoves(r, c, moves, directions[:])
}

func (gs *GameState) GetBishopMoves(r, c int, moves *[]Move) {
	directions := [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	gs.getSlidingMoves(r, c, moves, directions[:])
}

func (gs *GameState) GetQueenMoves(r, c int, moves *[]Move) {
	gs.GetRookMoves(r, c, moves)
	gs.GetBishopMoves(r, c, moves)
}

// Reusable slider logic for Rooks, Bishops, and Queens
func (gs *GameState) getSlidingMoves(r, c int, moves *[]Move, directions [][2]int) {
	piece := gs.Board[r][c]
	for _, d := range directions {
		for i := 1; i < 8; i++ {
			endRow := r + d[0]*i
			endCol := c + d[1]*i
			if endRow >= 0 && endRow < 8 && endCol >= 0 && endCol < 8 {
				endPiece := gs.Board[endRow][endCol]
				if endPiece == '-' {
					*moves = append(*moves, Move{r, c, endRow, endCol, piece, '-', false, false, false, 0})
				} else if isEnemy(gs, endPiece) {
					*moves = append(*moves, Move{r, c, endRow, endCol, piece, endPiece, false, false, false, 0})
					break
				} else {
					break
				}
			} else {
				break
			}
		}
	}
}

func (gs *GameState) GetKnightMoves(r, c int, moves *[]Move) {
	directions := [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	gs.getSteppingMoves(r, c, moves, directions[:])
}

func (gs *GameState) GetKingMoves(r, c int, moves *[]Move) {
	directions := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	gs.getSteppingMoves(r, c, moves, directions[:])
}

// Reusable stepping logic for Knights and Kings
func (gs *GameState) getSteppingMoves(r, c int, moves *[]Move, directions [][2]int) {
	piece := gs.Board[r][c]
	for _, d := range directions {
		endRow := r + d[0]
		endCol := c + d[1]
		if endRow >= 0 && endRow < 8 && endCol >= 0 && endCol < 8 {
			endPiece := gs.Board[endRow][endCol]
			if endPiece == '-' || isEnemy(gs, endPiece) {
				*moves = append(*moves, Move{r, c, endRow, endCol, piece, endPiece, false, false, false, 0})
			}
		}
	}
}

func (gs *GameState) GetCastleMoves(r, c int, moves *[]Move) {
	if gs.SquareUnderAttack(r, c) {
		return
	}
	if (gs.WhiteToMove && gs.CurrentCastleRights.Wks) || (!gs.WhiteToMove && gs.CurrentCastleRights.Bks) {
		if gs.Board[r][c+1] == '-' && gs.Board[r][c+2] == '-' {
			if !gs.SquareUnderAttack(r, c+1) && !gs.SquareUnderAttack(r, c+2) {
				*moves = append(*moves, Move{r, c, r, c + 2, gs.Board[r][c], '-', false, false, true, 0})
			}
		}
	}
	if (gs.WhiteToMove && gs.CurrentCastleRights.Wqs) || (!gs.WhiteToMove && gs.CurrentCastleRights.Bqs) {
		if gs.Board[r][c-1] == '-' && gs.Board[r][c-2] == '-' && gs.Board[r][c-3] == '-' {
			if !gs.SquareUnderAttack(r, c-1) && !gs.SquareUnderAttack(r, c-2) {
				*moves = append(*moves, Move{r, c, r, c - 2, gs.Board[r][c], '-', false, false, true, 0})
			}
		}
	}
}
