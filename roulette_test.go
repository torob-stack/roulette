package main

import "testing"

func TestColorOf(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "green"},
		{1, "red"}, {2, "black"}, {3, "red"},
		{11, "black"}, {12, "red"}, {13, "black"}, {14, "red"},
		{36, "red"},
	}
	for _, tt := range tests {
		if got := colorOf(tt.n); got != tt.want {
			t.Errorf("colorOf(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestEvaluateNumber(t *testing.T) {
	tests := []struct {
		choice string
		win    int
		want   int
	}{
		{"17", 17, 35},
		{"0", 0, 35},
		{"18", 17, 0},
	}
	for _, tt := range tests {
		got := EvaluateNumber(Bet{Type: "number", Choice: tt.choice, Stake: 100}, tt.win, "")
		if got != tt.want {
			t.Errorf("EvaluateNumber(choice=%s, win=%d) = %d, want %d", tt.choice, tt.win, got, tt.want)
		}
	}
}

func TestEvaluateColor(t *testing.T) {
	tests := []struct {
		choice string
		winNum int
		winCol string
		want   int
	}{
		{"red", 3, "red", 1},
		{"black", 22, "black", 1},
		{"red", 22, "black", 0},
		{"black", 0, "green", 0}, // 0 is green, should lose for red/black
	}
	for _, tt := range tests {
		got := EvaluateColor(Bet{Type: "color", Choice: tt.choice, Stake: 100}, tt.winNum, tt.winCol)
		if got != tt.want {
			t.Errorf("EvaluateColor(choice=%s, win=%d/%s) = %d, want %d", tt.choice, tt.winNum, tt.winCol, got, tt.want)
		}
	}
}

func TestEvaluateOddEven(t *testing.T) {
	tests := []struct {
		choice string
		win    int
		want   int
	}{
		{"odd", 17, 1},
		{"even", 18, 1},
		{"odd", 0, 0},  // 0 loses
		{"even", 0, 0}, // 0 loses
		{"odd", 18, 0},
	}
	for _, tt := range tests {
		got := EvaluateOddEven(Bet{Type: "odd_even", Choice: tt.choice, Stake: 100}, tt.win, "")
		if got != tt.want {
			t.Errorf("EvaluateOddEven(%s, %d) = %d, want %d", tt.choice, tt.win, got, tt.want)
		}
	}
}

func TestEvaluateLowHigh(t *testing.T) {
	tests := []struct {
		choice string
		win    int
		want   int
	}{
		{"low", 1, 1}, {"low", 18, 1}, {"low", 19, 0}, {"low", 0, 0},
		{"high", 19, 1}, {"high", 36, 1}, {"high", 18, 0}, {"high", 0, 0},
	}
	for _, tt := range tests {
		got := EvaluateLowHigh(Bet{Type: "low_high", Choice: tt.choice, Stake: 100}, tt.win, "")
		if got != tt.want {
			t.Errorf("EvaluateLowHigh(%s, %d) = %d, want %d", tt.choice, tt.win, got, tt.want)
		}
	}
}

func TestEvaluateDozen(t *testing.T) {
	tests := []struct {
		choice string
		win    int
		want   int
	}{
		{"1st", 1, 2}, {"1st", 12, 2}, {"1st", 13, 0}, {"1st", 0, 0},
		{"2nd", 13, 2}, {"2nd", 24, 2}, {"2nd", 25, 0},
		{"3rd", 25, 2}, {"3rd", 36, 2}, {"3rd", 24, 0},
	}
	for _, tt := range tests {
		got := EvaluateDozen(Bet{Type: "dozen", Choice: tt.choice, Stake: 100}, tt.win, "")
		if got != tt.want {
			t.Errorf("EvaluateDozen(%s, %d) = %d, want %d", tt.choice, tt.win, got, tt.want)
		}
	}
}

func TestEvaluateColumn(t *testing.T) {
	// Column mapping: col1 => numbers with (n%3==1), col2 => (n%3==2), col3 => (n%3==0)
	tests := []struct {
		choice string
		win    int
		want   int
	}{
		{"col1", 1, 2}, {"col1", 4, 2}, {"col1", 34, 2}, {"col1", 2, 0}, {"col1", 0, 0},
		{"col2", 2, 2}, {"col2", 5, 2}, {"col2", 35, 2}, {"col2", 3, 0}, {"col2", 0, 0},
		{"col3", 3, 2}, {"col3", 6, 2}, {"col3", 36, 2}, {"col3", 4, 0}, {"col3", 0, 0},
	}
	for _, tt := range tests {
		got := EvaluateColumn(Bet{Type: "column", Choice: tt.choice, Stake: 100}, tt.win, "")
		if got != tt.want {
			t.Errorf("EvaluateColumn(%s, %d) = %d, want %d", tt.choice, tt.win, got, tt.want)
		}
	}
}

func TestSumStakes(t *testing.T) {
	bets := []Bet{
		{Type: "color", Choice: "red", Stake: 150},
		{Type: "number", Choice: "17", Stake: 25},
		{Type: "dozen", Choice: "2nd", Stake: 100},
	}
	if got, want := sumStakes(bets), 275; got != want {
		t.Errorf("sumStakes = %d, want %d", got, want)
	}
}

func TestMoney(t *testing.T) {
	got := money(12345)
	want := "£123.45"
	if got != want {
		t.Errorf("money(12345)=%s want %s", got, want)
	}
}

func TestZeroLosesCommonBets(t *testing.T) {
	if got := EvaluateOddEven(Bet{Type: "odd_even", Choice: "odd", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("odd on 0 should lose")
	}
	if got := EvaluateOddEven(Bet{Type: "odd_even", Choice: "even", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("even on 0 should lose")
	}
	if got := EvaluateLowHigh(Bet{Type: "low_high", Choice: "low", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("low on 0 should lose")
	}
	if got := EvaluateLowHigh(Bet{Type: "low_high", Choice: "high", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("high on 0 should lose")
	}
	if got := EvaluateColor(Bet{Type: "color", Choice: "red", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("red on 0 should lose")
	}
	if got := EvaluateColor(Bet{Type: "color", Choice: "black", Stake: 100}, 0, "green"); got != 0 {
		t.Errorf("black on 0 should lose")
	}
}
