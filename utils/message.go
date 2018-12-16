package utils

import (
	"strings"

	"github.com/jerray/chatroom/pb"
)

// ParseInput parse user input.
//
// * `#login [name]` register a user name
// * `@[name] [message]` send message to a user
// * `[message]` send message to all users
func ParseInput(text string) *pb.Event {
	switch true {
	case strings.HasPrefix(text, "#"):
		args := strings.Split(text, " ")
		return Command(args[0], args[1:])
	case strings.HasPrefix(text, "@"):
		args := strings.SplitN(text, " ", 2)
		if len(args) < 2 {
			return nil
		}
		return Message(args[0], args[1])
	}

	return Message("", text)
}

func Command(command string, args []string) *pb.Event {
	switch command {
	case "#login":
		return &pb.Event{
			Command: &pb.Event_Login{
				Login: &pb.Client{
					Name: args[0],
				},
			},
		}
	}
	return nil
}

func Message(to, text string) *pb.Event {
	message := &pb.Message{
		Content: text,
	}
	if to != "" {
		message.To = &pb.Client{
			Name: strings.TrimPrefix(to, "@"),
		}
	}
	return &pb.Event{
		Command: &pb.Event_Message{
			Message: message,
		},
	}
}
