package lines

import (
	"log"
	"strconv"
	"strings"
)

func ParseCommaSeparatedStrings(line string) []string {
	var ss []string
	for _, s := range strings.Split(line, ",") {
		ss = append(ss, s)
	}
	return ss
}
func ParseCommaSeparatedInts(line string) []int {
	var ns []int
	for _, s := range strings.Split(line, ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		ns = append(ns, n)
	}
	return ns
}
