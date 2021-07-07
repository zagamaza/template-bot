package bot

import (
	"context"
	"github.com/almaznur91/splitty/internal/api"
	"github.com/go-pkgz/syncs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const start string = "/start"

//actions
const (
	createRoom api.Action = "create_room"
	viewStart  api.Action = "start"
)

// Interface is a bot reactive spec. response will be sent if "send" result is true
type Interface interface {
	OnMessage(ctx context.Context, update *api.Update) (api.TelegramMessage, error)
	HasReact(update *api.Update) bool
}

// SuperUser defines interface checking ig user name in su list
type SuperUser interface {
	IsSuper(userName string) bool
}

// MultiBot combines many bots to one virtual
type MultiBot []Interface

// OnMessage pass msg to all bots and collects reposnses (combining all of them)
//noinspection GoShadowedVar
func (b MultiBot) OnMessage(ctx context.Context, update *api.Update) (api.TelegramMessage, error) {

	resps := make(chan api.TelegramMessage)
	errors := make(chan error)

	wg := syncs.NewSizedGroup(4)
	for _, bot := range b {
		bot := bot
		wg.Go(func(ctx context.Context) {
			if bot.HasReact(update) {
				resp, err := bot.OnMessage(ctx, update)
				if err != nil {
					errors <- err
				} else {
					resps <- resp
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(resps)
		close(errors)
	}()

	message := &api.TelegramMessage{Chattable: []tgbotapi.Chattable{}}
	var eror error

tobreake:
	for {
		select {
		case r, ok := <-resps:
			if !ok {
				break tobreake
			}
			message.Chattable = append(message.Chattable, r.Chattable...)
			message.InlineConfig = r.InlineConfig
			message.CallbackConfig = r.CallbackConfig
			message.Redirect = r.Redirect
			message.Send = true

		case err, ok := <-errors:
			if !ok {
				break tobreake
			}
			eror = err

		default:
		}
	}

	return *message, eror
}

func (b MultiBot) HasReact(u *api.Update) bool {
	var hasReact bool
	for _, bot := range b {
		hasReact = hasReact && bot.HasReact(u)
	}
	return hasReact
}
