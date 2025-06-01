package message

type Message struct {
	Content          string
	MentionsEveryone bool
}

func NewMessage(content string, mentionsEveryone bool) *Message {
	return &Message{
		Content:          content,
		MentionsEveryone: mentionsEveryone,
	}
}
