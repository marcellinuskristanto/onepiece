package helper

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

// Urljoin join url path
func Urljoin(sections ...string) string {
	u, err := url.Parse(sections[0])
	if err != nil {
		log.Fatal(err)
	}
	for _, section := range sections[1:] {
		u.Path = path.Join(u.Path, section)
	}
	return u.String()
}

// DownloadFile from url
func DownloadFileAndReturn(siteurl string) (*os.File, string, error) {
	response, err := http.Get(siteurl)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		defer os.Remove(file.Name())
		defer file.Close()
		return nil, "", err
	}
	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		defer os.Remove(file.Name())
		defer file.Close()
		return nil, "", err
	}

	return file, response.Header.Get("Content-type"), nil
}

// GetEnv default value
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvInt integer
func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		return strconv.Itoa(value)
	}
	return fallback
}
