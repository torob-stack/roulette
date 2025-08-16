package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	flagSeed    = flag.Int64("seed", 0, "set RNG seed (0 = time-based)")
	flagFast    = flag.Bool("fast", false, "faster spin animation")
	flagNoColor = flag.Bool("no-color", false, "disable ANSI colors")
)

var useColor = true
var spinSleep = 120 * time.Millisecond

func main() {
	flag.Parse()

	if *flagSeed != 0 {
		rand.Seed(*flagSeed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	if *flagFast {
		spinSleep = 10 * time.Millisecond
	}
	if *flagNoColor || os.Getenv("NO_COLOR") != "" {
		useColor = false
	}

	reader := bufio.NewReader(os.Stdin)
	balance := 10000 // £100.00

	for {
		if balance <= 0 {
			fmt.Println("Game Over - you're out of funds")
			break
		}

		fmt.Printf("Current balance: %s\n", money(balance))
		fmt.Print("Type '1' to play roulette or anything else to quit: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "1" {
			fmt.Println("OK bye bye")
			break
		}
		balance = game(reader, balance) // capture updated balance
	}
}

func game(reader *bufio.Reader, balance int) int {
	// Show table before betting (no highlight yet)
	printAsciiTable(-1)

	var bets []Bet
	remaining := balance

	// Collect one or more bets
	for {
		betType := getOneOf(reader,
			"Choose bet type (number/color/odd_even/low_high/dozen/column): ",
			[]string{"number", "color", "odd_even", "low_high", "dozen", "column"},
		)

		switch betType {
		case "number":
			bets = append(bets, CollectNumber(reader, remaining))
		case "color":
			bets = append(bets, CollectColor(reader, remaining))
		case "odd_even":
			bets = append(bets, CollectOddEven(reader, remaining))
		case "low_high":
			bets = append(bets, CollectLowHigh(reader, remaining))
		case "dozen":
			bets = append(bets, CollectDozen(reader, remaining))
		case "column":
			bets = append(bets, CollectColumn(reader, remaining))
		}

		// Update remaining funds for this round
		remaining -= bets[len(bets)-1].Stake
		if remaining <= 0 {
			fmt.Println("You've used all available funds for this round.")
			break
		}

		if !getYesNo(reader, "Add another bet? (y/n): ") {
			break
		}
	}

	if len(bets) == 0 {
		fmt.Println("No bets placed. Balance unchanged.")
		return balance
	}

	printBetsTable(bets)

	// Commit stakes
	totalStake := sumStakes(bets)
	balance -= totalStake

	// Spin + derive properties
	winNum := roulette()
	winColor := colorOf(winNum)

	spinAnimation(winNum)
	fmt.Printf("\n🎡 Result: %s\n",
		colorText(fmt.Sprintf("%d (%s)", winNum, winColor), winColor))

	// Settle bets
	totalWinnings := 0
	for _, b := range bets {
		eval, ok := evaluators[b.Type]
		if !ok {
			fmt.Printf("Skipping unknown bet type: %s\n", b.Type)
			continue
		}
		// Net multiplier: 35 for number, 1 for color, etc.
		mult := eval(b, winNum, winColor)
		won := b.Stake * mult
		totalWinnings += won
		fmt.Printf("Bet: %-10s %-8s Stake: %s  ->  Net x%d  = %s\n",
			b.Type, b.Choice, money(b.Stake), mult, money(won))
	}

	// Update balance and show summary
	balance += totalWinnings

	net := totalWinnings - totalStake
	grossReturn := totalWinnings + totalStake // what came back to you this round

	fmt.Printf("\nSummary: Staked %s | Returned %s | Net %+s | New balance %s\n\n",
		money(totalStake), money(grossReturn), money(net), money(balance))

	return balance
}

// ---------- Collectors ----------

func CollectColor(reader *bufio.Reader, remaining int) Bet {
	color := getOneOf(reader, "Pick a color (red/black): ", []string{"red", "black"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "color", Choice: color, Stake: stake}
}

func CollectNumber(reader *bufio.Reader, remaining int) Bet {
	num := getIntInRange(reader, "Pick a number (0–36): ", 0, 36)
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "number", Choice: strconv.Itoa(num), Stake: stake}
}

func CollectOddEven(reader *bufio.Reader, remaining int) Bet {
	choice := getOneOf(reader, "Pick (odd/even): ", []string{"odd", "even"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "odd_even", Choice: choice, Stake: stake}
}

func CollectLowHigh(reader *bufio.Reader, remaining int) Bet {
	choice := getOneOf(reader, "Pick (low/high): ", []string{"low", "high"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "low_high", Choice: choice, Stake: stake}
}

func CollectDozen(reader *bufio.Reader, remaining int) Bet {
	choice := getOneOf(reader, "Pick dozen (1st/2nd/3rd): ", []string{"1st", "2nd", "3rd"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "dozen", Choice: choice, Stake: stake}
}

func CollectColumn(reader *bufio.Reader, remaining int) Bet {
	choice := getOneOf(reader, "Pick column (col1/col2/col3): ", []string{"col1", "col2", "col3"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "column", Choice: choice, Stake: stake}
}

// ---------- Input helpers ----------

func getIntInRange(r *bufio.Reader, prompt string, min int, max int) int {
	for {
		fmt.Print(prompt)
		in := readLine(r)
		n, err := strconv.Atoi(in)
		if err == nil && n >= min && n <= max {
			return n
		}
		fmt.Println("invalid input")
	}
}

func getOneOf(reader *bufio.Reader, prompt string, choices []string) string {
	for {
		fmt.Print(prompt)
		in := strings.ToLower(readLine(reader))
		for _, choice := range choices {
			if in == choice {
				return in
			}
		}
		fmt.Printf("invalid input please input: %v\n", choices)
	}
}

func readLine(r *bufio.Reader) string {
	str, _ := r.ReadString('\n')
	return strings.TrimSpace(str)
}

// ---------- Model & evaluation ----------

type Bet struct {
	Type   string // number | color | odd_even | low_high | dozen | column
	Choice string // e.g., "17", "red", "odd", "1st", "col2"
	Stake  int    // pennies
}

type Evaluator func(b Bet, winNum int, winColor string) int

var evaluators = map[string]Evaluator{
	"number":   EvaluateNumber,
	"color":    EvaluateColor,
	"odd_even": EvaluateOddEven,
	"low_high": EvaluateLowHigh,
	"dozen":    EvaluateDozen,
	"column":   EvaluateColumn,
}

func colorOf(n int) string {
	if n == 0 {
		return "green"
	}
	// True roulette red numbers
	red := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
	if slices.Contains(red, n) {
		return "red"
	}
	return "black"
}

// Return net multiplier (35 for a correct single number, else 0)
func EvaluateNumber(b Bet, winNum int, _ string) int {
	num, _ := strconv.Atoi(b.Choice)
	if winNum == num {
		return 35
	}
	return 0
}

// Return net multiplier (1 for a correct color, else 0)
func EvaluateColor(b Bet, _ int, winColor string) int {
	if strings.ToLower(winColor) == strings.ToLower(b.Choice) {
		return 1
	}
	return 0
}

func EvaluateOddEven(b Bet, winNum int, _ string) int {
	if winNum == 0 {
		return 0
	}
	switch b.Choice {
	case "odd":
		if winNum%2 == 1 {
			return 1
		}
	case "even":
		if winNum%2 == 0 {
			return 1
		}
	}
	return 0
}

func EvaluateLowHigh(b Bet, winNum int, _ string) int {
	if winNum == 0 {
		return 0
	}
	switch b.Choice {
	case "low":
		if 1 <= winNum && winNum <= 18 {
			return 1
		}
	case "high":
		if 19 <= winNum && winNum <= 36 {
			return 1
		}
	}
	return 0
}

func EvaluateDozen(b Bet, winNum int, _ string) int {
	if winNum == 0 {
		return 0
	}
	switch b.Choice {
	case "1st": // 1–12
		if 1 <= winNum && winNum <= 12 {
			return 2
		}
	case "2nd": // 13–24
		if 13 <= winNum && winNum <= 24 {
			return 2
		}
	case "3rd": // 25–36
		if 25 <= winNum && winNum <= 36 {
			return 2
		}
	}
	return 0
}

func EvaluateColumn(b Bet, winNum int, _ string) int {
	if winNum == 0 {
		return 0
	}
	// columns repeat every 3: col1 = 1,4,7...; col2 = 2,5,8...; col3 = 3,6,9...
	col := winNum % 3
	if col == 0 {
		col = 3
	}
	switch b.Choice {
	case "col1":
		if col == 1 {
			return 2
		}
	case "col2":
		if col == 2 {
			return 2
		}
	case "col3":
		if col == 3 {
			return 2
		}
	}
	return 0
}

func sumStakes(bets []Bet) int {
	total := 0
	for _, b := range bets {
		total += b.Stake
	}
	return total
}

// ---------- Yes/No & stake ----------

func getYesNo(reader *bufio.Reader, prompt string) bool {
	for {
		fmt.Print(prompt)
		in := strings.ToLower(readLine(reader))
		if in == "y" || in == "yes" {
			return true
		}
		if in == "n" || in == "no" {
			return false
		}
		fmt.Println("Please enter y/n.")
	}
}

// stake that cannot exceed remaining funds for this round
func getStakeWithinBalance(reader *bufio.Reader, prompt string, remaining int) int {
	for {
		fmt.Printf("%s (available %s): ", prompt, money(remaining))
		in := readLine(reader)

		// allow entering decimals like 2.50
		n, err := strconv.ParseFloat(in, 64)
		if err == nil && n > 0 {
			pennies := int(math.Round(n * 100))
			if pennies > 0 && pennies <= remaining {
				return pennies
			}
		}
		fmt.Println("Invalid amount. Must be > 0 and ≤ available.")
	}
}

// ---------- Spin / RNG ----------

func roulette() int {
	return rand.Intn(37)
}

// European wheel order (clockwise) starting at 0
var wheel = []int{0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30,
	8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26}

func spinAnimation(finalNum int) {
	// locate finalNum index
	idx := 0
	for i, n := range wheel {
		if n == finalNum {
			idx = i
			break
		}
	}
	// pick a start so we roll onto the result
	steps := 14 + rand.Intn(8) // 14–21 hops
	start := (idx - steps + len(wheel)*3) % len(wheel)

	fmt.Println("Spinning the wheel...")
	for s := 0; s <= steps; s++ {
		n := wheel[(start+s)%len(wheel)]
		c := colorOf(n)
		out := fmt.Sprintf("[%2d %5s]", n, c)
		if s == steps {
			out = colorText(out, c) // color final tick
		} else {
			out = ansiDim + out + ansiReset
		}
		fmt.Print(out)
		if s != steps {
			fmt.Print(" → ")
		}
		time.Sleep(spinSleep)
	}
	fmt.Println()
}

// ---------- Pretty printing / centering ----------

// ANSI colours (use bright gray for "black" so it's visible on dark themes)
const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
	ansiBlk   = "\033[90m" // bright black / gray
	ansiGrn   = "\033[32m"
	ansiDim   = "\033[2m"
)

func colorText(s, c string) string {
	if !useColor {
		return s
	}
	switch strings.ToLower(c) {
	case "red":
		return ansiRed + s + ansiReset
	case "black":
		return ansiBlk + s + ansiReset
	case "green":
		return ansiGrn + s + ansiReset
	default:
		return s
	}
}

// crude centering to ~80 columns; tweak width if you like
func centerLine(line string, width int) string {
	if len(line) >= width {
		return line
	}
	spaces := (width - len(line)) / 2
	return strings.Repeat(" ", spaces) + line
}

func centerPrintf(width int, format string, a ...any) {
	fmt.Print(centerLine(fmt.Sprintf(format, a...), width))
}

func printBetsTable(bets []Bet) {
	if len(bets) == 0 {
		return
	}
	width := 80
	fmt.Println(centerLine("\nYour bets this round:", width))
	fmt.Println(centerLine("--------------------------------", width))
	centerPrintf(width, "%-10s %-8s %8s\n", "Type", "Choice", "Stake")
	fmt.Println(centerLine("--------------------------------", width))
	total := 0
	for _, b := range bets {
		centerPrintf(width, "%-10s %-8s %8s\n", b.Type, b.Choice, money(b.Stake))
		total += b.Stake
	}
	fmt.Println(centerLine("--------------------------------", width))
	centerPrintf(width, "Total stake: %s\n\n", money(total))
}

// Pre-round left-aligned roulette table (like your mockup).
// Shows rows 1..36 (1,2,3 at the top), no highlight post-spin.
func printAsciiTable(_ int) {
	left := " " // overall left margin (tweak if you want)
	left0 := "    "
	left1 := "      "
	sep := "     " // spacing between columns

	// fixed-width numeric cell then colourize (keeps alignment)
	cell := func(n int) string {
		return colorText(fmt.Sprintf("%2d", n), colorOf(n)) // " 1", "12", "36"
	}

	fmt.Println()
	fmt.Println(left0 + "Roulette Table")
	fmt.Println(left + "(Col1   Col2   Col3)")

	// zero box, aligned to the grid
	fmt.Println(left1 + "  +----+")
	fmt.Println(left1 + "  | " + colorText(fmt.Sprintf("%2d", 0), "green") + " |")
	fmt.Println(left1 + "  +----+")

	// build one sample row to auto-size the border
	border := "+-----------------+"

	fmt.Println(left + border)
	for r := 1; r <= 12; r++ {
		c1 := 3*r - 2 // 1,4,7,...,34
		c2 := 3*r - 1 // 2,5,8,...,35
		c3 := 3 * r   // 3,6,9,...,36
		row := fmt.Sprintf("| %s%s%s%s%s|", cell(c1), sep, cell(c2), sep, cell(c3))
		fmt.Println(left + row)
	}
	fmt.Println(left + border)

	// legend
	fmt.Println(left + "types of play & payouts:")
	fmt.Println(left + "----------------------------------------------")
	fmt.Println(left + "Col1 / Col2 / Col3 ................. 2:1")
	fmt.Println(left + "1–12 / 13–24 / 25–36 (dozens) ...... 2:1")
	fmt.Println(left + "Low(1–18) / High(19–36) ............ 1:1")
	fmt.Println(left + "Odd / Even .......................... 1:1")
	fmt.Println(left + "Red / Black ......................... 1:1")
	fmt.Println(left + "Straight (single number) ............ 35:1")
	fmt.Println()
}

// money formats pennies as £#.##
func money(pennies int) string {
	return fmt.Sprintf("£%.2f", float64(pennies)/100.0)
}
