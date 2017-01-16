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
	flag.Var(&users, "u", "Set a notified user ID")
	msg := flag.String("m", "", "Set a message")
	fpath := flag.String("f", "", "Set a comment file")
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
	comment, err := makeCommentText(msg, fpath)
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

func makeCommentText(msg, fpath *string) (string, error) {
	var comment string
	if msg != nil {
		comment = *msg
	}
	if fpath == nil {
		return comment, nil
	}
	bytes, err := ioutil.ReadFile(*fpath)
	if err != nil {
		return "", err
	}
	return comment + "\n" + string(bytes), nil
}

func debug(format string, v ...interface{}) {
	if verbose != nil && *verbose {
		log.Printf(format, v)
	}
}
