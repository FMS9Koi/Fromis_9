package nugugame

// Init when the bot starts up
import (
	"regexp"
	"strings"

	"github.com/Seklfreak/Robyul2/helpers"
	"github.com/bwmarrin/discordgo"
)

// module struct
type Module struct{}

var gameGenders = map[string]string{
	"boy":   "boy",
	"boys":  "boy",
	"girl":  "girl",
	"girls": "girl",
	"mixed": "mixed",
}

func (m *Module) Init(session *discordgo.Session) {
	go func() {

		// refresh idols in difficulties
		idolsByDifficultyMutex.Lock()
		getModuleCache(NUGUGAME_DIFFICULTY_IDOLS_KEY, &idolsByDifficulty)
		idolsByDifficultyMutex.Unlock()

		// regex for idol and group names
		var err error
		alphaNumericRegex, err = regexp.Compile("[^a-zA-Z0-9가-힣]+")
		helpers.Relax(err)

		currentNuguGames = make(map[string][]*nuguGame)

		startDifficultyCacheLoop()

		// load all images and information
		loadMiscImages()
	}()
}

// Uninit called when bot is shutting down
func (m *Module) Uninit(session *discordgo.Session) {

}

// Will validate if the passed command entered is used for this plugin
func (m *Module) Commands() []string {
	return []string{
		"nugugame",
		"ng",
	}
}

// Main Entry point for the plugin
func (m *Module) Action(command string, content string, msg *discordgo.Message, session *discordgo.Session) {
	if !helpers.ModuleIsAllowed(msg.ChannelID, msg.ID, msg.Author.ID, helpers.ModulePermGames) {
		return
	}

	// process text after the initial command
	commandArgs := strings.Fields(content)

	if command == "nugugame" || command == "ng" {
		if len(commandArgs) > 0 {

			switch commandArgs[0] {
			case "list":

				helpers.RequireRobyulMod(msg, func() {
					listIdolsByDifficulty(msg, commandArgs)
				})
			case "refresh-nugugame":
				helpers.RequireRobyulMod(msg, func() {
					manualRefreshDifficulties(msg)
				})
			default:
				startNuguGame(msg, commandArgs)
			}
		} else {
			startNuguGame(msg, commandArgs)
		}
	}
}

func (m *Module) OnMessage(content string, msg *discordgo.Message, session *discordgo.Session) {
	for _, game := range getNuguGamesByChannelID(msg.ChannelID) {
		if game.WaitingForGuess {
			// if the game is not multiplayer, check the message author is the one who created the game
			if !game.IsMultigame && msg.Author.ID != game.User.ID {
				continue
			}
			game.GuessChannel <- msg
			break
		}
	}
}

///// Unused functions requried by ExtendedPlugin interface
func (m *Module) OnReactionAdd(reaction *discordgo.MessageReactionAdd, session *discordgo.Session) {
}
func (m *Module) OnMessageDelete(msg *discordgo.MessageDelete, session *discordgo.Session) {
}
func (m *Module) OnGuildMemberAdd(member *discordgo.Member, session *discordgo.Session) {
}
func (m *Module) OnGuildMemberRemove(member *discordgo.Member, session *discordgo.Session) {
}
func (m *Module) OnReactionRemove(reaction *discordgo.MessageReactionRemove, session *discordgo.Session) {
}
func (m *Module) OnGuildBanAdd(user *discordgo.GuildBanAdd, session *discordgo.Session) {
}
func (m *Module) OnGuildBanRemove(user *discordgo.GuildBanRemove, session *discordgo.Session) {
}