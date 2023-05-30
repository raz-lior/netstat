package statfiles

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)

const HEADER_LINE = 1

type NetworkStats struct {
	Line string
	LocalAddress string
	RemoteAddress string
	State string
	Inode string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseNetAddress(addr string) string {

	part1, err1 := strconv.ParseInt(addr[6:8], 16, 32)
	check(err1)

	part2, err2 := strconv.ParseInt(addr[4:6], 16, 32)
	check(err2)

	part3, err3 := strconv.ParseInt(addr[2:4], 16, 32)
	check(err3)

	part4, err4 := strconv.ParseInt(addr[0:2], 16, 32)
	check(err4)


	port, err := strconv.ParseInt(addr[9:13],16, 32)
	check(err)

	return fmt.Sprintf("%d.%d.%d.%d:%d", part1, part2, part3, part4, port)
}

func parseRowWord(pos int, word string) string {
	switch pos {
	case 1:
		return word
	case 2:
		return parseNetAddress(word)
	case 3:
		return parseNetAddress(word)
	case 4:
		return word
	case 10:
		return word
	default:
		return ""
	}
}

func parseHeaderWord(pos int, word string) string {
	switch pos {
	case 1:
		return word
	case 2:
		return word
	case 3:
		return word
	case 4:
		return word
	case 12:
		return word
	default:
		return ""
	}
}

func getParser(row int) func(int,string) string {

	var parser func (int, string) string

	if row == HEADER_LINE {
		parser = parseHeaderWord
	} else {
		parser = parseRowWord
	}

	return parser
}

func ParseNetStats() []NetworkStats {
	tcpData, err := os.Open("/proc/net/tcp")
	defer tcpData.Close()
	check(err)

	stats := make([]NetworkStats, 0)
	scanner := bufio.NewScanner(tcpData)
	scanner.Scan() // skipping the header

	for scanner.Scan() {

		line := scanner.Text()

		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords)

		lineWordCount := 0
		var stat NetworkStats
		row := make([]string, 0)
		for wordScanner.Scan() {
			lineWordCount++
			row = append( row, parseRowWord(lineWordCount, wordScanner.Text()) )
		}
		stat.Line = row[0]
		stat.LocalAddress = row[1]
		stat.RemoteAddress = row[2]
		stat.State = row[3]
		stat.Inode = row[9]

		stats = append(stats, stat)
	}

	if err1 := scanner.Err(); err1 != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err1)
	}

	return stats
}
