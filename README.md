# gero

A framework to write Discord bots in Go. Based on [github.com/bwmarrin/discordgo](https://github.com/bwmarrin/discordgo).

# Getting started

First, we must create a new instance of `gero.Bot`:

```go
bot := gero.NewBot()
```

If we want, before anything, we can configure the logger:

```go
if err := bot.ConfigureLogger(verbose, logfile); err != nil {
	panic(err)
}
```

Then, we must create a session:

```go
if err := bot.CreateSession(token); err != nil {
	panic(err)
}
```

`token` is the Oauth2 token you get for your bot user in the Applications section of the official Discord website.

After we have a successfully created session, we should initiate a WSS connection with the Discord server:

```go
if err := bot.OpenWssConnection(); err != nil {
	panic(err)
}

defer bot.CloseWssConnection()
```

Finally, we make the program wait for an exit signal:

```go
bot.WaitProgramExit()
```

This is enough to build a program, but it does nothing. If we want to actually do something, we must either act directly with the `discordgo.Session` through the `bot.Session()` method, or use the intuitive commands handler provided by gero:

```go
prefixConfig := &gero.PrefixConfig{
	CheckForPrefix:		true,
	Prefixes:			  []string{"mybot ", "mbot ", "!"},
	ConsiderMentionPrefix: true,
}

commands := map[string]gero.CommandHandler{
	"^hello world +(.+?)$": handleHelloWorld,
	"^help$":               handleHelp,
}

if err := bot.HandleCommands(prefixConfig, commands); err != nil {
	panic(err)
}

// ...

func handleHelloWorld(b *gero.Bot, m *discordgo.MessageCreate, params []string) error {
	fmt.Println(params)
}

func handleHelp(b *gero.Bot, m *discordgo.MessageCreate, params []string) error {
	fmt.Println(params)
}
```

Here is an example of a whole program:

```go
package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/xrash/gero"
)

func main() {
	var token, logfile string
	flag.StringVar(&token, "token", "", "Discord Application token.")
	flag.StringVar(&logfile, "logfile", "", "Writtable logfile filename.")
	flag.Parse()

	bot := gero.NewBot()

	if err := bot.ConfigureLogger(true, logfile); err != nil {
		panic(err)
	}

	if err := bot.CreateSession(token); err != nil {
		panic(err)
	}

	if err := bot.ActivateStatsCollector(); err != nil {
		panic(err)
	}

	prefixConfig := &gero.PrefixConfig{
		CheckForPrefix:        true,
		Prefixes:              []string{"mybot ", "mbot ", "!"},
		ConsiderMentionPrefix: true,
	}

	commands := map[string]gero.CommandHandler{
		"^hello world +(.+?)$": handleHelloWorld,
		"^help$":               handleHelp,
	}

	if err := bot.HandleCommands(prefixConfig, commands); err != nil {
		panic(err)
	}

	if err := bot.OpenWssConnection(); err != nil {
		panic(err)
	}

	defer bot.CloseWssConnection()

	bot.WaitProgramExit()
}

func handleHelloWorld(b *gero.Bot, m *discordgo.MessageCreate, params []string) error {
	fmt.Println(params)
	return nil
}

func handleHelp(b *gero.Bot, m *discordgo.MessageCreate, params []string) error {
	fmt.Println(params)
	return nil
}
```
