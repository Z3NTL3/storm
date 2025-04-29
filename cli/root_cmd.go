package cli

import (
	"fmt"
	"log"
	"os"
	"time"
	"z3ntl3/storm/bot/client"
	"z3ntl3/storm/globals"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lightup",
	Short: "Starts stress testing using L7",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New()
		channel := make(chan client.MessageContext, client.PoolSize)
		var i uint64 = 0

		// run on the background
		go func() {
			for {
				i++
				go c.Stress(*c.Rsrc.Proxies.Next(), i, channel)

				if c.ShouldExit() {
					close(channel)
					break
				}
			}
		}()

		for {
			msg, ok := <-channel
			if !ok {
				c.SigExit()
				os.Exit(0)
			}

			if msg.Kill && msg.Err != nil {
				c.SigExit()
				log.Fatal(msg.Err)
			}

			if msg.Kill && msg.Msg != "" {
				c.SigExit()
				fmt.Printf("%s\n", msg.Msg)
				return
			}

			if msg.Kill {
				c.SigExit()
				return
			}

			if msg.Msg != "" {
				fmt.Printf("h%s\n", msg.Msg)
			}
		}
	},
}

func addFlags() {
	opts := []struct {
		data_ref any
		value    any
		name     string
		usage    string
		required bool
	}{
		{
			data_ref: &globals.TargetURL,
			name:     "target",
			usage:    "Target URI, including the scheme, either 'http' or 'https'",
			value:    "",
			required: true,
		}, {
			data_ref: &globals.Timeout,
			name:     "timeout",
			usage:    "General timeout for both the socket and general proxy transport",
			value:    "5s",
			required: true,
		}, {
			data_ref: &globals.ProxyFile,
			name:     "proxy_file",
			usage:    "Relative file path to your file with proxies",
			value:    "data/proxy.txt",
		}, {
			data_ref: &globals.Accepts,
			name:     "accepts_file",
			usage:    "Relative file path to your file with HTTP accept headers",
			value:    "data/accepts.txt",
		}, {
			data_ref: &globals.Headers,
			name:     "headers_file",
			usage:    "Relative file path to your file with HTTP headers",
			value:    "data/headers.txt",
		},
	}

	for _, v := range opts {
		switch v.name {
		case "timeout", "target":
			if v.name == "timeout" {
				d, err := time.ParseDuration(v.value.(string))
				if err != nil {
					log.Fatal(err)
				}

				rootCmd.Flags().
					DurationVar(v.data_ref.(*time.Duration), v.name, d, v.usage)
			} else {
				rootCmd.Flags().
					StringVar(v.data_ref.(*string), v.name, v.value.(string), v.usage)
			}

			rootCmd.MarkFlagRequired(v.name)
		default:
			rootCmd.Flags().
				StringVar(v.data_ref.(*string), v.name, v.value.(string), v.usage)
		}
	}
}

func Execute() {
	addFlags()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
