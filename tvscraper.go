package main

import (
	"fmt"
	"tvscraper/clients"
	eztv "tvscraper/eztv"
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

	fmt.Printf("Found %v magnet links from 1EzTV\n", len(*magnetLinks))

	//TODO definir nome da pasta baseado na busca ou argumento
	//sendMagnetLinksToQBitTorrent(magnetLinks)

	fmt.Println("Fim")
}

func sendMagnetLinksToQBitTorrent(magnetLinks *[]string) {
	for _, mag := range *magnetLinks {
		username, _ := utils.GetConfig("username")
		password, _ := utils.GetConfig("password")

		qbitClient := &clients.QBitTorrentClient{}
		qbitClient.SetSettings("192.168.0.182", "8081", username, password, false)

		err := qbitClient.SendMagnet(mag)
		if err != nil {
			fmt.Println("Erro ao enviar magnet para qbitTorrent")
		} else {
			fmt.Println("Magnet enviado com sucesso")
		}
	}
}
