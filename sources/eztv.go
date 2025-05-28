package eztv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	documentcrawler "tvscraper/utils"

	"golang.org/x/net/html"
)

// TODO criar mecanismo para verificar disponibilidade e buscar links alternativos
var mainUrl string = "https://eztvx.to"

type EpInfo struct {
	EpName string
	EpLink string
}

type EpisodesList struct {
	Episodes *[]EpInfo
}

func FetchEztvSeason(seriesName string, season string) (*[]string, error) {

	//provisório
	// episodeNumbers := make([]int16, 10)
	// for i := 0; i < len(episodeNumbers); i++ {
	// 	episodeNumbers[i] = (int16)(i + 1)
	// }

	episodeNumbers := make([]int16, 1)
	for i := 0; i < len(episodeNumbers); i++ {
		episodeNumbers[i] = (int16)(i + 1)
	}

	fmt.Printf("Fetching season %s for series %s from EZTV...\n", season, seriesName)
	foundLinks := make([]string, 0)

	for _, epNum := range episodeNumbers {
		epInfos, err := fetchEpisode(seriesName, season, epNum)
		if err != nil {
			fmt.Printf("Error fetching episode %v - season %s of '%s'\n", epNum, season, seriesName)
			continue
		}

		epInfo := findBestMatch(epInfos)
		if epInfo != nil {
			fmt.Printf("\nBest match for episode %v s%s is %s\n", epNum, season, epInfo.EpName)
			foundLinks = append(foundLinks, epInfo.EpLink)
		}
	}

	return getMagnetLinks(&foundLinks)
}

func fetchEpisode(seriesName string, season string, epNumber int16) (*[]EpInfo, error) {

	fmt.Printf("\n\nFetching episode %v from season %s of %s from EZTV...\n", epNumber, season, seriesName)
	search := fmt.Sprintf("%s-%s%s", strings.Replace(seriesName, " ", "-", -1), getSeasonString(season), getEpisodeNumberString(epNumber))
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

	foundEpisodes, err := parseSearchResults(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing search results: %s\n", err.Error())
		return nil, fmt.Errorf("error parsing search results")
	}

	return foundEpisodes, nil
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

func parseSearchResults(reader io.Reader) (*[]EpInfo, error) {

	doc, err := documentcrawler.NewDocumentCrawler(reader)
	if err != nil {
		fmt.Printf("Error creating document crawler: %s\n", err.Error())
		return nil, fmt.Errorf("creating document crawler")
	}

	results := &[]EpInfo{}
	doc.ForEachElement(func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" && documentcrawler.HasClass(n, "epinfo") {
			epInfo := &EpInfo{EpName: n.FirstChild.Data}

			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Printf("Found link: %s\n", attr.Val)
					epInfo.EpLink = attr.Val
					*results = append(*results, *epInfo)
				}
			}
		}
	})

	// for i := 0; i < len(*results); i++ {
	// 	fmt.Printf("Result link: %s\n", (*results)[i])
	// }

	return results, nil
}

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

func findBestMatch(episodes *[]EpInfo) *EpInfo {
	episodes = filterByResolution(episodes)
	return &(*episodes)[0]
}

func filterByResolution(episodes *[]EpInfo) *[]EpInfo {
	episodesByResolution := make(map[string]EpisodesList)
	preferedResolutionIndex := 0
	minimumResolutionIndex := 0

	for resolutionIndex, resolution := range availableResolutions {
		if resolution == preferedResolution {
			preferedResolutionIndex = resolutionIndex
		}
		if resolution == minimumResolution {
			minimumResolutionIndex = resolutionIndex
		}

		for _, ep := range *episodes {
			if strings.Contains(ep.EpName, " "+resolution+" ") {
				_, ok := episodesByResolution[resolution]
				if !ok {
					epList := make([]EpInfo, 0)
					episodesByResolution[resolution] = EpisodesList{Episodes: &epList}
				}
				//fmt.Printf("\nAdding ep %s to resolution %s\n", ep, resolution)
				*episodesByResolution[resolution].Episodes = append(*episodesByResolution[resolution].Episodes, ep)
			}
		}
	}

	epList, ok := episodesByResolution[preferedResolution]
	if ok {
		return epList.Episodes
	}

	if preferedResolutionIndex <= minimumResolutionIndex {
		return &[]EpInfo{}
	}

	for i := len(availableResolutions) - 1; i >= minimumResolutionIndex; i-- {
		res := availableResolutions[i]
		epList, ok := episodesByResolution[res]
		if ok {
			return epList.Episodes
		}
	}

	return &[]EpInfo{}
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
					fmt.Printf("Found torrent link on episode url %s\n", link)
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

func getEpisodeNumberString(epNum int16) string {
	if epNum < 10 {
		return "e" + strconv.Itoa(int(epNum))
	}
	return "e0" + strconv.Itoa(int(epNum))
}
