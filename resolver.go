package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	log "github.com/sirupsen/logrus"
)

var client *http.Client

type editorJsLinkImage struct {
	URL string `json:"url"`
}

type editorJsLinkMeta struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       editorJsLinkImage `json:"image"`
}

type editorJsLinkInfo struct {
	Success int              `json:"success"`
	Meta    editorJsLinkMeta `json:"meta"`
}

const urlRegex = `(http|ftp|https)://([\w+?\.\w+])+([a-zA-Z0-9\~\!\@\#\$\%\^\&\*\(\)_\-\=\+\\\/\?\.\:\;\'\,]*)?`

func newClient() *http.Client {
	if client != nil {
		return client
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr}
	return client
}

func getBody(url string) string {
	res, err := newClient().Get(url)
	if err != nil {
		log.Errorf("Error on getting by address: '%v'", url)
		log.Error(err)
	}
	defer res.Body.Close()
	buf := bytes.NewBuffer([]byte{})
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		log.Error("Error on reading of response body")
	}
	return buf.String()
}

func getTargetFromRequest(r *http.Request) (string, error) {
	qsValues := r.URL.Query()
	if val, ok := qsValues["url"]; ok {
		return val[0], nil
	}
	return "", errors.New("cannot to get target url from request")
}

func checkURL(url string) bool {
	ok, err := regexp.MatchString(urlRegex, url)
	if err != nil {
		log.Errorf("Cannot to check url")
	}
	return ok
}

func getImageURL(og *opengraph.OpenGraph) string {
	if len(og.Images) != 0 {
		return og.Images[0].URL
	}
	return ""
}

func ogToJSON(og *opengraph.OpenGraph) string {
	info := editorJsLinkInfo{
		Success: 1,
		Meta: editorJsLinkMeta{
			Title:       og.Title,
			Description: og.Description,
			Image: editorJsLinkImage{
				URL: getImageURL(og),
			},
		},
	}
	json, err := json.Marshal(info)
	if err != nil {
		log.Errorf("Cannot to parse response to json")
	}
	return string(json)
}

func ogHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("New request for '%v'", r.RequestURI)
	targetPath, err := getTargetFromRequest(r)
	if err != nil {
		log.Error("targetPath is not defined in incomming request")
		return
	}
	if !checkURL(targetPath) {
		log.Errorf("Invalid target url: '%v'", targetPath)
		return
	}
	// get body and parse OG
	html := getBody(targetPath)
	og := opengraph.NewOpenGraph()
	err = og.ProcessHTML(strings.NewReader(html))
	if err != nil {
		log.Debug("Error on processing html")
	}
	// add header
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOW_ORIGIN"))
	fmt.Fprint(w, ogToJSON(og))
}

func main() {
	log.SetLevel(log.DebugLevel)
	http.HandleFunc("/fetchUrl", ogHandler)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	log.Infof("starting of server on host:%v port:%v", host, port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%v:%v", host, port),
		nil,
	))
}
