# Stoopid Chess Bot

**Stoopid Chess Bot** is a hybrid chess engine that combines a responsive **Pygame frontend** with a high-performance **Golang backend**.

- **Python + Pygame** for graphics, animations, and user input  
- **Go** for raw speed and deep search performance  

The result: a chess bot capable of searching **millions of positions per second** using **Negamax + Alpha-Beta pruning**.

## Features

- Playable graphical chess interface using Pygame  
- Strong AI powered by Golang engine  
- Fast move generation and evaluation  
- Smooth Python ↔ Go communication through standard I/O  

## Prerequisites

- Python 3.x
- Pygame

## Installation

```bash
git clone https://github.com/DevanshuTripathi/stoopid-chess-bot.git
cd stoopid-chess-bot
pip install -r requirements.txt
```

## Running the Game

```bash
python game.py
```

## Engine Details

### Search Techniques

- Negamax Search  
- Alpha-Beta Pruning  
- Move Ordering (MVV-LVA)  
- Quiescence Search  
- Distance to Mate  

### Evaluation

- Material Values  
- Piece-Square Tables (PST)

