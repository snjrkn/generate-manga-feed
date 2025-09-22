package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"github.com/snjrkn/generate-manga-feed/pkg/generatemangafeed"
)

type Command struct {
	Desc string
	Run  func() (string, error)
}

var commands = map[string]Command{
	"kuragefarm": {
		Desc: "Print kurage farm RSS feed",
		Run:  generatemangafeed.KurageFarm,
	},
	"comicdaysoneshot": {
		Desc: "Print comic-days oneshot RSS feed",
		Run:  generatemangafeed.ComicdaysOneshot,
	},
	"comicdaysnewcomer": {
		Desc: "Print comic-days newcomer RSS feed",
		Run:  generatemangafeed.ComicdaysNewcomer,
	},
	"andsofa": {
		Desc: "Print andsofa RSS feed",
		Run:  generatemangafeed.AndSofa,
	},
	"toti": {
		Desc: "Print toti RSS feed",
		Run:  generatemangafeed.Toti,
	},
	"matogrosso": {
		Desc: "Print matogrosso RSS feed",
		Run:  generatemangafeed.Matogrosso,
	},
	"kuragebunchaward": {
		Desc: "Print kuragebunch award RSS feed",
		Run:  generatemangafeed.KurageBunchAward,
	},
	"comicessaygekijo": {
		Desc: "Print comic-essay gekijo RSS feed",
		Run:  generatemangafeed.ComicEssayGekijo,
	},
	"comiplexoneshot": {
		Desc: "Print comiplex oneshot RSS feed",
		Run:  generatemangafeed.ComiplexOneshot,
	},
	"comicboostoneshot": {
		Desc: "Print comic-boost oneshot RSS feed",
		Run:  generatemangafeed.ComicBoostOneshot,
	},
	"younganimaloneshot": {
		Desc: "Print younganimal oneshot RSS feed",
		Run:  generatemangafeed.YoungAnimalOneshot,
	},
	"comicessaycontest": {
		Desc: "Print comic-essay contest RSS feed",
		Run:  generatemangafeed.ComicEssayContest,
	},
	"comicbunchkaiaward": {
		Desc: "Print comicbunch-kai award RSS feed",
		Run:  generatemangafeed.ComicBunchKaiAward,
	},
	"comicbunchkaioneshot": {
		Desc: "Print comicbunch-kai oneshot RSS feed",
		Run:  generatemangafeed.ComicBunchKaiOneshot,
	},
	"afternoonaward": {
		Desc: "Print afternoon award RSS feed",
		Run:  generatemangafeed.AfternoonAward,
	},
	"shonenmagazineaward": {
		Desc: "Print shonen magazine award RSS feed",
		Run:  generatemangafeed.ShonenMagazineAward,
	},
	"shonenmagazinerise": {
		Desc: "Print shonen magazine rise RSS feed",
		Run:  generatemangafeed.ShonenMagazineRise,
	},
	"championcrossoneshot": {
		Desc: "Print champion cross oneshot RSS feed",
		Run:  generatemangafeed.ChampionCrossOneshot,
	},
	"kuragebunchoneshot": {
		Desc: "Print kuragebunch oneshot RSS feed",
		Run:  generatemangafeed.KurageBunchOneshot,
	},
	"comicactiononeshot": {
		Desc: "Print comic-acticon oneshot RSS feed",
		Run:  generatemangafeed.ComicActionOneshot,
	},
}

var version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		showHelp()
		os.Exit(0)
	}

	subcommand := args[0]

	switch subcommand {
	case "version":
		showVersion()
	case "help":
		showHelp()
	default:
		cmd, ok := commands[subcommand]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unknown subcommand '%s'\n", subcommand)
			showHelp()
			os.Exit(1)
		}

		if checkInternetAccess() {
			str, err := cmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(str)
		} else {
			fmt.Fprintln(os.Stderr, "Error: no internet access")
			os.Exit(1)
		}
	}
}

func checkInternetAccess() bool {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func showVersion() {
	fmt.Printf("generate-manga-feed version: %s\n", version)
}

func showHelp() {
	fmt.Println("Usage: generate-manga-feed [subcommand]")
	fmt.Println("\nAvailable subcommands:")
	fmt.Println("  version               Show application version")
	fmt.Println("  help                  Show this help message")
	fmt.Println("  --------")
	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)
	for _, name := range commandNames {
		cmd := commands[name]
		fmt.Printf("  %-21s %s\n", name, cmd.Desc)
	}
	fmt.Println("")
}
