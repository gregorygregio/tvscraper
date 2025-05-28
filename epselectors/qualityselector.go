package epselectors

import (
	"fmt"
	"strings"
)

var availableResolutions []string = []string{
	"240p",
	"360p",
	"480p",
	"720p",
	"1080p",
	"1440p",
	"2160p",
	"4320p",
}

var minimumResolution string = "480p"
var preferedResolution string = "1080p"

var optionIndexesLoaded = false
var minimumResolutionIndex int16
var preferedResolutionIndex int16

func (sel *EpSelector) AcceptQuality() bool {
	loadResolutionOptions()

	epResolution := sel.getResolution()
	fmt.Printf("%s EpResolution: %s\n", sel.CleanTitle, epResolution)
	if epResolution == preferedResolution {
		return true
	}

	for i := minimumResolutionIndex; i < int16(len(availableResolutions)); i++ {
		if epResolution == availableResolutions[i] {
			return true
		}
	}

	return false
}

func (sel *EpSelector) getResolution() string {
	for _, resolution := range availableResolutions {
		if strings.Contains(sel.CleanTitle, "+"+resolution+"+") {
			return resolution
		}
	}

	return ""
}

func loadResolutionOptions() {
	if optionIndexesLoaded {
		return
	}

	for resolutionIndex, resolution := range availableResolutions {
		if resolution == preferedResolution {
			preferedResolutionIndex = int16(resolutionIndex)
		}
		if resolution == minimumResolution {
			minimumResolutionIndex = int16(resolutionIndex)
		}
	}
	optionIndexesLoaded = true
}
