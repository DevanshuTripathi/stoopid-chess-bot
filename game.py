import pygame
import engine
import random
import subprocess

# Size
WIDTH = 512
HEIGHT = 512
DIMENSION = 8

SQ_SIZE = WIDTH // DIMENSION
MAX_FPS = 15

IMAGES = {}

def load_sprites() :
    pieces = ['wp', 'wR', 'wN', 'wB', 'wQ', 'wK', 'bp', 'bR', 'bN', 'bB', 'bQ', 'bK']
    for piece in pieces:
        # We load the image, then immediately scale it to the size of our board squares
        image = pygame.image.load(f"assets/{piece}.png")
        IMAGES[piece] = pygame.transform.scale(image, (SQ_SIZE, SQ_SIZE))

def draw_game_state(screen, gs, validMoves, sqSelected) :
    draw_board(screen)
    highlight_squares(screen, gs, validMoves, sqSelected)
    draw_pieces(screen, gs.board)

def draw_board(screen) :
    colors = [pygame.Color("white"), pygame.Color("dark gray")]

    for row in range(DIMENSION):
        for col in range(DIMENSION):
            color = colors[((row + col) % 2)]
            pygame.draw.rect(screen, color, pygame.Rect(col * SQ_SIZE, row * SQ_SIZE, SQ_SIZE, SQ_SIZE))

def draw_pieces(screen, board) :
    for row in range(DIMENSION):
        for col in range(DIMENSION):
            piece = board[row][col]
            if piece != "--": # If the square isn't empty
                screen.blit(IMAGES[piece], pygame.Rect(col * SQ_SIZE, row * SQ_SIZE, SQ_SIZE, SQ_SIZE))

def highlight_squares(screen, gs, validMoves, sqSelected) :
    if sqSelected != () :
        r, c = sqSelected

        if gs.board[r][c][0] == ('w' if gs.whiteToMove else 'b') :
            s = pygame.Surface((SQ_SIZE, SQ_SIZE))
            s.set_alpha(100)
            s.fill(pygame.Color('blue'))
            screen.blit(s, (c * SQ_SIZE, r * SQ_SIZE))

            s.fill(pygame.Color('yellow'))
            for move in validMoves:
                if move.startRow == r and move.startCol == c:
                    screen.blit(s, (move.endCol * SQ_SIZE, move.endRow * SQ_SIZE))

def draw_end_game_text(screen, text):
    # Initialize a Pygame font (SysFont uses system fonts, 'Helvetica' or 'Arial' is safe)
    font = pygame.font.SysFont("Helvetica", 32, True, False)
    
    # Render the drop shadow first (slightly offset)
    textObject = font.render(text, 0, pygame.Color('Black'))
    textLocation = pygame.Rect(0, 0, WIDTH, HEIGHT).move(
                                                            WIDTH / 2 - textObject.get_width() / 2 + 2, 
                                                            HEIGHT / 2 - textObject.get_height() / 2 + 2
                                                        )
    screen.blit(textObject, textLocation)
    
    # Render the actual text on top
    textObject = font.render(text, 0, pygame.Color('Dark Orange')) # Feel free to change the color!
    screen.blit(textObject, textLocation.move(-2, -2))

def main() :
    pygame.init()
    screen = pygame.display.set_mode((WIDTH, HEIGHT))
    pygame.display.set_caption("Stoopid Chess Bot")
    clock = pygame.time.Clock()

    gs = engine.GameState()
    validMoves = gs.getValidMoves()
    moveMade = False

    load_sprites()
    sqSelected = ()
    playerClicks = []

    running = True
    playerOne = True
    playerTwo = False
    gameOver = False

    go_engine = subprocess.Popen(['./engine.exe'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True)

    while running:
        humanTurn = (gs.whiteToMove and playerOne) or (not gs.whiteToMove and playerTwo)

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                go_engine.stdin.write("quit\n") # Tell Go to shut down
                go_engine.stdin.flush()
                go_engine.terminate() 
                running = False
            
            elif event.type == pygame.MOUSEBUTTONDOWN :
                location = pygame.mouse.get_pos()
                col = location[0] // SQ_SIZE
                row = location[1] // SQ_SIZE

                if sqSelected == (row, col):
                    sqSelected = ()
                    playerClicks = []
                
                else:
                    sqSelected = (row, col)
                    playerClicks.append(sqSelected)
                
                if len(playerClicks) == 2 :
                    move = engine.Move(playerClicks[0], playerClicks[1], gs.board)

                    for i in range(len(validMoves)):
                        if move == validMoves[i] :
                            gs.makeMove(validMoves[i])
                            moveMade = True
                            sqSelected = ()
                            playerClicks = []
                            print(move.getChessNotation())
                            break
                    if not moveMade:
                        playerClicks = [sqSelected]

                    sqSelected = () 
                    playerClicks = []
            
            elif event.type == pygame.KEYDOWN:
                if event.key == pygame.K_r: # Press 'r' to reset the board
                    gs = engine.GameState()
                    validMoves = gs.getValidMoves()
                    sqSelected = ()
                    playerClicks = []
                    moveMade = False
                    gameOver = False
        
        if not gameOver and not humanTurn:
            # 1. Prepare the move history to sync the Go Engine
            move_history = [m.getChessNotation() + ('q' if m.isPawnPromotion else '') for m in gs.moveLog]
            
            # 2. Send position to Go and wait for 'ready' signal
            go_engine.stdin.write(f"position moves {' '.join(move_history)}\n")
            go_engine.stdin.flush()
            
            # Wait for the Go engine to confirm synchronization
            while True:
                response = go_engine.stdout.readline().strip()
                if response == "ready":
                    break
            
            # 3. Tell Go to calculate the best move
            go_engine.stdin.write("go\n")
            go_engine.stdin.flush()
            
            # 4. Read the best move from Go (e.g., "bestmove e7e5")
            best_move_raw = go_engine.stdout.readline().strip()
            
            if best_move_raw.startswith("bestmove") and best_move_raw != "bestmove none":
                move_code = best_move_raw.split(" ")[1]
                
                # 5. Apply the Go move to your local Python GameState
                for m in validMoves:
                    # Match the notation (including promotion suffix)
                    promo = 'q' if m.isPawnPromotion else ''
                    if m.getChessNotation() + promo == move_code:
                        gs.makeMove(m)
                        moveMade = True
                        break
    
        if moveMade:
            validMoves = gs.getValidMoves()
            moveMade = False

            if gs.checkmate or gs.stalemate:
                gameOver = True
                if gs.stalemate:
                    text = 'Stalemate!'
                else:
                    text = 'Black wins by Checkmate!' if gs.whiteToMove else 'White wins by Checkmate!'
                draw_end_game_text(screen, text)

        draw_game_state(screen, gs, validMoves, sqSelected)

        if gameOver:
            draw_end_game_text(screen, text)

        pygame.display.flip()
        clock.tick(MAX_FPS)

    pygame.quit()

if __name__ == "__main__" :
    main()