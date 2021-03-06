package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// def-matcher.go
// author: Seong Yong-ju <sei40kr@gmail.com>

type Def struct {
	alias string
	abbr  string
}

func ParseDef(line string) Def {
	alias := make([]rune, 0, 1024)
	abbr := make([]rune, 0, 1024)

	afterEscape := false
	inQuote := false
	inRightExp := false
	for _, aRune := range line {
		if aRune == '\\' {
			afterEscape = !afterEscape

			if afterEscape {
				continue
			}
		}

		if aRune == '\'' && !afterEscape {
			inQuote = !inQuote
		} else if aRune == '=' && !inQuote {
			inRightExp = true
		} else if !inRightExp {
			alias = append(alias, aRune)
		} else {
			abbr = append(abbr, aRune)
		}

		afterEscape = false
	}

	return Def{alias: string(alias), abbr: string(abbr)}
}

func MatchDef(defs []Def, command string) (*Def, bool) {
	sort.Slice(defs, func(i, j int) bool {
		return len(defs[j].abbr) <= len(defs[i].abbr)
	})

	var candidate Def
	isFullMatch := false

	for true {
		var match Def
		for _, def := range defs {

			if command == def.abbr {
				match = def
				isFullMatch = true
				break
			} else if strings.HasPrefix(command, def.abbr) {
				match = def
				break
			}
		}

		if match != (Def{}) {
			command = fmt.Sprintf("%s%s", match.alias, command[len(match.abbr):])
			candidate = match
		} else {
			break
		}
	}

	if candidate != (Def{}) {
		return &candidate, isFullMatch
	} else {
		return nil, false
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid number of arguments")
		os.Exit(1)
	}

	defs := make([]Def, 0, 512)

	scanner := bufio.NewScanner(bufio.NewReaderSize(os.Stdin, 1024))
	for scanner.Scan() {
		line := scanner.Text()
		defs = append(defs, ParseDef(line))
	}

	command := os.Args[1]
	match, isFullMatch := MatchDef(defs, command)
	if match != nil {
		if isFullMatch {
			fmt.Printf("%s\n", match.alias)
		} else {
			fmt.Printf("%s%s\n", match.alias, command[len(match.abbr):])
		}
	}
}
