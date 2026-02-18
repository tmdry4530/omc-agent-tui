package arena

import "github.com/chamdom/omc-agent-tui/pkg/schema"

// SpriteLines is the fixed 3-line CLCO mascot sprite.
// Every agent card renders this at the top, tinted by role color.
var SpriteLines = [3]string{
	"\u2590\u259B\u2588\u2588\u2588\u259C\u258C",     // ▐▛███▜▌
	"\u259D\u259C\u2588\u2588\u2588\u2588\u2588\u259B\u2598", // ▝▜█████▛▘
	"\u2598\u2598   \u259D\u259D",                      // ▘▘   ▝▝
}

// ASCIIFallback is the 3-line ASCII fallback when block chars break.
var ASCIIFallback = [3]string{
	" /===\\ ",
	"|#####|",
	" \\   / ",
}

// GetSprite returns the 3-line sprite (unicode or ASCII fallback).
func GetSprite(useUnicode bool) [3]string {
	if useUnicode {
		return SpriteLines
	}
	return ASCIIFallback
}

// stateIndicators maps agent states to unicode/ASCII indicator pairs.
var stateIndicators = map[schema.AgentState][2]string{
	schema.StateRunning:   {"\u25CF", "*"},
	schema.StateWaiting:   {"\u25CB", "o"},
	schema.StateBlocked:   {"\u26A0", "!"},
	schema.StateError:     {"\u2718", "x"},
	schema.StateDone:      {"\u2714", "+"},
	schema.StateIdle:      {"\u2500", "-"},
	schema.StateFailed:    {"\u2716", "X"},
	schema.StateCancelled: {"\u2205", "~"},
}

// GetStateIndicator returns the unicode or ASCII indicator for a state.
func GetStateIndicator(state schema.AgentState, useUnicode bool) string {
	if pair, ok := stateIndicators[state]; ok {
		if useUnicode {
			return pair[0]
		}
		return pair[1]
	}
	if useUnicode {
		return "\u2500"
	}
	return "-"
}
