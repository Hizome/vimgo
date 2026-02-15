package sgf

import (
	"fmt"
	"strings"
)

// ToSGFCoord converts x, y coordinates to SGF style (e.g., 0,0 -> "aa").
func ToSGFCoord(x, y int) string {
	if x < 0 || y < 0 {
		return ""
	}
	return string(rune('a'+x)) + string(rune('a'+y))
}

// FromSGFCoord converts SGF style coordinates to x, y (e.g., "pd" -> 15, 3).
func FromSGFCoord(coord string) (int, int, error) {
	if len(coord) != 2 {
		return -1, -1, fmt.Errorf("invalid SGF coord length: %s", coord)
	}
	x := int(coord[0] - 'a')
	y := int(coord[1] - 'a')
	return x, y, nil
}

// EncodeMove creates an SGF move string, e.g., "B[pd]" or "W[aa]".
func EncodeMove(color string, x, y int) string {
	if x < 0 || y < 0 {
		return fmt.Sprintf("%s[]", color) // Pass
	}
	return fmt.Sprintf("%s[%s]", color, ToSGFCoord(x, y))
}

// SimpleSGFWriter creates a basic SGF string from a sequence of moves.
func SimpleSGFWriter(size int, moves []string) string {
	var sb strings.Builder
	sb.WriteString("(;GM[1]FF[4]CA[UTF-8]")
	sb.WriteString(fmt.Sprintf("SZ[%d]", size))
	for _, m := range moves {
		sb.WriteString(";")
		sb.WriteString(m)
	}
	sb.WriteString(")")
	return sb.String()
}

// ParseSGF parses a simple SGF string and returns the board size and a list of moves.
func ParseSGF(content string) (int, []string, error) {
	size := 19
	var moves []string

	inVal := false
	propKey := ""
	valBuf := ""
	
	// Normalize content
	content = strings.TrimSpace(content)
	
	for i := 0; i < len(content); i++ {
		char := content[i]
		
		if inVal {
			if char == ']' {
				inVal = false
				// Process property
				switch propKey {
				case "SZ":
					fmt.Sscanf(valBuf, "%d", &size)
				case "B", "W":
					moves = append(moves, fmt.Sprintf("%s[%s]", propKey, valBuf))
				case "AB", "AW": 
					// Setup stones, treat as moves for visual simplicity or handle separately? 
					// For now, let's just focus on B/W moves.
				}
			} else {
				valBuf += string(char)
			}
			continue
		}

		switch char {
		case '[':
			inVal = true
			valBuf = ""
		case '(', ')', ';':
			propKey = ""
		default:
			if char > ' ' { 
				propKey += string(char)
			}
		}
	}

	return size, moves, nil
}
