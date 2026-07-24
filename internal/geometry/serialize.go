package geometry

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func SerializeCommands(commands []Command) (string, error) {
	if len(commands) == 0 {
		return "", fmt.Errorf("cannot serialize an empty command stream")
	}
	absolute, err := encodeCommands(commands, false)
	if err != nil {
		return "", err
	}
	relative, err := encodeCommands(commands, true)
	if err != nil {
		return "", err
	}
	candidates := []string{absolute, relative}
	if len(relative) < len(absolute) {
		candidates[0], candidates[1] = relative, absolute
	}
	var failures []string
	for _, candidate := range candidates {
		parsed, parseErr := ParsePath(candidate)
		if parseErr != nil {
			failures = append(failures, parseErr.Error())
			continue
		}
		if commandsEqual(commands, parsed) {
			return candidate, nil
		}
		failures = append(failures, "serialized path did not round-trip")
	}
	return "", fmt.Errorf("no round-tripping serialization: %s", strings.Join(failures, "; "))
}

func encodeCommands(commands []Command, relative bool) (string, error) {
	var b strings.Builder
	var x, y float64
	lastOp := ""
	for _, c := range commands {
		want := map[string]int{"M": 2, "L": 2, "Q": 4, "Z": 0}[c.Op]
		if _, ok := map[string]int{"M": 2, "L": 2, "Q": 4, "Z": 0}[c.Op]; !ok || len(c.Values) != want {
			return "", fmt.Errorf("invalid command %q", c.Op)
		}
		op := c.Op
		if relative {
			op = strings.ToLower(op)
		}
		omit := c.Op == lastOp && (c.Op == "L" || c.Op == "Q")
		if !omit {
			b.WriteString(op)
		}
		vals := append([]float64(nil), c.Values...)
		if relative {
			switch c.Op {
			case "M", "L":
				vals[0] -= x
				vals[1] -= y
			case "Q":
				vals[0] -= x
				vals[1] -= y
				vals[2] -= x
				vals[3] -= y
			}
		}
		for i, v := range vals {
			writePathNumber(&b, v, i > 0 || omit)
		}
		switch c.Op {
		case "M", "L":
			x, y = c.Values[0], c.Values[1]
		case "Q":
			x, y = c.Values[2], c.Values[3]
		}
		if c.Op == "Z" {
			lastOp = ""
		} else {
			lastOp = c.Op
		}
	}
	return b.String(), nil
}

func writePathNumber(b *strings.Builder, v float64, separator bool) {
	s := number(v)
	if strings.HasPrefix(s, "-0.") {
		s = "-" + strings.TrimPrefix(s, "-0")
	} else if strings.HasPrefix(s, "0.") {
		s = strings.TrimPrefix(s, "0")
	}
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	if s == "-0" || s == "-" || s == "" {
		s = "0"
	}
	if separator && len(s) > 0 && s[0] != '-' && s[0] != '+' {
		b.WriteByte(' ')
	}
	b.WriteString(s)
}

func ParsePath(path string) ([]Command, error) {
	var out []Command
	i := 0
	op := byte(0)
	var x, y float64
	for {
		skipSpaces := func() {
			for i < len(path) && (unicode.IsSpace(rune(path[i])) || path[i] == ',') {
				i++
			}
		}
		skipSpaces()
		if i >= len(path) {
			break
		}
		if isCommand(path[i]) {
			op = path[i]
			i++
			if op == 'Z' || op == 'z' {
				out = append(out, Command{Op: "Z"})
				op = 0
				continue
			}
		}
		if op == 0 {
			return nil, fmt.Errorf("missing command at %d", i)
		}
		upper := strings.ToUpper(string(op))
		n := 2
		if upper == "Q" {
			n = 4
		}
		values := make([]float64, n)
		for k := 0; k < n; k++ {
			skipSpaces()
			start := i
			if i < len(path) && (path[i] == '-' || path[i] == '+') {
				i++
			}
			dots := 0
			for i < len(path) && ((path[i] >= '0' && path[i] <= '9') || path[i] == '.') {
				if path[i] == '.' {
					dots++
				}
				i++
			}
			if start == i || dots > 1 {
				return nil, fmt.Errorf("invalid number at %d", start)
			}
			v, err := strconv.ParseFloat(path[start:i], 64)
			if err != nil {
				return nil, err
			}
			values[k] = v
		}
		if op >= 'a' && op <= 'z' {
			for k := 0; k < len(values); k += 2 {
				values[k] += x
				values[k+1] += y
			}
		}
		out = append(out, Command{Op: upper, Values: values})
		if upper == "M" || upper == "L" {
			x, y = values[0], values[1]
		} else {
			x, y = values[2], values[3]
		}
		if upper == "M" {
			if op == 'M' {
				op = 'L'
			} else {
				op = 'l'
			}
		}
	}
	return out, nil
}
func isCommand(b byte) bool { return strings.ContainsRune("MmLlQqZz", rune(b)) }
func commandsEqual(a, b []Command) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Op != b[i].Op || len(a[i].Values) != len(b[i].Values) {
			return false
		}
		for j := range a[i].Values {
			if mathAbs(a[i].Values[j]-b[i].Values[j]) > .005000000001 {
				return false
			}
		}
	}
	return true
}
