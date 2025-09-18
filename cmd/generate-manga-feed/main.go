package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/snjrkn/generate-manga-feed/internal/service"
)

type Command struct {
	Desc string
	Run  func() (string, error)
}

var commands = map[string]Command{
	"kuragefarm": {
		Desc: "Print kurage farm RSS feed",
		Run:  service.KurageFarm().MakeFeed,
	},
	"comicdaysoneshot": {
		Desc: "Print comicdays oneshot RSS feed",
		Run:  service.ComicDaysOneshot().MakeFeed,
	},
	"comicdaysnewcomer": {
		Desc: "Print comicdays newcomer RSS feed",
		Run:  service.ComicDaysNewcomer().MakeFeed,
	},
	"andsofa": {
		Desc: "Print andsofa RSS feed",
		Run:  service.AndSofa().MakeFeed,
	},
	"toti": {
		Desc: "Print toti RSS feed",
		Run:  service.Toti().MakeFeed,
	},
	"matogrosso": {
		Desc: "Print matogrosso RSS feed",
		Run:  service.Matogrosso().MakeFeed,
	},
	"kuragebunchaward": {
		Desc: "Print kuragebunch award RSS feed",
		Run:  service.KurageBunchAward().MakeFeed,
	},
	"comicessaygekijo": {
		Desc: "Print comic essay gekijo RSS feed",
		Run:  service.ComicEssayGekijo().MakeFeed,
	},
	"comiplexoneshot": {
		Desc: "Print comiplex oneshot RSS feed",
		Run:  service.ComiplexOneshot().MakeFeed,
	},
	"comicboostoneshot": {
		Desc: "Print comicboost oneshot RSS feed",
		Run:  service.ComicBoostOneshot().MakeFeed,
	},
	"younganimaloneshot": {
		Desc: "Print younganimal oneshot RSS feed",
		Run:  service.YoungAnimalOneshot().MakeFeed,
	},
	"comicessaycontest": {
		Desc: "Print comic essay contest RSS feed",
		Run:  service.ComicEssayContest().MakeFeed,
	},
	"comicbunckaiaward": {
		Desc: "Print comicbunchi kai award RSS feed",
		Run:  service.ComicBunchKaiAward().MakeFeed,
	},
	"comicbunckaioneshot": {
		Desc: "Print comicbunchi kai oneshot RSS feed",
		Run:  service.ComicBunchKaiOneshot().MakeFeed,
	},
	"afternoonaward": {
		Desc: "Print afternoon award RSS feed",
		Run:  service.AfternoonAward().MakeFeed,
	},
	"shonenmagazineaward": {
		Desc: "Print shonen magazine award RSS feed",
		Run:  service.ShonenMagazineAward().MakeFeed,
	},
	"shonenmagazinerise": {
		Desc: "Print shonen magazine rise RSS feed",
		Run:  service.ShonenMagazineRise().MakeFeed,
	},
	"championcrossoneshot": {
		Desc: "Print championcross oneshot RSS feed",
		Run:  service.ChampionCrossOneshot().MakeFeed,
	},
	"kuragebunchoneshot": {
		Desc: "Print kuragebunch oneshot RSS feed",
		Run:  service.KurageBunchOneshot().MakeFeed,
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

		str, err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(str)
	}
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
