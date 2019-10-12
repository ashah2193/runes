package main

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"strings"

	"golang.org/x/text/unicode/runenames"
)

type CharName struct {
	Char rune
	Name string
}

func (cn CharName) String() string {
	return fmt.Sprintf("%U\t%c\t%v", cn.Char, cn.Char, cn.Name)
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func containsAll(haystack, needles []string) bool {
	for _, s := range needles {
		if !contains(haystack, s) {
			return false
		}
	}
	return true
}

func scan(start, end rune) []CharName {
	result := []CharName{}
	for r := start; r < end; r++ {
		name := runenames.Name(r)
		hasName := !strings.HasPrefix(name, "<")
		if len(name) > 0 && hasName {
			result = append(result, CharName{r, name})
		}
	}
	return result
}

func filter(sample []CharName, words []string) []CharName {
	result := []CharName{}
	for i, s := range words {
		words[i] = strings.ToUpper(s)
	}
	for _, cn := range sample {
		name := strings.Replace(cn.Name, "-", " ", -1)
		parts := strings.Fields(name)
		if containsAll(parts, words) {
			result = append(result, CharName{cn.Char, cn.Name})
		}
	}
	return result
}

func report(words ...string) {
	result := filter(scan(' ', unicode.MaxRune), words)
	for _, r := range result {
		fmt.Println(fmt.Sprint(r))
	}
	fmt.Printf("%d character found", len(result))
}

func main() {
	if len(os.Args) > 1 {
		start := time.Now()
		report(os.Args[1:]...)
		elapsed := time.Since(start)
		fmt.Println("Program took %s", elapsed)
	} else {
		fmt.Println("Please provide at leat one word")
	}
}
