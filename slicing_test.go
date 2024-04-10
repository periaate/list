package list

import (
	"fmt"
	"log/slog"
	"testing"
)

func TestSlice(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	longLen := 1000
	baseInput := []string{"0", "1", "2", "3", "4", "5"}
	longInput := make([]string, 0, longLen)
	for i := range longLen {
		longInput = append(longInput, fmt.Sprintf("%v", i))
	}
	bl := len(baseInput)

	exp := []struct {
		pattern     string
		input       []string
		expected    []string
		fail        bool
		description string
	}{
		// Simple cases
		{
			"[:]",
			baseInput,
			baseInput,
			false,
			"empty slice",
		},
		{
			"[0:]",
			baseInput,
			baseInput,
			false,
			"from start",
		},
		{
			"[:0]",
			baseInput,
			[]string{},
			false,
			"to start",
		},
		// Positive index
		{
			"[2:4]",
			baseInput,
			baseInput[2:4],
			false,
			"middle slice",
		},
		{
			"[2:]",
			baseInput,
			baseInput[2:],
			false,
			"from middle",
		},
		// Negative index
		{
			"[-1]",
			baseInput,
			[]string{"5"},
			false,
			"select last element",
		},
		{
			"[-10]",
			longInput,
			[]string{"990"},
			false,
			"select last element, long input",
		},
		{
			fmt.Sprintf("[-1%s10]", string(pageToken)),
			longInput,
			longInput[longLen-1*10 : longLen],
			false,
			"select last 10, long input",
		},
		{
			"[-1:]",
			baseInput,
			[]string{"5"},
			false,
			"last element via negative index",
		},
		{
			"[:-1]",
			baseInput,
			baseInput[:bl-1],
			false,
			"all but last element",
		},
		{
			"[-3:-1]",
			baseInput,
			baseInput[bl-3 : bl-1],
			false,
			"last two elements via negative index",
		},
		// Relative cases
		{
			"[-3:+1]",
			baseInput,
			baseInput[bl-3 : bl-3+1],
			false,
			"last three elements via negative index",
		},
		{
			"[2:+3]",
			baseInput,
			baseInput[2 : 2+3],
			false,
			"middle slice via relative index",
		},
		// Pagination
		{
			fmt.Sprintf("[0%s50]", string(pageToken)),
			longInput,
			longInput[:1*50],
			false,
			"first page of 50",
		},
		{
			fmt.Sprintf("[1%s50]", string(pageToken)),
			longInput,
			longInput[1*50 : 2*50],
			false,
			"second page of 50",
		},
		{
			fmt.Sprintf("[1:3%s10]", string(pageToken)),
			longInput,
			longInput[1*10 : 3*10],
			false,
			"select two pages of 10",
		},
		{
			fmt.Sprintf("[-3:-1%s10]", string(pageToken)),
			longInput,
			longInput[longLen-3*10 : longLen-1*10],
			false,
			"last two pages of 10",
		},
		{
			fmt.Sprintf("[%s10]", string(pageToken)),
			longInput,
			longInput[0*10 : 1*10],
			false,
			"implicit page 0 of size 10",
		},
		// Fail cases
		{
			"[0:a]",
			baseInput,
			nil,
			true,
			"non-integer slice",
		},
		{
			"[0:1:-2]",
			baseInput,
			nil,
			true,
			"negative step",
		},
		{
			"[0:1:2]",
			baseInput,
			nil,
			true,
			"too many colons",
		},
	}

	for _, e := range exp {
		t.Run(e.description, func(t *testing.T) {
			out, err := Slice(e.pattern, e.input)
			if err != nil && !e.fail {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && e.fail {
				t.Errorf("expected error, got none")
			}
			if !e.fail && !equal(out, e.expected) {
				t.Errorf("expected %v, got %v", e.expected, out)
			}
		})
	}

}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
