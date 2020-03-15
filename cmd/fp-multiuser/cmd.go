package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gofrp/fp-multiuser/pkg/server"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var (
	showVersion bool

	bindAddr  string
	tokenFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "l", "127.0.0.1:7200", "bind address")
	rootCmd.PersistentFlags().StringVarP(&tokenFile, "token_file", "f", "./tokens", "token file")
}

var rootCmd = &cobra.Command{
	Use:   "fp-multiuser",
	Short: "fp-multiuser is the server plugin of frp to support multiple users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version)
			return nil
		}
		tokens, err := ParseTokensFromFile(tokenFile)
		if err != nil {
			log.Printf("parse tokens from file %s error: %v", tokenFile, err)
			return err
		}
		s, err := server.New(server.Config{
			BindAddress: bindAddr,
			Tokens:      tokens,
		})
		if err != nil {
			return err
		}
		s.Run()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ParseTokensFromFile(file string) (map[string]string, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	rows := strings.Split(string(buf), "\n")
	for _, row := range rows {
		kvs := strings.SplitN(row, "=", 2)
		if len(kvs) == 2 {
			ret[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
		}
	}
	return ret, nil
}
