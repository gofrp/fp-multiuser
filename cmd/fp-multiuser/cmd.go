package main

import (
	"errors"
	"github.com/gofrp/fp-multiuser/pkg/server"
	"github.com/gofrp/fp-multiuser/pkg/server/controller"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
)

const version = "0.0.2"

var (
	showVersion bool
	configFile  string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frps-multiuser")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./frps-multiuser.ini", "config file of frps-multiuser")
}

var rootCmd = &cobra.Command{
	Use:   "frps-multiuser",
	Short: "frps-multiuser is the server plugin of frp to support multiple users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			log.Println(version)
			return nil
		}
		common, tokens, ports, domains, subdomains, iniFile, err := ParseConfigFile(configFile)
		if err != nil {
			log.Printf("fail to start frps-multiuser : %v", err)
			return nil
		}
		s, err := server.New(controller.HandleController{
			CommonInfo: common,
			Tokens:     tokens,
			Ports:      ports,
			Domains:    domains,
			Subdomains: subdomains,
			ConfigFile: configFile,
			IniFile:    iniFile,
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

func ParseConfigFile(file string) (controller.CommonInfo, map[string]controller.TokenInfo, map[string][]string, map[string][]string, map[string][]string, *ini.File, error) {
	common := controller.CommonInfo{}
	users := make(map[string]controller.TokenInfo)
	ports := make(map[string][]string)
	domains := make(map[string][]string)
	subdomains := make(map[string][]string)

	iniFile, err := ini.LoadSources(ini.LoadOptions{
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
		return common, nil, nil, nil, nil, iniFile, err
	}

	commonSection, err := iniFile.GetSection("common")
	if err != nil {
		log.Printf("fail to get [common] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	pluginAddr := commonSection.Key("plugin_addr").Value()
	if len(pluginAddr) != 0 {
		common.PluginAddr = pluginAddr
	} else {
		common.PluginAddr = "0.0.0.0"
	}
	pluginPort := commonSection.Key("plugin_port").Value()
	if len(pluginPort) != 0 {
		port, err := strconv.Atoi(pluginPort)
		if err != nil {
			return common, nil, nil, nil, nil, iniFile, err
		}
		common.PluginPort = port
	} else {
		common.PluginPort = 7200
	}
	common.User = commonSection.Key("admin_user").Value()
	common.Pwd = commonSection.Key("admin_pwd").Value()

	usersSection, err := iniFile.GetSection("users")
	if err != nil {
		log.Printf("fail to get [user] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}

	disabledSection, err := iniFile.GetSection("disabled")
	if err != nil {
		log.Printf("fail to get [disabled] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}

	keys := usersSection.Keys()
	for _, key := range keys {
		var token = controller.TokenInfo{
			User:    key.Name(),
			Token:   key.Value(),
			Comment: key.Comment,
			Status:  !(disabledSection.HasKey(key.Name()) && disabledSection.Key(key.Name()).Value() == "disable"),
		}
		users[token.User] = token
	}

	portsSection, err := iniFile.GetSection("ports")
	if err != nil {
		log.Printf("fail to get [ports] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range portsSection.Keys() {
		var user = key.Name()
		var port = strings.Split(key.Value(), ",")
		ports[user] = port
	}

	domainsSection, err := iniFile.GetSection("domains")
	if err != nil {
		log.Printf("fail to get [domains] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range domainsSection.Keys() {
		var user = key.Name()
		var domain = strings.Split(key.Value(), ",")
		domains[user] = domain
	}

	subdomainsSection, err := iniFile.GetSection("subdomains")
	if err != nil {
		log.Printf("fail to get [subdomains] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range subdomainsSection.Keys() {
		var user = key.Name()
		var subdomain = strings.Split(key.Value(), ",")
		subdomains[user] = subdomain
	}

	return common, users, ports, domains, subdomains, iniFile, nil
}
