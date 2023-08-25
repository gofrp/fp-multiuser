package main

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gofrp/fp-multiuser/pkg/server"

	"github.com/spf13/cobra"
)

const version = "0.0.2"

var (
	showVersion bool

	bindAddr  string
	tokenFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "l", "127.0.0.1:7200", "bind address")
	rootCmd.PersistentFlags().StringVarP(&tokenFile, "token_file", "c", "./tokens.ini", "token file")
}

var rootCmd = &cobra.Command{
	Use:   "frps-multiuser",
	Short: "frps-multiuser is the server plugin of frp to support multiple users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version)
			return nil
		}
		tokens, err := ParseTokensFromFile(tokenFile)
		if err != nil {
			log.Printf("fail to start frps-multiuser")
			return nil
		}
		s, err := server.New(server.Config{
			BindAddress: bindAddr,
			Tokens:      tokens,
		})
		if err != nil {
			return err
		}
		err = s.Run()
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ParseTokensFromFile(file string) (map[string]string, error) {
	ret := make(map[string]string)

	i, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:         false,
		InsensitiveSections: false,
		InsensitiveKeys:     false,
		IgnoreInlineComment: true,
		AllowBooleanKeys:    true,
	}, file)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			log.Printf("token file %s not found", file)
		} else {
			log.Printf("fail to parse token file %s : %v", file, err)
		}
		return nil, err
	}

	t, err := i.GetSection("user")
	if err != nil {
		log.Printf("fail to parse token file %s : %v", file, err)
		return nil, err
	}

	keys := t.Keys()
	for _, key := range keys {
		ret[strings.TrimSpace(key.Name())] = strings.TrimSpace(key.Value())
	}

	return ret, nil
}
