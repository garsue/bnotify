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

func (u *users) String() string {
	return fmt.Sprint(*u)
}

func (u *users) Set(value string) error {
	*u = append(*u, value)
	return nil
}

func main() {
	var users users
	flag.Var(&users, "u", "Set a notified user ID")
	fpath := flag.String("f", "", "Set a comment file")
	flag.Parse()
	if fpath == nil || len(*fpath) == 0 {
		log.Fatalln("Must set a comment file")
	}
	domain := os.Getenv("BACKLOG_DOMAIN")
	spaceID := os.Getenv("BACKLOG_SPACE_ID")
	issue := os.Getenv("BACKLOG_ISSUE")
	apiKey := os.Getenv("BACKLOG_API_KEY")
	endpoint := fmt.Sprintf(
		"https://%s.%s/api/v2/issues/%s/comments?apiKey=%s",
		spaceID,
		domain,
		issue,
		apiKey,
	)
	data := url.Values{}
	comment, err := loadComment(*fpath)
	log.Println(comment)
	if err != nil {
		log.Fatalln(err)
	}
	data.Set("content", comment)
	for _, user := range users {
		data.Add("notifiedUserId[]", user)
	}

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
	fmt.Print(string(body))
}

func loadComment(fpath string) (string, error) {
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
