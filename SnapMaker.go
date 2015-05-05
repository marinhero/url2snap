/*
** SnapMaker.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"text/template"

	"github.com/marinhero/wootricChallenge/urlbox"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

//S3 bucket name
const BUCKETNAME = "snapshotswootric"

//Page provides an endpoint of valuable information
//needed by the fronted
type Page struct {
	Title    string
	S3Link   string
	Messages string
}

func validURL(url string) bool {
	r := regexp.MustCompile("(http(s*))://([\\w]*.)*")
	return r.MatchString(url)
}

func validForm(f url.Values) (data urlbox.ShotData, strErr string) {
	rawURL := f.Get("url")
	width, _ := strconv.ParseUint(f.Get("width"), 10, 0)
	height, _ := strconv.ParseUint(f.Get("height"), 10, 0)
	if validURL(rawURL) != true {
		strErr = fmt.Sprintf("Invalid format encountered in URL")
		return
	}
	if width > 0 &&
		height > 0 {
		u, _ := url.Parse(rawURL)
		data = urlbox.ShotData{u.Host, uint(width), uint(height)}
		strErr = ""
		return
	}
	strErr = fmt.Sprintf("Invalid width|height values")
	return
}

//GetFile opens the file in local disk and returns its content
func getFile(fileName string) (bytes []byte, filetype string) {
	file, _ := os.Open(fileName)
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	bytes = make([]byte, size)

	buffer := bufio.NewReader(file)
	buffer.Read(bytes)

	filetype = http.DetectContentType(bytes)
	return
}

//UploadToS3 sends the snapshot to S3 storage and returns a public link
func uploadToS3(data urlbox.ShotData) (s3Link string) {
	if data.URL != "" && data.Width > 0 && data.Height > 0 {
		auth, _ := aws.EnvAuth()
		client := s3.New(auth, aws.USEast)
		demoBucket := client.Bucket(BUCKETNAME)

		fileName := urlbox.GetFileName(data)
		bytes, filetype := getFile(fileName)
		demoBucket.Put(fileName, bytes, filetype, s3.PublicRead)
		s3Link = demoBucket.URL(fileName)
		os.Remove(fileName)
	}
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "index.html")
	index := Page{Title: "SnapMaker - By Marin Alcaraz"}

	t, err := template.ParseFiles(lp)
	if err != nil {
		log.Fatal("[!]indexHandler:", err)
	}
	if r.Method == "POST" {
		r.ParseForm()
		data, strErr := validForm(r.Form)
		if strErr != "" {
			index.Messages = strErr
		} else {
			if urlbox.CreateShot(data) == "OK" {
				index.S3Link = uploadToS3(data)
			} else {
				index.Messages = "API returned KO"
			}
		}
	}
	t.Execute(w, index)
}

func initWebServer() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func main() {
	initWebServer()
}
