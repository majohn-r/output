package output

import "fmt"

// ListDecorator contains the data needed for creating list decorations
type ListDecorator struct {
	enabled    bool
	numeric    bool
	itemNumber uint8
}

func newListDecorator(enabled, numeric bool) *ListDecorator {
	return &ListDecorator{
		enabled:    enabled,
		numeric:    numeric,
		itemNumber: 1, // correct if numeric is true, unused if false
	}
}

// Decorator generates the appropriate decoration for lists (and typically, this is the empty string)
func (ld *ListDecorator) Decorator() string {
	if !ld.enabled {
		return ""
	}
	if ld.numeric {
		s := fmt.Sprintf("%2d. ", ld.itemNumber)
		ld.itemNumber++
		return s
	}
	return "● "
}
