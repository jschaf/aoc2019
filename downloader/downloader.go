package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	client := new(http.Client)
	sessionCookie := ""
	year := 2019
	day := 01
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	checkError(err)

	cookie := new(http.Cookie)
	cookie.Name, cookie.Value = "session", sessionCookie
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Got non 200 status code: %d", resp.StatusCode)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
