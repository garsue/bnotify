package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/garsue/bnotify"
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
	quote := flag.Bool("q", false, "Quote a message from remote or a file")
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
	comment, err := bnotify.NewComment(msg, quote, msgURL, fpath)
	if err != nil {
		log.Fatalln(err)
	}
	debugln(comment)
	data.Set("content", comment.String())
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
	debugln(string(body))
}

func debugln(v ...interface{}) {
	if verbose != nil && *verbose {
		log.Println(v...)
	}
}

func debugf(format string, v ...interface{}) {
	if verbose != nil && *verbose {
		log.Printf(format, v...)
	}
}
