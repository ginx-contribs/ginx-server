package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show app version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion(Version)
	},
}

func printVersion(v string) {
	if v == "" {
		v = "unknown"
	}
	fmt.Printf("%s %s/%s\n", v, runtime.GOOS, runtime.GOARCH)
}
