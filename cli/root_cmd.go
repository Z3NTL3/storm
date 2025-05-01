package cli

import (
	"fmt"
	"log"
	"os"
	"time"
	"z3ntl3/storm/bot/client"
	"z3ntl3/storm/globals"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), globals.Timeout)
					defer cancel()

					done := make(chan int)
					go c.Stress(*c.Rsrc.Proxies.Next(), i, channel, done)

					// when the wrapping goroutine exits the inner goroutine also terminates, which is when operation is done
					// or timeout occurs
					select {
					case <-done:
					case <-ctx.Done():
						return
					}
				}()

				if c.ShouldExit() {
					close(channel)
					break
				}

			}
		}()

		for {
			if c.ShouldExit() {
				break
			}

			msg, ok := <-channel
			if !ok {
				c.SigExit()
				break
			}

			if msg.Kill {
				c.SigExit()

				// precedence over msg when both set
				// notice that we only record major errors with [msg.Kill]
				if msg.Err != nil {
					log.Fatal(msg.Err)
				}

				if msg.Msg != "" {
					log.Fatal(msg.Msg)
				}
			}

			if msg.Msg != "" {
				fmt.Println(msg.Msg)
				continue
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
			data_ref: &globals.ProxyProto,
			name:     "proto",
			usage:    "Proxy protocol",
			value:    "http",
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
		}, {
			data_ref: &globals.Refs,
			name:     "refs_file",
			usage:    "Relative file path to your file with HTTP referer headers",
			value:    "data/refs.txt",
		}, {
			data_ref: &globals.UAs,
			name:     "uas_file",
			usage:    "Relative file path to your file with HTTP user-agent headers",
			value:    "data/uas.txt",
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
