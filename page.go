package gatherProxyParsing

import (
    "io/ioutil"
    "net/http"
    "time"
    "fmt"
    "strings"
    "net/url"
    "bytes"
)

const timeout time.Duration = 10000

func getUserAgent() (string) {
    return "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36"
}

func getPage(parseUrl string, proxyType string, listIndex int, minUptime uint32) ([]byte, error) {
    var client http.Client

    transport := http.Transport{}

    client = http.Client{
        Transport: &transport,
        Timeout:   time.Duration(time.Millisecond * timeout),
    }

    data := url.Values{}
    data.Add(typeParam, proxyType)
    data.Add(indexParam, fmt.Sprint(listIndex))
    data.Add(uptimeParam, fmt.Sprint(minUptime))

    request, err := http.NewRequest("POST", parseUrl, bytes.NewBufferString(data.Encode()))
    request.Header.Add("User-Agent", getUserAgent())
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    response, err := client.Do(request)

    if err != nil {
        return []byte(""), err
    }

    defer response.Body.Close()
    text, err := ioutil.ReadAll(response.Body)

    if err != nil {
        return []byte(""), err
    }

    return text, nil
}

func testPage(parseUrl string, proxy Proxy, testStr string) (bool, error) {
    var client http.Client

    transport := http.Transport{}

    proxyUrl, _ := url.Parse(fmt.Sprintf("%v://%v:%v", proxy.Type, proxy.Ip, proxy.Port))

    transport = http.Transport{
        Proxy: http.ProxyURL(proxyUrl),
    }

    client = http.Client{
        Transport: &transport,
        Timeout:   time.Duration(time.Millisecond * timeout),
    }

    request, err := http.NewRequest("GET", parseUrl, nil)
    request.Header.Add("User-Agent", getUserAgent())
    response, err := client.Do(request)

    if err != nil {
        return false, err
    }

    defer response.Body.Close()
    text, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return false, err
    }

    if testStr != "" && !strings.Contains(string(text), testStr) {
        return false, nil
    }

    return true, nil
}
