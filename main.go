package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isCharDevice(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	fileMode := fileInfo.Mode()
	return (fileMode & os.ModeCharDevice) != 0
}

func buildBalloon(lines []string, maxwidth int) string {
	var ret []string

	top := " " + strings.Repeat("_", maxwidth+2)
	ret = append(ret, top)

	borders := []string{"/", "\\", "\\", "/", "|", "<", ">"}
	for index, line := range lines {
		if len(lines) == 1 {
			ret = append(ret, fmt.Sprintf("%s %s %s", borders[5], line, borders[6]))
			continue
		}

		if index == 0 {
			ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[0], line, borders[1]))
		} else if index == len(lines)-1 {
			ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[2], line, borders[3]))
		} else {
			ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[4], line, borders[4]))
		}
	}

	bottom := " " + strings.Repeat("-", maxwidth+2)
	ret = append(ret, bottom)

	return strings.Join(ret, "\n")
}

func calculateMaxWidth(lines []string) int {
	var max int
	for _, line := range lines {
		length := utf8.RuneCountInString(line)
		if length > max {
			max = length
		}
	}
	return max
}

func tabsToSpaces(lines []string) {
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.Replace(lines[i], "\t", "    ", -1)
	}
}

func normalizeStringsLength(lines []string, maxwidth int) {
	for i := 0; i < len(lines); i++ {
		lines[i] = lines[i] + strings.Repeat(" ", maxwidth-utf8.RuneCountInString(lines[i]))

	}
}

func printFigure(name string) {
	var cow = `         \  ^__^
          \ (oo)\_______
	    (__)\       )\/\
	        ||----w |
	        ||     ||
		`
	var stegosaurus = `         \                      .       .
          \                    / ` + "`" + `.   .' "
           \           .---.  <    > <    >  .---.
            \          |    \  \ - ~ ~ - /  /    |
          _____           ..-~             ~-..-~
         |     |   \~~~\\.'                    ` + "`" + `./~~~/
        ---------   \__/                         \__/
       .'  O    \     /               /       \  "
      (_____,    ` + "`" + `._.'               |         }  \/~~~/
       ` + "`" + `----.          /       }     |        /    \__/
             ` + "`" + `-.      |       /      |       /      ` + "`" + `. ,~~|
                 ~-.__|      /_ - ~ ^|      /- _      ` + "`" + `..-'
                      |     /        |     /     ~-.     ` + "`" + `-. _  _  _
                      |_____|        |_____|         ~ - . _ _ _ _ _>

	`

	switch name {
	case "cow":
		fmt.Println(cow)
	case "stegosaurus":
		fmt.Println(stegosaurus)
	default:
		fmt.Println("Unknown figure")
	}
}

func main() {
	var figure string
	flag.StringVar(&figure, "f", "cow", "the figure name. Valid values are `cow` and `stegosaurus`")
	flag.Parse()

	if isCharDevice(os.Stdin) {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gocowsay")
		return
	}

	var lines []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())

	tabsToSpaces(lines)
	maxwidth := calculateMaxWidth(lines)
	normalizeStringsLength(lines, maxwidth)
	balloon := buildBalloon(lines, maxwidth)

	fmt.Println(balloon)
	printFigure(figure)
	fmt.Println()
}
