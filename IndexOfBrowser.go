package IndexOfBrowser

import (
	"crypto/tls"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Browser struct {
	Url        string
	cwd        string
	httpClient *http.Client
	toPop      []string
}

type Entry struct {
	Directory bool
	Name      string
	Url       string
}

func NewBrowser(baseUrl string) *Browser {
	//string final /
	baseUrl = strings.TrimRight(baseUrl, "/")
	urlParts, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	rootUrl := strings.Replace(baseUrl, urlParts.Path, "", -1)
	return &Browser{
		Url: rootUrl,
		cwd: urlParts.Path,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxConnsPerHost:       2,
				ResponseHeaderTimeout: 2 * time.Second,
				ExpectContinueTimeout: 2 * time.Second,
			},
			Timeout: 5 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (browser *Browser) Cwd() string {
	return browser.cwd + "/"
}

func (browser *Browser) Ls() (entries []Entry, err error) {
	// Get current page
	req, err := http.NewRequest(http.MethodGet, browser.Url+browser.cwd+"/", nil)
	if err != nil {
		return entries, err
	}
	resp, err := browser.httpClient.Do(req)
	if err != nil {
		return entries, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Request.URL.String())
		return entries, errors.New("error getting page")
	}
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return entries, err
	}
	if !strings.HasPrefix(document.Find("title").Text(), "Index of /") {
		return entries, errors.New("not a directory listing")
	}
	//find all trs
	trs := document.Find("tr")
	trs.Each(func(i int, selection *goquery.Selection) {
		//skip 2
		if i <= 1 {
			return
		}
		if strings.Contains(selection.Text(), "Parent") {
			return
		}
		a := selection.Find("a")
		link, found := a.Attr("href")
		if !found || link == browser.cwd {
			return
		}
		//get tds
		entries = append(entries, Entry{
			Directory: strings.HasSuffix(link, "/"),
			Url:       browser.Url + browser.cwd + link,
			Name:      link,
		})
	})
	return entries, nil
}

func (browser *Browser) ChDir(dir string) {
	browser.cwd = strings.TrimRight(dir, "/")
}

func (browser *Browser) Pushd(dir string) {
	browser.toPop = append(browser.toPop, browser.cwd)
	browser.ChDir(dir)
}

func (browser *Browser) Popd() {
	if len(browser.toPop) > 0 {
		browser.cwd, browser.toPop = browser.toPop[len(browser.toPop)-1], browser.toPop[:len(browser.toPop)-1]
	}
}
