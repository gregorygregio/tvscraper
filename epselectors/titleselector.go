package epselectors

import "strings"

func (sel *EpSelector) AcceptTitle() bool {
	return strings.Contains(sel.CleanTitle, sel.SeriesName)
}
