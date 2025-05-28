package epselectors

import (
	"regexp"
	"strings"
	"tvscraper/models"
)

type EpSelector struct {
	Episode    *models.EpInfo
	CleanTitle string
	SeriesName string
}

var nonWord = regexp.MustCompile(`[\W]`)
var specialChars = regexp.MustCompile(`[^\d\p{Latin}\+]`)
var doublePlus = regexp.MustCompile(`\+{2,}`)

func (sel *EpSelector) cleanName(epname string) string {
	epname = strings.ReplaceAll(epname, "&", "and")
	epname = nonWord.ReplaceAllString(epname, "+")
	epname = specialChars.ReplaceAllString(epname, "")

	epname = doublePlus.ReplaceAllString(epname, "")

	return strings.ToLower(strings.Trim(epname, "+"))
}

func (sel *EpSelector) AcceptOrReject() bool {
	sel.CleanTitle = sel.cleanName(sel.Episode.EpName)
	sel.SeriesName = sel.cleanName(sel.SeriesName)
	//fmt.Printf("Episode title cleaned %s\n", sel.CleanTitle)

	if !sel.AcceptQuality() {
		return false
	}
	if !sel.AcceptTitle() {
		return false
	}

	return true
}
