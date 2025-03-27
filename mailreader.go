package main

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func GetLastEmail(username, password string) (string, error) {
	// Подключение к IMAP серверу Mail.ru
	c, err := client.DialTLS("imap.mail.ru:993", nil)
	if err != nil {
		return "", fmt.Errorf("не удалось подключиться: %v", err)
	}
	defer c.Logout()

	// Авторизация
	if err := c.Login(username, password); err != nil {
		return "", fmt.Errorf("не удалось авторизоваться: %v", err)
	}

	// Открытие INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return "", fmt.Errorf("не удалось открыть INBOX: %v", err)
	}

	if mbox.Messages == 0 {
		return "Нет писем", nil
	}

	// Получение последнего письма
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	messages := make(chan *imap.Message, 1)
	err = c.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении письма: %v", err)
	}

	msg := <-messages
	if msg == nil {
		return "Письмо не найдено", nil
	}

	from := msg.Envelope.From[0]
	return fmt.Sprintf("От: %s@%s\nТема: %s", from.MailboxName, from.HostName, msg.Envelope.Subject), nil
}
