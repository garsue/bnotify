package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type users []string

var (
	verbose *bool
)

func (u *users) String() string {
	return fmt.Sprint(*u)
}

func (u *users) Set(value string) error {
	*u = append(*u, value)
	return nil
}

func main() {
	// Load command line options
	var users users
	flag.Var(&users, "u", "Specifies a notified user ID")
	msg := flag.String("m", "", "Sets a message string")
	msgURL := flag.String("r", "", "Specifies a URL to load a remote message")
	fpath := flag.String("f", "", "Specifies a message file")
	verbose = flag.Bool("v", false, "Verbose mode")
	flag.Parse()

	// Get setting from env vars
	domain := os.Getenv("BACKLOG_DOMAIN")
	spaceID := os.Getenv("BACKLOG_SPACE_ID")
	issue := os.Getenv("BACKLOG_ISSUE")
	apiKey := os.Getenv("BACKLOG_API_KEY")

	// Compose a request
	endpoint := fmt.Sprintf(
		"https://%s.%s/api/v2/issues/%s/comments?apiKey=%s",
		spaceID,
		domain,
		issue,
		apiKey,
	)
	data := url.Values{}
	comment, err := makeCommentText(msg, msgURL, fpath)
	if err != nil {
		log.Fatalln(err)
	}
	debug(comment)
	data.Set("content", comment)
	for _, user := range users {
		data.Add("notifiedUserId[]", user)
	}

	// Post
	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusCreated {
		log.Fatalln(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	debug(string(body))
}

func makeCommentText(msg, msgURL, fpath *string) (string, error) {
	var comment string
	if msg != nil {
		comment = *msg
	}
	fromFile, err := loadFromFile(fpath)
	if err != nil {
		return "", err
	}
	if len(comment) > 0 {
		comment += "\n\n"
	}
	comment += fromFile
	fromURL, err := loadFromURL(msgURL)
	if err != nil {
		return "", err
	}
	if len(comment) > 0 {
		comment += "\n\n"
	}
	comment += fromURL
	return comment, nil
}

func loadFromFile(fpath *string) (string, error) {
	if fpath == nil {
		return "", nil
	}
	fp := *fpath
	if len(fp) == 0 {
		return "", nil
	}
	bytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func loadFromURL(msgURL *string) (string, error) {
	if msgURL == nil {
		return "", nil
	}
	urlStr := *msgURL
	if len(urlStr) == 0 {
		return "", nil
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func debug(format string, v ...interface{}) {
	if verbose != nil && *verbose {
		if len(v) > 0 {
			log.Printf(format, v)
		} else {
			log.Println(format)
		}
	}
}
