package main

import (
	"fmt"
	"local-utils/netstat/statfiles"
)

// TODO: read all the inodes from the /proc/process_id/fd folders
// TODO: compare the inodes from the net files to the inodes of the fd files to get process id

func main() {
	stats := statfiles.ParseNetStats()

	for _, stat := range stats {
		fmt.Println(stat.Line, stat.LocalAddress, stat.RemoteAddress, stat.State, stat.Inode)
	}
}
