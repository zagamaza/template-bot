package bot

import (
	"context"
	"github.com/almaznur91/splitty/internal/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatStateService interface {
	Save(ctx context.Context, u *api.ChatState) error
	DeleteById(ctx context.Context, id primitive.ObjectID) error
	FindByUserId(ctx context.Context, userId int) (*api.ChatState, error)
	CleanChatState(ctx context.Context, state *api.ChatState)
}

type ButtonService interface {
	Save(ctx context.Context, u *api.Button) (primitive.ObjectID, error)
	SaveAll(ctx context.Context, b ...*api.Button) ([]*api.Button, error)
}

//StartScreen send /room, after click on the button 'Присоединиться'
type StartScreen struct {
	css ChatStateService
	bs  ButtonService
	us  UserService
	cfg *Config
}

// NewStartScreen makes a bot for SO
func NewStartScreen(s ChatStateService, bs ButtonService, us UserService, cfg *Config) *StartScreen {
	return &StartScreen{
		css: s,
		bs:  bs,
		us:  us,
		cfg: cfg,
	}
}

// StartScreen keys
func (s StartScreen) HasReact(u *api.Update) bool {
	if hasAction(u, viewStart) {
		return true
	} else if isPrivate(u) {
		return u.Message != nil && u.Message.Text == start
	} else {
		return u.Message != nil && u.Message.Text == start+"@"+s.cfg.BotName
	}
}

// OnMessage returns one entry
func (s *StartScreen) OnMessage(ctx context.Context, u *api.Update) (api.TelegramMessage, error) {
	defer s.css.CleanChatState(ctx, u.ChatState)

	var screen tgbotapi.Chattable
	cb := api.NewButton(createRoom, new(api.CallbackData))
	if _, err := s.bs.SaveAll(ctx, cb); err != nil {
		return api.TelegramMessage{}, err
	}
	screen = createScreen(u, I18n(u.User, "scrn_main"), &[][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(I18n(u.User, "btn_create_room"), cb.ID.Hex())},
	})

	return api.TelegramMessage{
		Chattable: []tgbotapi.Chattable{screen},
		Send:      true,
	}, nil
}
