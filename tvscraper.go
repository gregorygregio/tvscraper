package main

import (
	"fmt"
	eztv "tvscraper/sources"
	"tvscraper/utils"
)

func main() {
	fmt.Println("Iniciando scraper")

	utils.LoadArguments()

	magnetLinks, err := eztv.FetchEztvSeason("The Office", "01")
	if err != nil {
		fmt.Println("Error fetching from Eztv: ", err)
		return
	}

	fmt.Printf("Found %v magnet links from EzTV\n", len(*magnetLinks))

	for _, mag := range *magnetLinks {
		eztv.SendMagnetToQBitTorrent(mag)
	}

	fmt.Println("Fim")
}
