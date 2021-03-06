// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"github.com/almaznur91/splitty/internal/bot"
	"github.com/almaznur91/splitty/internal/events"
	"github.com/almaznur91/splitty/internal/handler"
	"github.com/almaznur91/splitty/internal/repository"
	"github.com/almaznur91/splitty/internal/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

func initApp(ctx context.Context, cfg *config) (*events.TelegramListener, func(), error) {
	botConfig := initBotConfig(cfg)
	botAPI, err := initTelegramApi(cfg, botConfig)
	if err != nil {
		return nil, nil, err
	}
	database, cleanup, err := initMongoConnection(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	mongoChatStateRepository := repository.NewChatStateRepository(database)
	chatStateService := service.NewChatStateService(mongoChatStateRepository)
	mongoButtonRepository := repository.NewButtonRepository(database)
	buttonService := service.NewButtonService(mongoButtonRepository)
	mongoUserRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(mongoUserRepository)
	startScreen := bot.NewStartScreen(chatStateService, buttonService, userService, botConfig)
	v := ProvideBotList(startScreen)
	errorHandler := handler.NewErrorHandler()
	telegramListener, err := initTelegramConfig(botAPI, v, buttonService, userService, chatStateService, errorHandler)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return telegramListener, func() {
		cleanup()
	}, nil
}

// wire.go:

var bots = wire.NewSet(bot.NewStartScreen)

func ProvideBotList(b2 *bot.StartScreen) []bot.Interface {
	return []bot.Interface{b2}
}
