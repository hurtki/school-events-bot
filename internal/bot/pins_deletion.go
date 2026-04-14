package bot

import tele "gopkg.in/telebot.v4"

func (b *Bot) DeletePinsHandle() {
	b.bot.Handle(tele.OnPinned, func(c tele.Context) error {
		pinnedPlateID := c.Message().ID + 1
		return b.bot.Delete(&tele.Message{ID: pinnedPlateID, Chat: &tele.Chat{ID: b.cfg.UpdatesChannel}})
	})
}
