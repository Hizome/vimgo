package vim

import (
	"strconv"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Visual
	Command
)

func (m Mode) String() string {
	switch m {
	case Normal:
		return "NORMAL"
	case Insert:
		return "INSERT"
	case Visual:
		return "VISUAL"
	case Command:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

type Handler struct {
	Mode          Mode
	CursorX       int
	CursorY       int
	BoardSize     int
	InputBuffer   string
	CommandBuffer string
	RepeatCount   int
}

func NewHandler(boardSize int) *Handler {
	return &Handler{
		Mode:      Normal,
		CursorX:   boardSize / 2,
		CursorY:   boardSize / 2,
		BoardSize: boardSize,
	}
}

// Action represents an intent derived from a keypress.
type Action struct {
	Type  ActionType
	Value string
	Count int
}

type ActionType int

const (
	ActionMove ActionType = iota
	ActionPlaceStone
	ActionEnterMode
	ActionCommand
	ActionUndo
	ActionRedo
	ActionPass
)

func (h *Handler) HandleKey(key string) *Action {
	switch h.Mode {
	case Normal:
		return h.handleNormalKey(key)
	case Command:
		return h.handleCommandKey(key)
	case Insert:
		return h.handleInsertKey(key)
	}
	return nil
}

func (h *Handler) handleInsertKey(key string) *Action {
	if key == "esc" {
		h.Mode = Normal
		return &Action{Type: ActionEnterMode, Value: "NORMAL"}
	}
	return nil
}

func (h *Handler) handleNormalKey(key string) *Action {
	// Handle numbers for repeat count
	if _, err := strconv.Atoi(key); err == nil {
		h.InputBuffer += key
		h.RepeatCount, _ = strconv.Atoi(h.InputBuffer)
		return nil
	}

	count := h.RepeatCount
	if count == 0 {
		count = 1
	}
	h.InputBuffer = ""
	h.RepeatCount = 0

	switch key {
	case "h":
		h.CursorX = max(0, h.CursorX-count)
		return &Action{Type: ActionMove}
	case "l":
		h.CursorX = min(h.BoardSize-1, h.CursorX+count)
		return &Action{Type: ActionMove}
	case "j":
		h.CursorY = min(h.BoardSize-1, h.CursorY+count)
		return &Action{Type: ActionMove}
	case "k":
		h.CursorY = max(0, h.CursorY-count)
		return &Action{Type: ActionMove}
	case "x":
		return &Action{Type: ActionPlaceStone, Count: count}
	case ":":
		h.Mode = Command
		h.CommandBuffer = ""
		return &Action{Type: ActionEnterMode, Value: "COMMAND"}
	case "u":
		return &Action{Type: ActionUndo, Count: count}
	case "i":
		h.Mode = Insert
		return &Action{Type: ActionEnterMode, Value: "INSERT"}
	}
	return nil
}

func (h *Handler) handleCommandKey(key string) *Action {
	if key == "enter" {
		cmd := h.CommandBuffer
		h.Mode = Normal
		h.CommandBuffer = ""
		return &Action{Type: ActionCommand, Value: cmd}
	}
	if key == "esc" {
		h.Mode = Normal
		h.CommandBuffer = ""
		return &Action{Type: ActionEnterMode, Value: "NORMAL"}
	}
	if key == "backspace" {
		if len(h.CommandBuffer) > 0 {
			h.CommandBuffer = h.CommandBuffer[:len(h.CommandBuffer)-1]
		}
		return nil
	}
	if len(key) == 1 {
		h.CommandBuffer += key
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
