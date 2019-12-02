package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	client := new(http.Client)
	year := 2019
	if len(os.Args) < 2 {
		log.Fatal("Need at least 1 arg, the day")
	}
	day, err := strconv.Atoi(os.Args[1])
	checkError(err)
	dayPrefix := strconv.Itoa(day)
	if day < 10 {
		dayPrefix = fmt.Sprintf("0%d", day)
	}

	inputUrl := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	input := getUrl(client, inputUrl)
	inputBytes, err := ioutil.ReadAll(input)
	checkError(err)
	inputPath := fmt.Sprintf("day%s/input.txt", dayPrefix)
	err = os.MkdirAll(filepath.Dir(inputPath), 0o666)
	checkError(err)
	inputFile, err := os.OpenFile(inputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	checkError(err)
	defer inputFile.Close()
	inputFile.Write(inputBytes)

	goPath := fmt.Sprintf("day%s/day%s.go", dayPrefix, dayPrefix)
	file, err := os.OpenFile(
		goPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0o666)
	checkError(err)
	defer file.Close()
	file.WriteString("package main\n\n")

	file.WriteString(`
func main() {
	part1()
	part2()
}
`)
	fmt.Fprintf(file, `
func part1() {
	ints, err := files.ReadAllLinesAsInts("%s/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	for _, n := range ints {
		total += doPart1(n)
	}
	fmt.Printf("Part 1: %%d\n", total)
}

func doPart1(n int) int {
	return 1
}
`, dayPrefix)

	fmt.Fprintf(file, `
func part2() {
	ints, err := files.ReadAllLinesAsInts("%s/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	for _, n := range ints {
		total += doPart2(n)
	}
	fmt.Printf("Part 2: %%d\n", total)
}

func doPart2(n int) int {
	return 1
}
`, dayPrefix)

}

func getUrl(client *http.Client, url string) io.ReadCloser {
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	cookie := createCookie()
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Got non 200 status code: %d for url %s", resp.StatusCode, url)
	}
	return resp.Body
}

func createCookie() *http.Cookie {
	path := "aoc_session.txt"
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open session value from %s", path)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read session file from %s", path)
	}
	return &http.Cookie{
		Name:  "session",
		Value: string(bytes),
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
