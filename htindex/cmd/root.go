// Copyright © 2019 Dmitry Mozzherin <dmozzherin@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gnames/htindex"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion string
	buildDate    string
	cfgFile      string
	opts         []htindex.Option
)

type config struct {
	Root   string
	Input  string
	Output string
	Jobs   int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "htindex",
	Short: "creates an index of scientific names in Hathi Trust corpus.",
	Long: `Hathi Trust is a large collection of public and private textual
	information. The htindex program allows to use its data to find in it
	scientific names.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag(cmd)
		opts = getOpts()
		opts = getFlags(opts, cmd)
		hti := htindex.NewHTindex(opts...)
		fmt.Println(hti)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string, date string) {
	buildVersion = ver
	buildDate = date
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("version", "v", false, "htindex version and build timestamp")
	rootCmd.Flags().StringP("root", "r", "", "root path to add to the input file data")
	rootCmd.Flags().StringP("input", "i", "", "path to the input data file")
	rootCmd.Flags().StringP("output", "o", "", "path to the output directory")
	rootCmd.Flags().IntP("jobs", "j", 200, "number of workers (jobs)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".htindex" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".htindex")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func versionFlag(cmd *cobra.Command) {
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if version {
		fmt.Printf("\nversion: %s\n\nbuild:   %s\n\n",
			buildVersion, buildDate)
		os.Exit(0)
	}
}

func getOpts() []htindex.Option {
	var opts []htindex.Option
	cfg := &config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	if cfg.Root != "" {
		opts = append(opts, htindex.OptRoot(cfg.Root))
	}
	if cfg.Input != "" {
		opts = append(opts, htindex.OptInput(cfg.Input))
	}
	if cfg.Output != "" {
		opts = append(opts, htindex.OptOutput(cfg.Output))
	}
	if cfg.Jobs != 0 {
		opts = append(opts, htindex.OptJobs(cfg.Jobs))
	}
	return opts
}

func getFlags(opts []htindex.Option, cmd *cobra.Command) []htindex.Option {
	root, err := cmd.Flags().GetString("root")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if root != "" {
		opts = append(opts, htindex.OptRoot(root))
	}
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if input != "" {
		opts = append(opts, htindex.OptInput(input))
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if output != "" {
		opts = append(opts, htindex.OptOutput(output))
	}
	jobs, err := cmd.Flags().GetInt("jobs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if jobs > 0 {
		opts = append(opts, htindex.OptJobs(jobs))
	}
	return opts
}
