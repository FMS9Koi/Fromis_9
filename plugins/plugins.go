package plugins

import (
    "github.com/bwmarrin/discordgo"
    "github.com/sn0w/Karen/utils"
    "strings"
)

// Plugin interface to enforce a basic structure
type Plugin interface {
    // List of commands and aliases
    Commands() []string

    // Action to execute on message receive
    Action(
    command string,
    content string,
    msg *discordgo.Message,
    session *discordgo.Session,
    )

    // Plugin constructor
    Init(session *discordgo.Session)
}

// List of active plugins
var PluginList = []Plugin{
    About{},
    Stats{},
    Ping{},
    Invite{},
    Giphy{},
    Google{},
    RandomCat{},
    Stone{},
    Roll{},
    Reminders{},
    //Music{},
}

// CallBotPlugin iterates through the list of registered
// plugins and tries to guess which one is the intended call
// Fist match wins.
//
// command - The command that triggered this execution
// content - The content without command
// msg     - The message object
// session - The discord session
func CallBotPlugin(command string, content string, msg *discordgo.Message, session *discordgo.Session) {
    // Iterate over all plugins
    for _, plug := range PluginList {
        // Iterate over all commands of the current plugin
        for _, cmd := range plug.Commands() {
            if command == cmd {
                go safePluginCall(command, strings.Trim(content, " "), msg, session, plug)
                break
            }
        }
    }
}

// Wrapper that catches any panics from plugins
// Arguments: Same as CallBotPlugin().
func safePluginCall(command string, content string, msg *discordgo.Message, session *discordgo.Session, plug Plugin) {
    defer func() {
        err := recover()

        if err != nil {
            utils.SendError(session, msg, err)
        }
    }()

    plug.Action(command, content, msg, session)
}

// Getter for this plugin list
func GetPlugins() []Plugin {
    return PluginList
}
