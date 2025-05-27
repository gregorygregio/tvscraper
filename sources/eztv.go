package eztv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"tvscraper/clients"
	"tvscraper/utils"
	documentcrawler "tvscraper/utils"

	"golang.org/x/net/html"
)

var mainUrl string = "https://eztvx.to"

func SendMagnetToQBitTorrent(magnetLink string) {
	fmt.Println("Iniciando SentMagnetToQBitTorrent")

	username, _ := utils.GetConfig("username")
	password, _ := utils.GetConfig("password")

	qbitClient := &clients.QBitTorrentClient{}
	qbitClient.SetSettings("192.168.0.182", "8081", username, password, false)

	err := qbitClient.SendMagnet(magnetLink)
	if err != nil {
		fmt.Println("Erro ao enviar magnet para qbitTorrent")
	} else {
		fmt.Println("Magnet enviado com sucesso")
	}

	fmt.Println("Fim SentMagnetToQBitTorrent")
}

func getFromUrl(targetUrl string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		fmt.Printf("Error creating request to %s: %s\n", targetUrl, err.Error())
		return nil, fmt.Errorf("error creating request to %s", targetUrl)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request to %s: %s\n", targetUrl, err.Error())
		return nil, fmt.Errorf("error making request to %s", targetUrl)
	}

	return resp, nil
}

func FetchEztvSeason(seriesName string, season string) (*[]string, error) {

	fmt.Printf("Fetching season %s for series %s from EZTV...\n", season, seriesName)

	search := fmt.Sprintf("%s-%s", strings.Replace(seriesName, " ", "-", -1), getSeasonString(season))
	targetUrl := fmt.Sprintf("%s/search/%s", getBaseUrl(), search)

	fmt.Printf("Fetching URL: %s \n", targetUrl)

	resp, err := getFromUrl(targetUrl)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %s\n", targetUrl, err.Error())
		return nil, fmt.Errorf("error fetching URL %s", targetUrl)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from %s", resp.StatusCode, targetUrl)
		return nil, fmt.Errorf("error: received status code %d from %s", resp.StatusCode, targetUrl)
	}

	foundLinks, err := parseSearchResults(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing search results: %s\n", err.Error())
		return nil, fmt.Errorf("error parsing search results")
	}

	return getMagnetLinks(foundLinks)
}

func parseSearchResults(reader io.Reader) (*[]string, error) {

	doc, err := documentcrawler.NewDocumentCrawler(reader)
	if err != nil {
		fmt.Printf("Error creating document crawler: %s\n", err.Error())
		return nil, fmt.Errorf("creating document crawler")
	}

	results := &[]string{}
	doc.ForEachElement(func(n *html.Node) {
		//fmt.Printf("Parsing %s\n", n.Data)
		if n.Type == html.ElementNode && n.Data == "a" && documentcrawler.HasClass(n, "epinfo") {
			//TODO filtrar o episódio para que não haja repetição do mesmo ep e para selecionar a resolução correta
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Printf("Found link: %s\n", attr.Val)
					*results = append(*results, attr.Val)
				}
			}
		}
	})

	for i := 0; i < len(*results); i++ {
		fmt.Printf("Result link: %s\n", (*results)[i])
	}

	return results, nil
}

func getMagnetLinks(links *[]string) (*[]string, error) {

	var wg sync.WaitGroup
	results := &[]string{}

	for _, link := range *links {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := getMagnetLink(link)
			if err == nil {
				*results = append(*results, result)
			}
		}()
	}

	wg.Wait()
	for i := 0; i < len(*results); i++ {
		fmt.Printf("Magnets Result link: %s\n", (*results)[i])
	}
	return results, nil
}

func getMagnetLink(link string) (string, error) {
	epUrl := getBaseUrl() + link
	fmt.Printf("Fetching episode URL: %s\n", epUrl)
	resp, err := getFromUrl(epUrl)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %s\n", epUrl, err.Error())
		return "", fmt.Errorf("error fetching URL %s", epUrl)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from %s", resp.StatusCode, epUrl)
		return "", fmt.Errorf("error: received status code %d from %s", resp.StatusCode, epUrl)
	}

	doc, err := documentcrawler.NewDocumentCrawler(resp.Body)
	if err != nil {
		fmt.Printf("Error creating document crawler: %s\n", err.Error())
		return "", fmt.Errorf("creating document crawler")
	}

	torrentLink := ""
	doc.ForEachElement(func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" && documentcrawler.HasAttr(n, "title", "Magnet Link") {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Printf("Found torrent link: %s\n", attr.Val)
					torrentLink = attr.Val
				}
			}
		}
	})

	return torrentLink, nil
}

func writeToFile(filename string, data string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", filename, err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("Error writing to file %s: %s\n", filename, err.Error())
		return
	}

	fmt.Printf("Data written to file %s\n", filename)
}

func getBaseUrl() string {
	return mainUrl
}
func getSeasonString(season string) string {
	if len(season) == 1 {
		return "s0" + season
	}
	return "s" + season
}
