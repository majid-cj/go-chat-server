package util

import (
	"fmt"
)

func GetChatId(url string, isReceiver bool) string {
	sender := GetURLIds(url)[0]
	receiver := GetURLIds(url)[1]
	if isReceiver {
		return fmt.Sprintf("%s-%s", receiver, sender)
	}
	return fmt.Sprintf("%s-%s", sender, receiver)
}
