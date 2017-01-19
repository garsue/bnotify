package bnotify

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Comment has comment data
type Comment struct {
	message           string
	quote             bool
	remoteMessageURL  string
	messageFromRemote string
	messageFilePath   string
	messageFromFile   string
}

// NewComment returns new Comment instance
func NewComment(msg *string, quote *bool, msgURL, fpath *string) (Comment, error) {
	c := Comment{
		message:          *msg,
		quote:            *quote,
		remoteMessageURL: *msgURL,
		messageFilePath:  *fpath,
	}
	err := c.loadFromFile()
	if err != nil {
		return Comment{}, err
	}
	err = c.loadFromURL()
	if err != nil {
		return Comment{}, err
	}
	return c, nil
}

func (c Comment) String() string {
	msg := c.message
	if len(c.messageFromFile) > 0 {
		if len(msg) > 0 {
			msg += "\n\n"
		}
		if c.quote {
			msg += "```\n"
		}
		msg += c.messageFromFile
		if c.quote {
			msg += "```"
		}
	}
	if len(c.messageFromRemote) > 0 {
		if len(msg) > 0 {
			msg += "\n\n"
		}
		if c.quote {
			msg += "```\n"
		}
		msg += c.messageFromRemote
		if c.quote {
			msg += "```"
		}
	}
	return msg
}

func (c *Comment) loadFromFile() error {
	if len(c.messageFilePath) == 0 {
		return nil
	}
	bytes, err := ioutil.ReadFile(c.messageFilePath)
	if err != nil {
		return err
	}
	c.messageFromFile = string(bytes)
	return nil
}

func (c *Comment) loadFromURL() error {
	if len(c.remoteMessageURL) == 0 {
		return nil
	}
	resp, err := http.Get(c.remoteMessageURL)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.messageFromRemote = string(bytes)
	return nil
}
