package clients

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type QBitTorrentSettings struct {
	Host     string
	Port     string
	Username string
	Password string
	UseSSL   bool
}

type QBitTorrentClient struct {
	Settings     *QBitTorrentSettings
	cookies      []*http.Cookie
	isAuthorized bool
}

func (c *QBitTorrentClient) SendMagnet(magnetLink string) error {
	err := c.authQBitTorrent()
	if err != nil {
		fmt.Printf("Error authorizing on QBitTorrent: %s", err.Error())
		return fmt.Errorf("Error authorizing on QBitTorrent")
	}

	qBitUrl := c.getUrl() + "/api/v2/torrents/add"

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(qBitUrl)
	jar.SetCookies(u, c.cookies)

	client := &http.Client{Jar: jar}
	formData := url.Values{
		"urls":     {magnetLink},
		"category": {"crawler"},
	}

	resp, err := client.PostForm(qBitUrl, formData)

	if err != nil {
		fmt.Printf("Error making request to %s: %s\n", qBitUrl, err.Error())
		return fmt.Errorf("Error making request to %s: %s\n", qBitUrl, err.Error())
	}

	defer resp.Body.Close()

	fmt.Printf("Received status code %d from %s\n", resp.StatusCode, qBitUrl)

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return fmt.Errorf("Error reading response body: %s\n", err.Error())
	}

	fmt.Printf("QBitTorrent Response: (%v) %s \n", n, buf.String())

	return nil
}

func (c *QBitTorrentClient) SetSettings(host string, port string, username string, password string, useSSL bool) {
	c.Settings = &QBitTorrentSettings{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		UseSSL:   useSSL,
	}
}

func (c *QBitTorrentClient) getUrl() string {
	return "http://192.168.0.182:8081"
}

func (c *QBitTorrentClient) authQBitTorrent() error {
	fmt.Println("Iniciando AuthQBitTorrent")
	qBitUrl := c.getUrl() + "/api/v2/auth/login"

	formData := url.Values{
		"username": {c.Settings.Username},
		"password": {c.Settings.Password},
	}

	resp, err := http.PostForm(qBitUrl, formData)
	if err != nil {
		fmt.Printf("Error making auth request to %s: %s\n", qBitUrl, err.Error())
		return fmt.Errorf("Error making auth request to %s: %s\n", qBitUrl, err.Error())
	}

	defer resp.Body.Close()

	fmt.Printf("Auth Received status code %d from %s\n", resp.StatusCode, qBitUrl)

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Printf("Error reading auth response body: %s\n", err.Error())
		return fmt.Errorf("Error reading auth response body: %s\n", err.Error())
	}

	fmt.Printf("QBitTorrent AUTH Response: (%v) %s \n", n, buf.String())

	c.cookies = resp.Cookies()
	return nil
}
