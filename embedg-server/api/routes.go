package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	embedgapp "github.com/merlinfuchs/embed-generator/embedg-app"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/assistant"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/auth"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/custom_bots"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/guilds"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/images"
	premium_handler "github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/premium"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/saved_messages"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/scheduled_messages"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/send_message"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/shared_messages"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/handlers/users"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/helpers"
	"github.com/merlinfuchs/embed-generator/embedg-server/api/session"
	"github.com/merlinfuchs/embed-generator/embedg-server/bot"
	"github.com/merlinfuchs/embed-generator/embedg-server/util"
	embedgsite "github.com/merlinfuchs/embed-generator/embedg-site"
	"github.com/spf13/viper"
)

func registerRoutes(app *fiber.App, stores *stores, bot *bot.Bot, managers *managers) {
	authHandler := auth.New(stores.pg, bot, managers.session)
	app.Get("/api/auth/login", authHandler.HandleAuthRedirect)
	app.Get("/api/auth/callback", authHandler.HandleAuthCallback)
	app.Post("/api/auth/exchange", helpers.WithRequestBody(authHandler.HandleAuthExchange))
	app.Get("/api/auth/logout", authHandler.HandleAuthLogout)

	sessionMiddleware := session.NewSessionMiddleware(managers.session)

	usersHandler := users.New(stores.pg, managers.premium)
	usersGroup := app.Group("/api/users", sessionMiddleware.SessionRequired())
	usersGroup.Get("/:userID", usersHandler.HandleGetUser)

	savedMessagesHandler := saved_messages.New(stores.pg, managers.access)
	savedMessagesGroup := app.Group("/api/saved-messages", sessionMiddleware.SessionRequired())
	savedMessagesGroup.Get("/", savedMessagesHandler.HandleListSavedMessages)
	savedMessagesGroup.Post("/", helpers.WithRequestBodyValidated(savedMessagesHandler.HandleCreateSavedMessage))
	savedMessagesGroup.Patch("/", helpers.WithRequestBodyValidated(savedMessagesHandler.HandleImportSavedMessages))
	savedMessagesGroup.Put("/:messageID", helpers.WithRequestBodyValidated(savedMessagesHandler.HandleUpdateSavedMessage))
	savedMessagesGroup.Delete("/:messageID", savedMessagesHandler.HandleDeleteSavedMessage)

	sharedMessageHandler := shared_messages.New(bot, stores.pg)
	sharedMessagesGroup := app.Group("/api/shared-messages")
	sharedMessagesGroup.Post("/", helpers.WithRequestBodyValidated(sharedMessageHandler.HandleCreateSharedMessage))
	sharedMessagesGroup.Get("/:messageID", sharedMessageHandler.HandleGetSharedMessage)

	assistantHandler := assistant.New(stores.pg, managers.access, managers.premium)
	app.Post("/api/assistant/message", sessionMiddleware.SessionRequired(), helpers.WithRequestBody(assistantHandler.HandleAssistantGenerateMessage))

	guildsHanlder := guilds.New(stores.pg, bot, managers.access, managers.premium)
	guildsGroup := app.Group("/api/guilds", sessionMiddleware.SessionRequired())
	guildsGroup.Get("/", guildsHanlder.HandleListGuilds)
	guildsGroup.Get("/:guildID", guildsHanlder.HandleGetGuild)
	guildsGroup.Get("/:guildID/channels", guildsHanlder.HandleListGuildChannels)
	guildsGroup.Get("/:guildID/roles", guildsHanlder.HandleListGuildRoles)
	guildsGroup.Get("/:guildID/emojis", guildsHanlder.HandleListGuildEmojis)
	guildsGroup.Get("/:guildID/stickers", guildsHanlder.HandleListGuildStickers)
	guildsGroup.Get("/:guildID/branding", guildsHanlder.HandleGetGuildBranding)

	sendMessageHandler := send_message.New(bot, stores.pg, managers.access, managers.actionParser)
	app.Post("/api/send-message/channel", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(sendMessageHandler.HandleSendMessageToChannel))
	app.Post("/api/send-message/webhook", helpers.WithRequestBodyValidated(sendMessageHandler.HandleSendMessageToWebhook))
	app.Post("/api/restore-message/channel", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(sendMessageHandler.HandleRestoreMessageFromChannel))
	app.Post("/api/restore-message/webhook", helpers.WithRequestBodyValidated(sendMessageHandler.HandleRestoreMessageFromWebhook))

	premiumHandler := premium_handler.New(stores.pg, bot, managers.access, managers.premium)
	app.Get("/api/premium/features", sessionMiddleware.SessionRequired(), premiumHandler.HandleGetFeatures)
	app.Get("/api/premium/entitlements", sessionMiddleware.SessionRequired(), premiumHandler.HandleListEntitlements)

	customBotHandler := custom_bots.New(stores.pg, bot, managers.access, managers.premium, managers.actionParser)
	app.Post("/api/custom-bot", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(customBotHandler.HandleConfigureCustomBot))
	app.Put("/api/custom-bot/presence", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(customBotHandler.HandleUpdateCustomBotPresence))
	app.Get("/api/custom-bot", sessionMiddleware.SessionRequired(), customBotHandler.HandleGetCustomBot)
	app.Delete("/api/custom-bot", sessionMiddleware.SessionRequired(), customBotHandler.HandleDisableCustomBot)
	app.Get("/api/custom-bot/commands", sessionMiddleware.SessionRequired(), customBotHandler.HandleListCustomCommands)
	app.Get("/api/custom-bot/commands/:commandID", sessionMiddleware.SessionRequired(), customBotHandler.HandleGetCustomCommand)
	app.Post("/api/custom-bot/commands", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(customBotHandler.HandleCreateCustomCommand))
	app.Put("/api/custom-bot/commands/:commandID", sessionMiddleware.SessionRequired(), helpers.WithRequestBodyValidated(customBotHandler.HandleUpdateCustomCommand))
	app.Delete("/api/custom-bot/commands/:commandID", sessionMiddleware.SessionRequired(), customBotHandler.HandleDeleteCustomCommand)
	app.Post("/api/custom-bot/commands/deploy", sessionMiddleware.SessionRequired(), customBotHandler.HandleDeployCustomCommands)
	app.Post("/api/gateway/:customBotID", customBotHandler.HandleCustomBotInteraction)

	imagesHandler := images.New(stores.pg, managers.access, managers.premium, stores.blob)
	app.Post("/api/images", sessionMiddleware.SessionRequired(), imagesHandler.HandleUploadImage)
	app.Get("/api/images/:imageID", sessionMiddleware.SessionRequired(), imagesHandler.HandleGetImage)
	app.Get("/cdn/images/:imageKey", imagesHandler.HandleDownloadImage)

	scheduledMessagesHandler := scheduled_messages.New(stores.pg, managers.access, managers.premium)
	scheduledMessagesGroup := app.Group("/api/scheduled-messages", sessionMiddleware.SessionRequired())
	scheduledMessagesGroup.Get("/", scheduledMessagesHandler.HandleListScheduledMessages)
	scheduledMessagesGroup.Post("/", helpers.WithRequestBodyValidated(scheduledMessagesHandler.HandleCreateScheduledMessage))
	scheduledMessagesGroup.Get("/:messageID", scheduledMessagesHandler.HandleGetScheduledMessage)
	scheduledMessagesGroup.Put("/:messageID", helpers.WithRequestBodyValidated(scheduledMessagesHandler.HandleUpdateScheduledMessage))
	scheduledMessagesGroup.Delete("/:messageID", scheduledMessagesHandler.HandleDeleteScheduledMessage)

	app.Get("/invite", func(c *fiber.Ctx) error {
		return c.Redirect(util.BotInviteURL(), 302)
	})

	app.Get("/discord", func(c *fiber.Ctx) error {
		return c.Redirect(viper.GetString("links.discord"), 302)
	})

	app.Get("/source", func(c *fiber.Ctx) error {
		return c.Redirect(viper.GetString("links.source"), 302)
	})

	app.Get("/premium", func(c *fiber.Ctx) error {
		return c.Redirect(fmt.Sprintf("https://discord.com/application-directory/%s/premium", viper.GetString("discord.client_id")), 302)
	})

	// Serve statix files
	app.Use("/app/", filesystem.New(filesystem.Config{
		Root:         http.FS(embedgapp.DistFS),
		Browse:       false,
		NotFoundFile: "dist/index.html",
		PathPrefix:   "/dist",
	}))

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(embedgsite.DistFS),
		Browse:       false,
		NotFoundFile: "dist/index.html",
		PathPrefix:   "/dist",
	}))
}
