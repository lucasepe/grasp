package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasepe/grasp/generator"
	"github.com/spf13/cobra"
)

const (
	banner = `┌─┐┬─┐┌─┐┌─┐┌─┐
│ ┬├┬┘├─┤└─┐├─┘
└─┘┴└─┴ ┴└─┘┴  `

	appSummary = "Create strong passwords using words that are easy for you to remember."

	optNoDigits  = "no-digits"
	optNoSymbols = "no-symbols"
	optSize      = "size"
	optNoNL      = "no-newline"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MinimumNArgs(2),
		Use:                   fmt.Sprintf("%s <KEYWORD_1> <KEYWORD_1> [... KEYWORD_n]", appName()),
		Short:                 appSummary,
		Long:                  fmt.Sprintf("%s\n%s", banner, appSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			length, err := getPasswordSize(cmd)
			if err != nil {
				return err
			}

			//.Allow repeat for longer passwords in order to avoid generation errors
			allowRepeat := false
			if length > 16 {
				allowRepeat = true
			}

			//.Calculate the number of digits
			noDigits, _ := cmd.Flags().GetBool(optNoDigits)
			noSymbols, _ := cmd.Flags().GetBool(optNoSymbols)

			nl := "\n"
			if ok, _ := cmd.Flags().GetBool(optNoNL); ok {
				nl = ""
			}

			//.Password Generation
			gen, err := generator.NewGenerator(args)
			if err != nil {
				return err
			}

			res, err := gen.Generate(length, noDigits, noSymbols, allowRepeat)
			if err != nil {
				return err
			}

			fmt.Printf("%s%s", res, nl)
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.Flags().BoolP(optNoDigits, "d", false, "do not use digits")
	rootCmd.Flags().BoolP(optNoSymbols, "x", false, "do not use symbols")
	rootCmd.Flags().BoolP(optNoNL, "n", false, "do not append a newline when print result")

	rootCmd.Flags().StringP(optSize, "s", "M", fmt.Sprintf("password length in t-shirt size [%s]", availableSizes()))
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}} - crafted with passion by Luca Sepe <luca.sepe@gmail.com>
`)
}

func appName() string {
	return filepath.Base(os.Args[0])
}
