/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"io"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "drago",
	Short: "A brief description of your application",
	Long:  `Usage: dragon [-version] [-help] [-autocomplete-(un)install] <command> [args]`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "","/etc/drago.d", "config file (default is /etc/drago.d)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("ui", false)
	viper.SetDefault("datadir", "/tmp/drago")
    viper.SetDefault("client.enabled", true)
	viper.SetDefault("client.datadir", "/tmp/drago")
	viper.SetDefault("client.serveraddr", "localhost:8080")
	viper.SetDefault("client.iface", "wg0")
	viper.SetDefault("client.syncinterval", 60)
	viper.SetDefault("client.wgkey", "")
	viper.SetDefault("client.jwt", "")
	viper.SetDefault("server.enabled", true)
	viper.SetDefault("server.bindaddr", ":8080")
	viper.SetDefault("server.secret", "drago")
	viper.SetDefault("server.network", "10.0.0.0/24")


	LoadConfig(cfgFile)


	viper.Set("ui",viper.Get("ui"))
	viper.Set("datadir",viper.Get("datadir"))
	viper.Set("client.enabled",viper.Get("client.enabled"))
	viper.Set("client.datadir",viper.Get("client.datadir"))
	viper.Set("client.iface",viper.Get("client.iface"))
	viper.Set("client.serveraddr",viper.Get("client.serveraddr"))
	viper.Set("client.syncinterval",viper.Get("client.syncinterval"))
	viper.Set("client.wgkey",viper.Get("client.wgkey"))
	viper.Set("client.jwt",viper.Get("client.jwt"))
	viper.Set("server.enabled",viper.Get("server.enabled"))
	viper.Set("server.bindaddr",viper.Get("server.bindaddr"))
	viper.Set("server.secret",viper.Get("server.secret"))
	viper.Set("server.network",viper.Get("server.network"))

}


func LoadConfig(path string) (error) {

	//check if dir or file
	fi, err := os.Stat(path)

	if err != nil {
		return err
	}

	//handle dir
	if fi.IsDir() {
		return LoadConfigDir(path)
	}

	return LoadConfigFile(path)
}

func LoadConfigDir(dir string) (error) {
	
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("configuration path must be a directory: %s", dir)
	}

	file := ""
	err = nil
	for err != io.EOF {
		var fis []os.FileInfo
		fis, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return err
		}

		for _, fi := range fis {
			// Ignore directories
			if fi.IsDir() {
				continue
			}

			// Only care about files that are valid to load.
			name := fi.Name()
			skip := true
			if strings.HasSuffix(name, ".yml") {
				skip = false
			}
			if skip || file != ""  {
				continue
			}

			path := filepath.Join(dir, name)
			file = path
		}
	}

	// Fast-path if we have no files
	if file == "" {
		return nil
	}

	err = LoadConfigFile(file)
	if err != nil {
		return fmt.Errorf("Error loading %s: %s", f, err)
	}

	return nil
}

func LoadConfigFile(path string) (error) {

	viper.SetConfigFile(path)
	err := viper.MergeInConfig()

	// If a config file is found, read it in.
	if err == nil {
			fmt.Println("==> Loaded configuration from ", viper.ConfigFileUsed())
	} else {
			fmt.Println("==> Error loading configuration: ", err)
	}
	return err
}