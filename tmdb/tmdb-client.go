package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getTmdbToken() string {
	return "[TOKEN]"
}

func GetEpisodes(series_id int32, season_number int32) ([]TmdbEpisode, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%v/season/%v", series_id, season_number)
	fmt.Printf("GetEpisodes: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request to %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Error creating request to %s: %s\n", url, err.Error())
	}
	client := http.Client{}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getTmdbToken()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request to %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Error executing request to %s: %s\n", url, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from %s", resp.StatusCode, url)
		return nil, fmt.Errorf("Error: received status code %d from %s", resp.StatusCode, url)
	}

	var decodedResponse TmdbSeasonResponse
	err = json.NewDecoder(resp.Body).Decode(&decodedResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
		return nil, fmt.Errorf("Error decoding response: %s", err.Error())
	}
	fmt.Printf("Decoded response len= %v\n", len(decodedResponse.Episodes))

	return decodedResponse.Episodes, nil
}

func GetSeries(series_name string) ([]TmdbSeries, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/tv?query=%s&include_adult=false&language=en-US&page=1", strings.Replace(series_name, " ", "%", -1))
	fmt.Printf("GetSeries: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request to %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Error creating request to %s: %s\n", url, err.Error())
	}
	client := http.Client{}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getTmdbToken()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request to %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Error executing request to %s: %s\n", url, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from %s", resp.StatusCode, url)
		return nil, fmt.Errorf("Error: received status code %d from %s", resp.StatusCode, url)
	}

	var decodedResponse TmdbResponse
	err = json.NewDecoder(resp.Body).Decode(&decodedResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
		return nil, fmt.Errorf("Error decoding response: %s", err.Error())
	}
	fmt.Printf("Decoded response len= %v\n", len(decodedResponse.Results))

	return decodedResponse.Results, nil
}

type TmdbResponse struct {
	Page          int32
	Results       []TmdbSeries
	Total_pages   int32
	Total_results int32
}

type TmdbSeries struct {
	Adult             bool
	Backdrop_path     string
	Genre_ids         []int32
	Id                int32
	Origin_country    []string
	Original_language string
	Original_name     string
	Overview          string
	Popularity        float32
	Poster_path       string
	First_air_date    string
	Name              string
	Vote_average      float32
	Vote_count        int32
}

type TmdbSeasonResponse struct {
	ID           string        `json:"_id"`
	AirDate      string        `json:"air_date"`
	Episodes     []TmdbEpisode `json:"episodes"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
	VoteAverage  float32       `json:"vote_average"`
}

type TmdbEpisode struct {
	AirDate        string       `json:"air_date"`
	EpisodeNumber  int          `json:"episode_number"`
	EpisodeType    string       `json:"episode_type"`
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Overview       string       `json:"overview"`
	ProductionCode string       `json:"production_code"`
	Runtime        int          `json:"runtime"`
	SeasonNumber   int          `json:"season_number"`
	ShowID         int          `json:"show_id"`
	StillPath      string       `json:"still_path"`
	VoteAverage    float32      `json:"vote_average"`
	VoteCount      int          `json:"vote_count"`
	Crew           []CrewMember `json:"crew"`
	GuestStars     []CastMember `json:"guest_stars"`
}

type CrewMember struct {
	Job                string  `json:"job"`
	Department         string  `json:"department"`
	CreditID           string  `json:"credit_id"`
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float32 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
}

type CastMember struct {
	Character          string  `json:"character"`
	CreditID           string  `json:"credit_id"`
	Order              int     `json:"order"`
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float32 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
}
