package statfiles

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)

const HEADER_LINE = 1

var TCP_STATE_CODE_MAP = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
	"0C": "NEW_SYN_RECV",
	"0D": "MAX_STATES",
}

type NetworkStats struct {
	Protocol string
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

func parseTcpFile(protocol string) []NetworkStats {

	filePath := "/proc/net/" + protocol

	tcpData, err := os.Open(filePath)
	defer tcpData.Close()
	check(err)

	stats := make([]NetworkStats, 0)
	lineScanner := bufio.NewScanner(tcpData)
	lineScanner.Scan() // skipping the header

	for lineScanner.Scan() {

		statParts := strings.Fields(lineScanner.Text())

		var stat NetworkStats

		stat.Protocol = protocol
		stat.LocalAddress = parseNetAddress(statParts[1])
		stat.RemoteAddress = parseNetAddress(statParts[2])
		stat.State = TCP_STATE_CODE_MAP[statParts[3]]
		stat.Inode = statParts[9]

		stats = append(stats, stat)
	}

	if err1 := lineScanner.Err(); err1 != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err1)
	}

	return stats
}

func ParseNetStats() []NetworkStats {

	stats := make([]NetworkStats,0)
	tcp := parseTcpFile("tcp")

	stats = append(stats, tcp...)
	return stats
}
