package main

import "fmt"

// TODO: read all the inodes from the /proc/process_id/fd folders
// TODO: compare the inodes from the net files to the inodes of the fd files to get process id

func main() {
	stats := ParseNetStats()

	for _, stat := range stats {
		fmt.Println(stat.line, stat.localAddress, stat.remoteAddress, stat.state, stat.inode)
	}
}
