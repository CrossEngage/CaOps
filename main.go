//go:generate make gen_version
package main

import "bitbucket.org/crossengage/athena/cmd"

func main() {
	cmd.Execute()
}
