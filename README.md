# Roulette CLI Game (Go)

A fully interactive **command-line roulette game** built in Go.  
This project simulates a European roulette wheel with multiple bet types, a visual table display, color-coded results, and balance tracking — all in a simple, text-based interface.

---

## Features

- **Multiple bet types**:
  - Straight (single number)
  - Red/Black
  - Odd/Even
  - Low/High (1–18 / 19–36)
  - Dozens (1st, 2nd, 3rd)
  - Columns (Col1, Col2, Col3)
- **Accurate European roulette colors & payouts**
- **Animated wheel spin** showing number sequence
- **Color-coded results** in the terminal
- **ASCII roulette table** for easy reference before betting
- Tracks **player balance** and ends the game when funds run out
- Supports **multiple bets in one round**
- Fully **tested** core game logic (`roulette_test.go`)

---

## Demo

```
Current balance: £100.00
Type '1' to play roulette or anything else to quit: 1

Roulette Table
(Col1   Col2   Col3)
  +----+
  |  0 |
  +----+
+-------------------+
|  1     2     3    |
|  4     5     6    |
...
+-------------------+
types of play & payouts:
Col1 / Col2 / Col3 ................. 2:1
1–12 / 13–24 / 25–36 (dozens) ...... 2:1
Low(1–18) / High(19–36) ............ 1:1
Odd / Even ......................... 1:1
Red / Black ........................ 1:1
Straight (single number) ........... 35:1
```

---

## Installation

1. **Clone or download** this repository.
2. Make sure you have **Go installed** (1.20+ recommended).  
   You can check with:
   ```bash
   go version
   ```
3. Run the game:
   ```bash
   go run main.go
   ```

---

## Running Tests

Unit tests are included for key game logic.

```bash
go test
```

---

## Project Structure

```
.
├── main.go             # Game logic and CLI
├── roulette_test.go    # Unit tests for game logic
├── go.mod              # Go module definition
├── README.md           # This file
└── .gitignore          # Ignore compiled binaries & temp files
```

---

## About

This project was built as a **portfolio piece** to showcase:
- CLI-based game design
- Structuring a Go project
- Using ASCII art for improved UX in terminal apps
- Writing and running unit tests in Go
- Clean, maintainable, and documented code
