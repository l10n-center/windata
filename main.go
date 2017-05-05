package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const (
	url = "https://www.microsoft.com/resources/msdn/goglobal/default.mspx"
)

type row struct {
	LCID int64  `json:"lcid"`
	Code string `json:"code"`
	Name string `json:"name"`
	ANSI int64  `json:"ansi"`
}

func main() {
	rows, err := parsePage(url)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
	fmt.Print("Save as file? (y/N) ")
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		answer = "n"
	}
	answer = strings.ToLower(answer)
	if answer == "y" || answer == "yes" {
		fmt.Print("Enter file name (default: windata): ")
		if _, err := fmt.Scanln(&answer); err != nil {
			answer = "windata"
		}
		if !strings.HasSuffix(answer, ".json") {
			answer += ".json"
		}
		f, err := os.Create(answer)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		defer f.Close()

		e := json.NewEncoder(f)
		e.SetIndent("", "   ")
		if err := e.Encode(rows); err != nil {
			fmt.Printf("ERROR: %s", err)
			os.Exit(1)
		}
	}
}

func parsePage(url string) ([]row, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	root := html.NewTokenizer(resp.Body)
	rows := []row{}
	r := row{}
	state := ""
	for {
		t := root.Next()
		switch t {
		case html.ErrorToken:

			return nil, root.Err()

		case html.StartTagToken:
			tn, _ := root.TagName()
			switch string(tn) {
			case "th":
				state = "th"
			case "td":
				switch state {
				case "":
					state = "lcid"
				case "lcid":
					state = "code"
				case "code":
					state = "name"
				case "name":
					state = "lang"
				case "lang":
					state = "local"
				case "local":
					state = "ansi"
				case "ansi":
					state = "oem"
				case "oem":
					state = "cntr"
				case "cntr":
					state = "lng"
				}
			}
		case html.EndTagToken:
			tn, _ := root.TagName()
			switch string(tn) {
			case "tr":
				switch state {
				case "th":
					state = ""
				case "lng"
				}
			}
		}
	}
	return []row{}, nil
}
