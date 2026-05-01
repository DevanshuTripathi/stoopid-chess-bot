class GameState() :
    def __init__(self) :
        self.board = [
            ["bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"],
            ["bp", "bp", "bp", "bp", "bp", "bp", "bp", "bp"],
            ["--", "--", "--", "--", "--", "--", "--", "--"],
            ["--", "--", "--", "--", "--", "--", "--", "--"],
            ["--", "--", "--", "--", "--", "--", "--", "--"],
            ["--", "--", "--", "--", "--", "--", "--", "--"],
            ["wp", "wp", "wp", "wp", "wp", "wp", "wp", "wp"],
            ["wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"]
        ]
        self.whiteToMove = True
        self.moveLog = []
        self.currentCastlingRight = CastleRights(True, True, True, True)
        self.castleRightsLog = [CastleRights( 
                                                self.currentCastlingRight.wks, 
                                                self.currentCastlingRight.bks, 
                                                self.currentCastlingRight.wqs, 
                                                self.currentCastlingRight.bqs
                                            )]
        self.enPassantPossible = ()
        self.whiteKingLocation = (7, 4) # Hardcoded starting locations
        self.blackKingLocation = (0, 4)
        self.checkmate = False
        self.stalemate = False
        self.enPassantPossibleLog = [self.enPassantPossible]
    
    def getValidMoves(self) :
        moves = self.getAllPossibleMoves()

        for i in range(len(moves) - 1, -1, -1): # Traverse backwards!
            self.makeMove(moves[i])
            
            # 3. Generate all opponent's moves and see if they attack our king
            self.whiteToMove = not self.whiteToMove # Switch turns back to check our own king
            if self.inCheck():
                moves.remove(moves[i]) # 4. If they attack your king, not a valid move!
            self.whiteToMove = not self.whiteToMove # Switch turns back
            
            # 5. Undo move
            self.undoMove()
            
        # --- Checkmate & Stalemate Detection ---
        if len(moves) == 0:
            if self.inCheck():
                self.checkmate = True
            else:
                self.stalemate = True
        else:
            self.checkmate = False
            self.stalemate = False
        
        if self.whiteToMove:
            self.getCastleMoves(self.whiteKingLocation[0], self.whiteKingLocation[1], moves)
        else:
            self.getCastleMoves(self.blackKingLocation[0], self.blackKingLocation[1], moves)
            
        return moves
    
    def getAllPossibleMoves(self) :
        moves = []
        for r in range(len(self.board)) :
            for c in range(len(self.board[r])) :
                turn = self.board[r][c][0] # Looks at the first character ('w', 'b', or '-')
                
                # If it's a piece belonging to the person whose turn it is
                if (turn == 'w' and self.whiteToMove) or (turn == 'b' and not self.whiteToMove) :
                    piece = self.board[r][c][1]
                    
                    if piece == 'p' :
                        self.getPawnMoves(r, c, moves)
                    elif piece == 'R' :
                        self.getRookMoves(r, c, moves)
                    elif piece == 'B' :
                        self.getBishopMoves(r, c, moves)
                    elif piece == 'Q' :
                        self.getQueenMoves(r, c, moves)
                    elif piece == 'N':
                        self.getKnightMoves(r, c, moves)
                    elif piece == 'K':
                        self.getKingMoves(r, c, moves)
        return moves
    
    def inCheck(self):
        """Determine if the current player is in check"""
        if self.whiteToMove:
            return self.squareUnderAttack(self.whiteKingLocation[0], self.whiteKingLocation[1])
        else:
            return self.squareUnderAttack(self.blackKingLocation[0], self.blackKingLocation[1])

    def squareUnderAttack(self, r, c):
        """Determine if the enemy can attack the square (r, c)"""
        self.whiteToMove = not self.whiteToMove # Switch to opponent's point of view
        oppMoves = self.getAllPossibleMoves()
        self.whiteToMove = not self.whiteToMove # Switch back
        
        for move in oppMoves:
            if move.endRow == r and move.endCol == c: # Square is under attack
                return True
        return False
    
    def getPawnMoves(self, r, c, moves) :
        if self.whiteToMove :
            if self.board[r-1][c] == "--" :
                moves.append(Move((r, c), (r-1, c), self.board))
                if r == 6 and self.board[r-2][c] == "--" :
                    moves.append(Move((r, c), (r-2, c), self.board))
            
            if c - 1 >= 0 : 
                if self.board[r-1][c-1][0] == 'b' : # Enemy piece to capture
                    moves.append(Move((r, c), (r-1, c-1), self.board))

            if c + 1 <= 7 : 
                if self.board[r-1][c+1][0] == 'b' : 
                    moves.append(Move((r, c), (r-1, c+1), self.board))
            
            if (r - 1, c - 1) == self.enPassantPossible:
                moves.append(Move((r, c), (r - 1, c - 1), self.board))
            if (r - 1, c + 1) == self.enPassantPossible:
                moves.append(Move((r, c), (r - 1, c + 1), self.board))
        
        else : # BLACK PAWN LOGIC
            if self.board[r+1][c] == "--" : # 1 square pawn advance
                moves.append(Move((r, c), (r+1, c), self.board))
                if r == 1 and self.board[r+2][c] == "--" : # 2 square advance
                    moves.append(Move((r, c), (r+2, c), self.board))
                    
            # Captures to the left
            if c - 1 >= 0 : 
                if self.board[r+1][c-1][0] == 'w' :
                    moves.append(Move((r, c), (r+1, c-1), self.board))
            # Captures to the right
            if c + 1 <= 7 : 
                if self.board[r+1][c+1][0] == 'w' :
                    moves.append(Move((r, c), (r+1, c+1), self.board))
            
            if (r + 1, c - 1) == self.enPassantPossible:
                moves.append(Move((r, c), (r + 1, c - 1), self.board))
            if (r + 1, c + 1) == self.enPassantPossible:
                moves.append(Move((r, c), (r + 1, c + 1), self.board))

    def getRookMoves(self, r, c, moves) :
        directions = ((-1, 0), (0, -1), (1, 0), (0, 1))

        enemyColor = "b" if self.whiteToMove else "w"

        for d in directions :
            for i in range(1, 8) :
                endRow = r + d[0] * i
                endCol = c + d[1] * i

                if 0 <= endRow < 8 and 0 <= endCol < 8 :
                    endPiece = self.board[endRow][endCol]

                    if endPiece == "--" :
                        moves.append(Move((r, c), (endRow, endCol), self.board))
                    elif endPiece[0] == enemyColor :
                        moves.append(Move((r, c), (endRow, endCol), self.board))
                        break
                    else :
                        break
                else :
                    break
    
    def getBishopMoves(self, r, c, moves) :
        directions = ((-1, -1), (-1, 1), (1, -1), (1, 1))

        enemyColor = "b" if self.whiteToMove else "w"
    
        for d in directions:
            for i in range(1, 8): # A bishop can move a maximum of 7 squares
                endRow = r + d[0] * i
                endCol = c + d[1] * i
                
                if 0 <= endRow < 8 and 0 <= endCol < 8: # Is it on board?
                    endPiece = self.board[endRow][endCol]
                    
                    if endPiece == "--": # Empty space
                        moves.append(Move((r, c), (endRow, endCol), self.board))
                    elif endPiece[0] == enemyColor: # Enemy piece
                        moves.append(Move((r, c), (endRow, endCol), self.board))
                        break # Stop sliding after a capture
                    else: # Friendly piece
                        break # Stop sliding
                else: # Off board
                    break
    
    def getQueenMoves(self, r, c, moves):
        """
        Get all the queen moves for the queen located at row r and col c and add these moves to the list.
        """
        self.getRookMoves(r, c, moves)
        self.getBishopMoves(r, c, moves)
    
    def getKnightMoves(self, r, c, moves) :
        knightMoves = ((-2, -1), (-2, 1), (-1, -2), (-1, 2), (1, -2), (1, 2), (2, -1), (2, 1))
        allyColor = "w" if self.whiteToMove else "b"

        for m in knightMoves:
            endRow = r + m[0]
            endCol = c + m[1]
            
            if 0 <= endRow < 8 and 0 <= endCol < 8: # Is it on board?
                endPiece = self.board[endRow][endCol]
                if endPiece[0] != allyColor: # If it's NOT an ally piece (so it's empty or enemy)
                    moves.append(Move((r, c), (endRow, endCol), self.board))
    
    def getKingMoves(self, r, c, moves) :
        kingMoves = ((-1, -1), (-1, 0), (-1, 1), (0, -1), (0, 1), (1, -1), (1, 0), (1, 1))
        allyColor = "w" if self.whiteToMove else "b"

        for i in range(8):
            endRow = r + kingMoves[i][0]
            endCol = c + kingMoves[i][1]
            
            if 0 <= endRow < 8 and 0 <= endCol < 8: # Is it on board?
                endPiece = self.board[endRow][endCol]
                if endPiece[0] != allyColor: # Not an ally piece
                    moves.append(Move((r, c), (endRow, endCol), self.board))
    
    def getCastleMoves(self, r, c, moves):
        if self.squareUnderAttack(r, c):
            return
        if (self.whiteToMove and self.currentCastlingRight.wks) or (not self.whiteToMove and self.currentCastlingRight.bks):
            self.getKingsideCastleMoves(r, c, moves)
        if (self.whiteToMove and self.currentCastlingRight.wqs) or (not self.whiteToMove and self.currentCastlingRight.bqs):
            self.getQueensideCastleMoves(r, c, moves)
    
    def getKingsideCastleMoves(self, r, c, moves):
        if self.board[r][c+1] == '--' and self.board[r][c+2] == '--':
            if not self.squareUnderAttack(r, c+1) and not self.squareUnderAttack(r, c+2):
                moves.append(Move((r, c), (r, c+2), self.board, isCastleMove=True))
    
    def getQueensideCastleMoves(self, r, c, moves):
        if self.board[r][c-1] == '--' and self.board[r][c-2] == '--' and self.board[r][c-3] == '--':
            if not self.squareUnderAttack(r, c-1) and not self.squareUnderAttack(r, c-2):
                moves.append(Move((r, c), (r, c-2), self.board, isCastleMove=True))

    def makeMove(self, move) :
        self.board[move.startRow][move.startCol] = "--"
        self.board[move.endRow][move.endCol] = move.pieceMoved
        self.moveLog.append(move)
        self.whiteToMove = not self.whiteToMove

        if move.pieceMoved == 'wK':
            self.whiteKingLocation = (move.endRow, move.endCol)
        elif move.pieceMoved == 'bK':
            self.blackKingLocation = (move.endRow, move.endCol)

        if move.isPawnPromotion:
            self.board[move.endRow][move.endCol] = move.pieceMoved[0] + 'Q'

        if move.isEnpassantMove:
            self.board[move.startRow][move.endCol] = "--"

        if move.pieceMoved[1] == 'p' and abs(move.startRow - move.endRow) == 2:
            self.enPassantPossible = ((move.startRow + move.endRow) // 2, move.startCol)
        else:
            self.enPassantPossible = ()
        
        self.updateCastleRights(move)
        self.castleRightsLog.append(CastleRights(
                                                    self.currentCastlingRight.wks, 
                                                    self.currentCastlingRight.bks, 
                                                    self.currentCastlingRight.wqs, 
                                                    self.currentCastlingRight.bqs
                                                ))
        
        if move.isCastleMove:
            if move.endCol - move.startCol == 2: # Kingside castle
                self.board[move.endRow][move.endCol - 1] = self.board[move.endRow][move.endCol + 1] # Move Rook
                self.board[move.endRow][move.endCol + 1] = '--' # Erase old Rook
            else: # Queenside castle
                self.board[move.endRow][move.endCol + 1] = self.board[move.endRow][move.endCol - 2] # Move Rook
                self.board[move.endRow][move.endCol - 2] = '--'
        
        self.enPassantPossibleLog.append(self.enPassantPossible)

    def undoMove(self):
        """Undo the last move made"""
        if len(self.moveLog) != 0: # Make sure there is a move to undo
            move = self.moveLog.pop()
            self.board[move.startRow][move.startCol] = move.pieceMoved
            self.board[move.endRow][move.endCol] = move.pieceCaptured
            self.whiteToMove = not self.whiteToMove # Swap turns back
            
            # Undo King Location
            if move.pieceMoved == 'wK':
                self.whiteKingLocation = (move.startRow, move.startCol)
            elif move.pieceMoved == 'bK':
                self.blackKingLocation = (move.startRow, move.startCol)
                
            # Undo En Passant
            if move.isEnpassantMove:
                self.board[move.endRow][move.endCol] = '--' # Leave landing square blank
                self.board[move.startRow][move.endCol] = move.pieceCaptured # Put pawn back
                
            self.enPassantPossibleLog.pop()
            self.enPassantPossible = self.enPassantPossibleLog[-1]
            
            # Undo Castling Rights
            self.castleRightsLog.pop() 
            newRights = self.castleRightsLog[-1]
            self.currentCastlingRight = CastleRights(newRights.wks, newRights.bks, newRights.wqs, newRights.bqs)
            
            # Undo Castling Move
            if move.isCastleMove:
                if move.endCol - move.startCol == 2: # Kingside
                    self.board[move.endRow][move.endCol + 1] = self.board[move.endRow][move.endCol - 1]
                    self.board[move.endRow][move.endCol - 1] = '--'
                else: # Queenside
                    self.board[move.endRow][move.endCol - 2] = self.board[move.endRow][move.endCol + 1]
                    self.board[move.endRow][move.endCol + 1] = '--'
    
    def updateCastleRights(self, move):
        """Update the castle rights given the move"""
        if move.pieceMoved == 'wK':
            self.currentCastlingRight.wks = False
            self.currentCastlingRight.wqs = False
        elif move.pieceMoved == 'bK':
            self.currentCastlingRight.bks = False
            self.currentCastlingRight.bqs = False
        elif move.pieceMoved == 'wR':
            if move.startRow == 7:
                if move.startCol == 0: # Left rook
                    self.currentCastlingRight.wqs = False
                elif move.startCol == 7: # Right rook
                    self.currentCastlingRight.wks = False
        elif move.pieceMoved == 'bR':
            if move.startRow == 0:
                if move.startCol == 0: # Left rook
                    self.currentCastlingRight.bqs = False
                elif move.startCol == 7: # Right rook
                    self.currentCastlingRight.bks = False

class Move() :
    ranksToRows = {"1": 7, "2": 6, "3": 5, "4": 4, "5": 3, "6": 2, "7": 1, "8": 0}
    rowsToRanks = {v: k for k, v in ranksToRows.items()}
    filesToCols = {"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7}
    colsToFiles = {v: k for k, v in filesToCols.items()}

    def __init__(self, startSq, endSq, board, isCastleMove=False) :
        self.startRow = startSq[0]
        self.startCol = startSq[1]
        self.endRow = endSq[0]
        self.endCol = endSq[1]
        self.pieceMoved = board[self.startRow][self.startCol]
        self.pieceCaptured = board[self.endRow][self.endCol]

        self.isCastleMove = isCastleMove

        self.isPawnPromotion = False
        if (self.pieceMoved == 'wp' and self.endRow == 0) or (self.pieceMoved == 'bp' and self.endRow == 7):
            self.isPawnPromotion = True
        
        self.isEnpassantMove = False
        if self.pieceMoved[1] == 'p' and (self.startCol != self.endCol) and self.pieceCaptured == "--":
            self.isEnpassantMove = True
    
    def __eq__(self, other) :
        if isinstance(other, Move) :
            return self.startRow == other.startRow and self.startCol == other.startCol and self.endRow == other.endRow and self.endCol == other.endCol
        return False
    
    def getChessNotation(self) :
        return self.getRankFile(self.startRow, self.startCol) + self.getRankFile(self.endRow, self.endCol)
    
    def getRankFile(self, r, c) :
        return self.colsToFiles[c] + self.rowsToRanks[r]

class CastleRights() :
    def __init__(self, wks, bks, wqs, bqs) :
        self.wks = wks # White King Side
        self.bks = bks # Black King Side
        self.wqs = wqs # White Queen Side
        self.bqs = bqs # Black Queen Side