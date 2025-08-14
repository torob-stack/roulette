package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	var balance float64 = 100
	for {
		fmt.Print("Type '1' to play roulette or anything else to quit: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "1" {
			fmt.Println("OK bye bye")
			break
		}
		game(reader, balance)
	}
}

func game(reader *bufio.Reader, balance float64) float64 {
	bets := []Bet{}
	remaining := balance

	// Collect one or more bets
	for {
		betType := getOneOf(reader, "Choose bet type (number/color): ", []string{"number", "color"})

		switch betType {
		case "number":
			bets = append(bets, CollectNumber(reader, remaining))
		case "color":
			bets = append(bets, CollectColor(reader, remaining))
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

	// Commit stakes
	totalStake := sumStakes(bets)
	balance -= totalStake

	// Spin + derive properties
	winNum := roulette()
	winColor := colorOf(winNum)
	fmt.Printf("\n🎡 Result: %d (%s)\n", winNum, winColor)

	// Settle bets
	totalWinnings := 0.0
	for _, b := range bets {
		eval, ok := evaluators[b.Type]
		if !ok {
			fmt.Printf("Skipping unknown bet type: %s\n", b.Type)
			continue
		}
		mult := eval(b, winNum, winColor)
		won := b.Stake * mult
		totalWinnings += won
		fmt.Printf("Bet: %-6s %-8s Stake: £%.2f  ->  Payout x%.1f  = £%.2f\n",
			b.Type, b.Choice, b.Stake, mult, won)
	}

	// Update balance and show summary
	balance += totalWinnings
	net := totalWinnings - totalStake
	fmt.Printf("\nSummary: Staked £%.2f | Returned £%.2f | Net %+0.2f | New balance £%.2f\n\n",
		totalStake, totalWinnings, net, balance)

	return balance
}

func CollectColor(reader *bufio.Reader, remaining float64) Bet {
	color := getOneOf(reader, "Pick a color (red/black): ", []string{"red", "black"})
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "color", Choice: color, Stake: stake}
}

func CollectNumber(reader *bufio.Reader, remaining float64) Bet {
	num := getIntInRange(reader, "Pick a number (0–36): ", 0, 36)
	stake := getStakeWithinBalance(reader, "Stake £", remaining)
	return Bet{Type: "number", Choice: strconv.Itoa(num), Stake: stake}
}

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

func getStake(reader *bufio.Reader, prompt string) float64 {
	for {
		fmt.Print(prompt)
		in := readLine(reader)
		n, err := strconv.ParseFloat(in, 64)
		if err == nil && n > 0 { // will add balance check later check if positive
			return n
		}
		fmt.Println("invalid! Stake must be positive")
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

type Bet struct {
	Type   string //number color oddeven
	Choice string // 17? red? odd?
	Stake  float64
}

type Evaluator func(b Bet, winNum int, winColor string) float64

var evaluators = map[string]Evaluator{
	"number": EvaluateNumber,
	"color":  EvaluateColor,
}

func colorOf(n int) string {
	red := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
	if n == 0 {
		return "green"
	}
	if slices.Contains(red, n) {
		return "red"
	} else {
		return "black"
	}
}

func EvaluateNumber(b Bet, winNum int, _ string) float64 {
	num, _ := strconv.Atoi(b.Choice)
	if winNum == num {
		return 36.0
	} else {
		return 0.0
	}
}

func EvaluateColor(b Bet, _ int, winColor string) float64 {
	if winColor == b.Choice {
		return 2.0
	} else {
		return 0.0
	}
}

func sumStakes(bets []Bet) float64 {
	total := 0.0
	for _, b := range bets {
		total += b.Stake
	}
	return total
}

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
func getStakeWithinBalance(reader *bufio.Reader, prompt string, remaining float64) float64 {
	for {
		fmt.Printf("%s (available £%.2f): ", prompt, remaining)
		in := readLine(reader)
		n, err := strconv.ParseFloat(in, 64)
		if err == nil && n > 0 && n <= remaining {
			return n
		}
		fmt.Println("Invalid amount. Must be > 0 and ≤ available.")
	}
}

func roulette() int {
	roulette := rand.Intn(37)
	return roulette
}
