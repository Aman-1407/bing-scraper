package main

import(
	"fmt"
	"net/http"
	"strings"
	"time"
	"math/rand"
	"net/url"
	"github.com/PuerkitoBio/goquery"
)

func randomUserAgent() string{
rand.Seed(time.Now().Unix())
randNum := rand.Int()%len(userAgents)
return userAgents[randNum]
}

func buildBingUrls(searchTerm, country string, pages, count int)([]string, error){
toScrape := []string{}
searchTerm = strings.Trim(searchTerm, " ")
searchTerm = strings.Replace(searchTerm, " ", "+", -1)
if countryCode, found := bingDomains[country]; found{
	for i :=0; i<pages ; i++{
		first := firstParameter(i,count);
		scrapeURL := fmt.Sprintf("https://bing.com/search?q=%s&first=%d&count=%d%s",searchTerm, first, count,countryCode)
		toScrape = append(toScrape, scrapeURL)
	}
}else{
	err := fmt.Errorf("country(%s)is currently not supported", country)
	return nil, err
}
return toScrape, nil		
}

func firstParameter(number, count int) int{
if number == 0{
	return number +1
}
return number*count +1
}


func getScrapeClient(proxyString interface{}) *http.Client{
switch V:=proxyString.(type){
case string:
	proxyUrl, _ := url.Parse(V)
	return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
default:
	return &http.Client{}
	}
}

func scrapeClientRequest(searchURL string, proxyString interface{})(*http.Response, error){

baseClient := getScrapeClient(proxyString)
req, _ := http.NewRequest("GET", searchURL, nil)
req.Header.Set("User-Agent", randomUserAgent())

res, err := baseClient.Do(req)
if res.StatusCode !=200 {
	err := fmt.Errorf("scraper received a non-200 status code suggesting a ban")
	return nil, err
}

if err != nil{
	return nil, err
}
return res, nil
}

func bingResultParser(response *http.Response, rank int)([]SearchResult, error){

doc, err := goquery.NewDocumentFromResponse(response)
if err != nil{
return nil, err
}
results := []SearchResult{}
sel := doc.Find("li.b_algo")
rank++

for i := range sel.Nodes{
	item := sel.Eq(i)
	linkTag := item.Find("a")
	link, _ := linkTag.Attr("href")
	titleTag := item.Find("h2")
	descTag := item.Find("div.b_caption p")
	desc := descTag.Text()
	title := titleTag.Text()
	link = strings.Trim(link, " ")
	if link != "" && link != "#" && !strings.HasPrefix(link, "/"){
		result := SearchResult{
			rank,
			link,
			title,
			desc,
			time.Now(),
		}
		results = append(results, result)
		rank++
	}
}
return results, err
}


func BingScrape(searchTerm, country string, proxyString interface{}, pages, count, backoff int)([]SearchResult, error){
results := []SearchResult{}

bingPages, err := buildBingUrls(searchTerm, country, pages, count)

if err!= nil{
	return nil, err
}

for _, page :=range bingPages{

	rank := len(results)
	res, err := scrapeClientRequest(page, proxyString)
	if err!=nil{
	return nil, err
	}
	data, err := bingResultParser(res, rank)
	if err != nil{
		return nil, err
	}
	for _, result := range data{
		results = append(results, result)
	}
	time.Sleep(time.Duration(backoff)*time.Second)
}
return results, nil
}

func main(){
create_database()
}