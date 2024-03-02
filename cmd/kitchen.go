/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/koksmat-com/koksmat/connectors"
	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var kitchenName string
var stationName string
var channelName string
var tenantName string = "365adm"
var journeyId string

func UnEscape(s string) string {
	ss, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return ss

}

// serveCmd represents the serve command
var kitchenCmd = &cobra.Command{
	Use:   "kitchen [[service]]",
	Short: "kitchen  ",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No service specified")
			return
		}
		service := args[0]
		switch service {
		case "list":
			k, err := kitchen.List()
			if err != nil {
				log.Fatalln(err)
			}
			printJSON(k)
			// restapi.All()

		default:

			log.Fatalln("Unknown service", service)
			return
		}
		//webserver.Run()
	},
}

var scriptcmd = &cobra.Command{
	Use:   "script [script]",
	Short: "Working with scripts",
	Long:  ``,
}

func cmd(use string, short string, long string, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   run,
	}
	cmd.Flags().StringVarP(&kitchenName, "kitchen", "k", "", "Kitchen (required)")
	cmd.MarkFlagRequired("kitchen")
	cmd.Flags().StringVarP(&stationName, "station", "s", "", "Station (required)")
	cmd.MarkFlagRequired("station")
	return cmd
}
func init() {

	rootCmd.AddCommand(kitchenCmd)
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "stations [kitchen]",
			Short: "List stations in kitchen",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				stations, err := kitchen.GetStations(name)
				if err != nil {
					log.Fatalln(err)
				}
				printJSON(stations)

				// kitchen := args[0]

			},
		})

	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "status [kitchen]",
			Short: "Get status of kitchen",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				status, err := kitchen.GetStatus(name, true)
				if err != nil {
					log.Fatalln(err)
				}
				printJSON(status)

				// kitchen := args[0]

			},
		})
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "build [kitchen]",
			Short: "Build kitchen",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				status, err := kitchen.Build(name)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(status)

				// kitchen := args[0]

			},
		})
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "open [kitchen]",
			Short: "Open kitchen",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				status, err := kitchen.Open(name)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(status)

				// kitchen := args[0]

			},
		})
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "launch [kitchen]",
			Short: "Launch kitchen",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				status, err := kitchen.Launch(name)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(status)

				// kitchen := args[0]

			},
		})
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "create [kitchen]",
			Short: "Create a new kitchen and change the current path to that",
			Args:  cobra.MinimumNArgs(1),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {
				name := args[0]
				err := kitchen.CreateKitchen(name)
				if err != nil {
					log.Fatalln(err)
				}
				printJSON("Created")

				// kitchen := args[0]

			},
		})
	kitchenCmd.AddCommand(scriptcmd)

	htmlCmd := &cobra.Command{
		Use:   "html [file]",
		Short: "Exports HTML from Markdown in script",
		Args:  cobra.MinimumNArgs(1),
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			root := viper.GetString("KITCHENROOT")
			filename, _ := url.QueryUnescape(args[0])
			file := path.Join(root, kitchenName, stationName, filename)
			markdown := ""
			switch filepath.Ext(file) {
			case ".ps1":
				md, _, err := kitchen.ReadMarkdownFromPowerShell(file)
				if err != nil {
					fmt.Println(err)
				}
				markdown = md
			case ".go":
				md, err := kitchen.ReadMarkdownFromGo(file)
				if err != nil {
					fmt.Println(err)
				}
				markdown = md
			default:
				fmt.Println("Unknown file type")
				return
			}

			html, _, err := kitchen.ParseMarkdown(false, filepath.Dir(file), markdown)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(html)

		},
	}
	scriptcmd.AddCommand(htmlCmd)
	markdownCmd := &cobra.Command{
		Use:   "markdown [file]",
		Short: "Exports Markdown in script",
		Args:  cobra.MinimumNArgs(1),
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			root := viper.GetString("KITCHENROOT")
			filename, _ := url.QueryUnescape(args[0])
			file := path.Join(root, kitchenName, stationName, filename)
			markdown := ""
			switch filepath.Ext(file) {
			case ".ps1":
				md, _, err := kitchen.ReadMarkdownFromPowerShell(file)
				if err != nil {
					fmt.Println(err)
				}
				markdown = md
			case ".go":
				md, err := kitchen.ReadMarkdownFromGo(file)
				if err != nil {
					fmt.Println(err)
				}
				markdown = md
			default:
				fmt.Println("Unknown file type")
				return
			}

			fmt.Println(markdown)

		},
	}
	scriptcmd.AddCommand(markdownCmd)
	htmlCmd.Flags().StringVarP(&kitchenName, "kitchen", "k", "", "Kitchen (required)")
	htmlCmd.MarkFlagRequired("kitchen")
	htmlCmd.Flags().StringVarP(&stationName, "station", "s", "", "Station (required)")
	htmlCmd.MarkFlagRequired("station")
	markdownCmd.Flags().StringVarP(&kitchenName, "kitchen", "k", "", "Kitchen (required)")
	markdownCmd.MarkFlagRequired("kitchen")
	markdownCmd.Flags().StringVarP(&stationName, "station", "s", "", "Station (required)")
	markdownCmd.MarkFlagRequired("station")
	metaCmd := &cobra.Command{
		Use:   "meta [file]",
		Short: "Exports Metadata from Markdown in script",
		Args:  cobra.MinimumNArgs(1),
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {

			fileName, err := url.QueryUnescape(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			metadata, err := kitchen.GetMetadata(kitchenName, stationName, fileName)
			printJSON(metadata)

		},
	}
	scriptcmd.AddCommand(metaCmd)
	metaCmd.Flags().StringVarP(&kitchenName, "kitchen", "k", "", "Kitchen (required)")
	metaCmd.MarkFlagRequired("kitchen")
	metaCmd.Flags().StringVarP(&stationName, "station", "s", "", "Station (required)")
	metaCmd.MarkFlagRequired("station")
	metaCmd.Flags().StringVarP(&journeyId, "journey", "j", "", "Journey ")

	scriptcmd.AddCommand(cmd("edit [file]", "Edit script", "", func(cmd *cobra.Command, args []string) {
		root := viper.GetString("KITCHENROOT")

		filename, _ := url.QueryUnescape(args[0])

		file := path.Join(root, kitchenName, stationName, filename)
		connectors.Execute("code", *&connectors.Options{}, file)
		fmt.Println("Opened", file)

	}))

	runcmd := cmd("run [file]", "Debug script", "", func(cmd *cobra.Command, args []string) {

		filename := UnEscape(args[0])
		mateContext, err := connectors.GetContext()
		if err != nil {
			log.Fatalln(err)
		}
		result, err := kitchen.Cook(true, mateContext.Tenant, kitchenName, stationName, journeyId, filename, nil)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)

	})
	runcmd.Flags().StringVarP(&channelName, "channel", "c", "", "Centrifugo channel to write back on")
	scriptcmd.AddCommand(runcmd)
	execcmd := cmd("execute [file]", "Execute script", "", func(cmd *cobra.Command, args []string) {

		filename := UnEscape(args[0])
		mateContext, err := connectors.GetContext()
		if err != nil {
			log.Fatalln(err)
		}
		sessionPath, err := kitchen.Cook(false, mateContext.Tenant, kitchenName, stationName, journeyId, filename, nil)
		if err != nil {
			log.Fatalln(err)
		}

		sessionPath2 := strings.Replace(string(sessionPath), "\n", "", -1)
		scriptPath := path.Join(string(sessionPath2), "run.ps1")

		ps1, err := os.ReadFile(scriptPath)
		if err != nil {

			log.Fatalln("could not run powershell script")
		}
		code := string(ps1)
		ps1args := strings.Join(args, " ")

		newcode := strings.ReplaceAll(code, "##ARGS##", ps1args)
		//log.Println(newcode)
		err = os.WriteFile(scriptPath, []byte(newcode), 0644)

		if err != nil {

			log.Fatalln("could not run powershell script")
		}

		oscmd := exec.Command("pwsh", "-f", "run.ps1", "-nologo", "-noprofile")

		oscmd.Dir = sessionPath2

		pipe, _ := oscmd.StdoutPipe()
		combinedOutput := []byte{}

		err = oscmd.Start()
		go func(p io.ReadCloser) {
			reader := bufio.NewReader(pipe)
			line, err := reader.ReadString('\n')
			for err == nil {
				//log.Print(line)
				combinedOutput = append(combinedOutput, []byte(line)...)
				line, err = reader.ReadString('\n')
			}
		}(pipe)
		err = oscmd.Wait()

		if err != nil {
			log.Println(fmt.Sprint(err), sessionPath2) //+ ": " + string(combinedOutput))
			log.Fatalln("could not run powershell script")
		}

		fmt.Println(string(combinedOutput))

	})
	execcmd.Flags().StringVarP(&channelName, "channel", "c", "", "Centrifugo channel to write back on")

	scriptcmd.AddCommand(execcmd)
	run2cmd := cmd("setup [file]", "Setup script", "", func(cmd *cobra.Command, args []string) {

		filename := UnEscape(args[0])
		mateContext, err := connectors.GetContext()
		if err != nil {
			log.Fatalln(err)
		}
		result, err := kitchen.Cook(false, mateContext.Tenant, kitchenName, stationName, journeyId, filename, nil)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)

	})
	scriptcmd.AddCommand(run2cmd)
	kitchenCmd.AddCommand(

		&cobra.Command{
			Use:   "boot",
			Short: "Boot kitchens",
			Args:  cobra.MinimumNArgs(0),
			Long:  ``,

			Run: func(cmd *cobra.Command, args []string) {

				err := kitchen.Boot()
				if err != nil {
					log.Fatalln(err)
				}

				// kitchen := args[0]

			},
		})
}
