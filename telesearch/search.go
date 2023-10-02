package telesearch

import (
	json "encoding/json"
	http "net/http"
    os "os"
    fmt "fmt"
    // log "log"
    ioutil "io/ioutil"
    strings "strings"
)

const (
    searchUrl = "https://www.googleapis.com/customsearch/v1"
)

type SearchItem struct {
    Title string `json:"title"`
    Description string `json:"snippet"`
    Url string `json:"link"`
}

type SearchResult struct {
    Success bool
    Results []SearchItem `json:"items"`
    Error string
}


func SearchGoogle(query string, pageOffset int, channel chan SearchResult) {
    googleApiKey := os.Getenv("GOOGLE_API_KEY")
    googleSearchId := os.Getenv("GOOGLE_SEARCH_ID")

    query = strings.Replace(query, " ", "%20", -1)

    queryUrl := fmt.Sprintf(
        "%s?key=%s&cx=%s&start=%d&q=%s",
        searchUrl,
        googleApiKey,
        googleSearchId,
        pageOffset,
        query,
    )

    resp, err := http.Get(queryUrl)

    if err != nil {
        fmt.Printf("Error fetching %s: %s", queryUrl, err)
        
        channel <- SearchResult{Error: fmt.Sprint(err)}
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    result := SearchResult{}
    err = json.Unmarshal([]byte(body), &result)

    if err != nil {
        channel <- SearchResult{Error: fmt.Sprint(err)}
        return
    }
    result.Success = true
    channel <- result
}