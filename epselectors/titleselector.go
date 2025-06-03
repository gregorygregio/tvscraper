package epselectors

import "strings"

func (sel *EpSelector) AcceptTitle() bool {
	//TODO verificar título todo até S01E01 para ver se não é maior
	return strings.Contains(sel.CleanTitle, sel.SeriesName)
}
