// Copyright 2022 Tristan Isham. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	// "errors"
	// "flag"
	// "fmt"
	// "html/template"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tristanisham/zvm/cli"
	"github.com/tristanisham/zvm/cli/meta"
	opts "github.com/urfave/cli/v2"

	"github.com/charmbracelet/log"

	_ "embed"
)

//go:embed help.txt
var helpTxt string

var zvm cli.ZVM

var zvmApp = &opts.App{
	Name:      "ZVM",
	HelpName:  "zvm",
	Version:   meta.VERSION,
	Copyright: "Copyright © 2022 Tristan Isham",
	Suggest:   true,
	Before: func(ctx *opts.Context) error {
		zvm = *cli.Initialize()
		return nil
	},
	// app-global flags
	Flags: []opts.Flag{
		&opts.StringFlag{
			Name:  "color",
			Usage: "enable or disable colored ZVM output",
			Value: "toggle",
			Action: func(ctx *opts.Context, val string) error {
				switch val {
				case "on", "yes", "enabled":
					zvm.Settings.YesColor()

				case "off", "no", "disabled":
					zvm.Settings.NoColor()

				default:
					zvm.Settings.ToggleColor()
				}

				return nil
			},
		},
	},
	Commands: []*opts.Command{
		{
			Name:    "install",
			Usage:   "download and install a version of Zig",
			Aliases: []string{"i"},
			Flags: []opts.Flag{
				&opts.BoolFlag{
					Name:    "zls",
					Aliases: []string{"z"},
					Usage:   "install ZLS",
				},
			},
			Description: "To install the latest version, use `master`",
			Args:        true,
			ArgsUsage:   " <ZIG VERSION>",
			Action: func(ctx *opts.Context) error {
				versionArg := strings.TrimPrefix(ctx.Args().First(), "v")

				if versionArg == "" {
					return errors.New("no version provided")
				}

				req := cli.ExtractInstall(versionArg)
				req.Version = strings.TrimPrefix(req.Version, "v")

				// Verify the requeste Zig version is good
				if err := zvm.ZigVersionIsValid(req.Package); err != nil {
					return err
				}

				// If ZLS install requested, verify that the versions match
				if ctx.Bool("zls") {
					if err := zvm.ZlsVersionIsValid(req.Package); err != nil {
						return err
					}
				}

				// Install Zig
				if err := zvm.Install(req.Package); err != nil {
					return err
				}

				// Install ZLS (if requested)
				if ctx.Bool("zls") {
					if err := zvm.InstallZls(req.Package); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			Name:  "use",
			Usage: "switch between versions of Zig",
			Args:  true,
			Action: func(ctx *opts.Context) error {
				versionArg := strings.TrimPrefix(ctx.Args().First(), "v")
				return zvm.Use(versionArg)
			},
		},
		{
			Name:    "list",
			Usage:   "list installed Zig versions",
			Aliases: []string{"ls"},
			Args:    true,
			Flags: []opts.Flag{
				&opts.BoolFlag{
					Name:    "all",
					Aliases: []string{"a"},
					Usage:   "list remote Zig versions available for download",
				},
			},
			Action: func(ctx *opts.Context) error {
				log.Debug("Version Map", "url", zvm.Settings.VersionMapUrl)
				if ctx.Bool("all") {
					return zvm.ListRemoteAvailable()
				} else {
					return zvm.ListVersions()
				}
			},
		},
		{
			Name:    "uninstall",
			Usage:   "remove an installed version of Zig",
			Aliases: []string{"rm"},
			Args:    true,
			Action: func(ctx *opts.Context) error {
				versionArg := strings.TrimPrefix(ctx.Args().First(), "v")
				return zvm.Uninstall(versionArg)
			},
		},
		{
			Name:  "clean",
			Usage: "remove build artifacts (good if you're a scrub)",
			Action: func(ctx *opts.Context) error {
				return zvm.Clean()
			},
		},
		{
			Name:  "upgrade",
			Usage: "self-update ZVM",
			Action: func(ctx *opts.Context) error {
				return zvm.Upgrade()
			},
		},
		{
			Name:  "vmu",
			Usage: "set ZVM's version map URL for custom Zig distribution servers",
			Args:  true,
			Action: func(ctx *opts.Context) error {
				url := ctx.Args().First()
				log.Debug("user passed vmu", "url", url)

				switch url {
				case "default":
					return zvm.Settings.ResetVersionMap()

				case "mach":
					if err := zvm.Settings.SetVersionMapUrl("https://machengine.org/zig/index.json"); err != nil {
						log.Info("Run `zvm vmu default` to reset your version map.")
						return err
					}

				default:
					if err := zvm.Settings.SetVersionMapUrl(url); err != nil {
						log.Info("Run `zvm vmu default` to reset your verison map.")
						return err
					}
				}

				return nil
			},
		},
	},
}

func main() {
	if _, ok := os.LookupEnv("ZVM_DEBUG"); ok {
		log.SetLevel(log.DebugLevel)
	}

	// run and report errors
	if err := zvmApp.Run(os.Args); err != nil {
		meta.CtaFatal(err)
	}
}

// if len(args) == 0 {
// helpMsg()
// zvm.AlertIfUpgradable()
// os.Exit(0)
// }
/*
		// zvm.AlertIfUpgradable()
		versionFlag := flag.Bool("version", false, "Print ZVM version information")
		// Install flags
		installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
		installDeps := installFlagSet.String("D", "", "Specify additional dependencies to install with Zig")

		// LS flags
		lsFlagSet := flag.NewFlagSet("ls", flag.ExitOnError)
		lsRemote := lsFlagSet.Bool("all", false, "List all available versions of Zig to install")

		// Global config
		sVersionMapUrl := flag.String("vmu", "", "Set ZVM's version map URL for custom Zig distribution servers")
		sColorToggle := flag.Bool("color", true, "Turn on or off ZVM's color output")
		flag.Parse()

		if *versionFlag {
			fmt.Println(meta.VerCopy)
			os.Exit(0)
		}

		if sVersionMapUrl != nil && len(*sVersionMapUrl) != 0 {
			log.Debug("user passed vmu", "url", *sVersionMapUrl)
			switch *sVersionMapUrl {
			case "default":
				if err := zvm.Settings.ResetVersionMap(); err != nil {
					meta.CtaFatal(err)
				}
			case "mach":
				if err := zvm.Settings.SetVersionMapUrl("https://machengine.org/zig/index.json"); err != nil {
					log.Info("Run `-vmu default` to reset your version map.")
					meta.CtaFatal(err)
				}

			default:

				if err := zvm.Settings.SetVersionMapUrl(*sVersionMapUrl); err != nil {
					log.Info("Run `-vmu default` to reset your version map.")
					meta.CtaFatal(err)
				}
			}

		}

		if sColorToggle != nil {
			if *sColorToggle != zvm.Settings.UseColor {
				if *sColorToggle {
					zvm.Settings.YesColor()
				} else {
					zvm.Settings.NoColor()
				}
			}

		}

		args = flag.Args()

	for i, arg := range args {

		switch arg {

		case "install", "i":
			installFlagSet.Parse(args[i+1:])
			installZls := false
			if *installDeps == "zls" {
				installZls = true
			}
			// signal to install zls after zig

			req := cli.ExtractInstall(args[len(args)-1])
			req.Version = strings.TrimPrefix(req.Version, "v")
			// log.Debug(req, "deps", *installDeps)

			if err := zvm.ValidateVersion(req.Package); err != nil {
				if errors.Is(err, cli.ErrInvalidZlsVersion) && !installZls {
					// TODO:(dogue) this if statement is stupid

					// don't report an error for the ZLS version
					// if ZLS wasn't requested to be installed
				} else {
					meta.CtaFatal(err)
				}
			}

			if err := zvm.Install(req.Package); err != nil {
				meta.CtaFatal(err)
			}

			if *installDeps != "" {
				switch *installDeps {
				case "zls":

					if err := zvm.InstallZls(req.Package); err != nil {
						meta.CtaFatal(err)
					}
				}
			}

			return
		case "use":
			if len(args) > i+1 {
				version := strings.TrimPrefix(args[i+1], "v")
				if err := zvm.Use(version); err != nil {
					meta.CtaFatal(err)
				}
				fmt.Printf("Switched to Zig v%s\n", version)
			}
			return

		case "ls":
			lsFlagSet.Parse(args[i+1:])
			log.Debug("Version Map", "url", zvm.Settings.VersionMapUrl)
			if *lsRemote {
				if err := zvm.ListRemoteAvailable(); err != nil {
					meta.CtaFatal(err)
				}
			} else {
				if err := zvm.ListVersions(); err != nil {
					meta.CtaFatal(err)
				}
			}

			return
		case "uninstall", "rm":
			if len(args) > i+1 {
				version := strings.TrimPrefix(args[i+1], "v")
				if err := zvm.Uninstall(version); err != nil {
					meta.CtaFatal(err)
				}
			}
			return

		case "sync":
			if err := zvm.Sync(); err != nil {
				meta.CtaFatal(err)
			}

		case "clean":
			// msg := "Clean is a beta command, and may not be included in the next release."
			// if zvm.Settings.UseColor {
			// 	fmt.Println(clr.Blue(msg))
			// } else {
			// 	fmt.Println(msg)
			// }

			if err := zvm.Clean(); err != nil {
				if zvm.Settings.UseColor {
					meta.CtaFatal(err)
				} else {
					meta.CtaFatal(err)
				}
			}
			return

		case "upgrade":
			if err := zvm.Upgrade(); err != nil {
				log.Error("this is a new command, and may have some issues.\nConsider reporting your problem on Github :)", "github", "https://github.com/tristanisham/zvm/issues")
				meta.CtaFatal(err)
			}

		case "version":
			fmt.Println(meta.VerCopy)
			return
		case "help":
			//zvm.Settings.UseColor
			helpMsg()

			return
			// Settings
		default:
			log.Fatalf("invalid argument %q. Please run `zvm help`.\n", arg)
		}

	}

}

func helpMsg() {
	helpTmpl, err := template.New("help").Parse(helpTxt)
	if err != nil {
		fmt.Printf("Sorry! There was a rendering error (%q). The version is %s\n", err, meta.VERSION)
		fmt.Println(helpTxt)
		return
	}

	if err := helpTmpl.Execute(os.Stdout, map[string]string{"Version": meta.VERSION}); err != nil {
		fmt.Printf("Sorry! There was a rendering error (%q). The version is %s\n", err, meta.VERSION)
		fmt.Println(helpTxt)
		return
	}
}
*/
