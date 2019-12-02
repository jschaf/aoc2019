package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	client := new(http.Client)
	sessionCookie := ""
	year := 2019
	if len(os.Args) < 2 {
		log.Fatal("Need at least 1 arg, the day")
	}
	day, err := strconv.Atoi(os.Args[1])
	checkError(err)
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	req, err := http.NewRequest("GET", url, nil)
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

	dayPrefix := strconv.Itoa(day)
	if day < 10 {
		dayPrefix = fmt.Sprintf("0%d", day)
	}
	path := fmt.Sprintf("day%s/day%s.go", dayPrefix, dayPrefix)
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0o666)
	checkError(err)
	defer file.Close()
	file.WriteString("package main\n\n")

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		fmt.Fprintf(file, "// %s\n", s.Text())
	}
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

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
