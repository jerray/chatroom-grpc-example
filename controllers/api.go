package controllers

import (
	"io"
	"log"
	"sync"

	"github.com/jerray/chatroom/pb"
)

type RoomController struct {
	clients sync.Map
}

type ClientStream struct {
	Stream pb.Chatroom_CheckInServer
	Name   string
}

func NewRoomController() *RoomController {
	return &RoomController{
		clients: sync.Map{},
	}
}

// CheckIn starts a new goroutine to handle a stream with a client.
func (r *RoomController) CheckIn(stream pb.Chatroom_CheckInServer) error {
	for {
		in, err := stream.Recv()
		ctx := stream.Context()
		clientID := ExtractClientID(ctx)
		var name string

		if err == io.EOF {
			r.clients.Delete(clientID)
			log.Println("client", clientID, "disconnected")
			return nil
		}
		if err != nil {
			r.clients.Delete(clientID)
			log.Println("client", clientID, "lost connection")
			return err
		}

		// store all clients in a sync.Map
		if client, ok := r.clients.Load(clientID); ok {
			name = client.(*ClientStream).Name
		} else {
			r.clients.Store(clientID, &ClientStream{
				Stream: stream,
			})
			log.Println("client", clientID, "connected")
		}

		from := &pb.Client{
			ClientId: clientID,
			Name:     name,
		}

		if login := in.GetLogin(); login != nil {
			name = login.GetName()
			r.clients.Store(clientID, &ClientStream{
				Stream: stream,
				Name:   name,
			})
			from.Name = name
			log.Println(name, "logged in")
			r.Broadcast(from, nil, &pb.Event{
				Command: &pb.Event_Login{
					Login: from,
				},
			})
		} else if message := in.GetMessage(); message != nil {
			content := message.GetContent()
			log.Println("client", clientID, "send message", content)
			r.Broadcast(from, message.GetTo(), &pb.Event{
				Command: &pb.Event_Message{
					Message: &pb.Message{
						From:    from,
						Content: content,
					},
				},
			})
		}
	}
}

// Braodcast route event to all clients except the sender. if parameter `to` is
// provided, only route event to clients match the name.
func (r *RoomController) Broadcast(from *pb.Client, to *pb.Client, event *pb.Event) {
	r.clients.Range(func(k interface{}, v interface{}) bool {
		key := k.(string)
		if key == from.ClientId {
			return true
		}
		client := v.(*ClientStream)

		if to != nil && to.Name != client.Name {
			return true
		}

		client.Stream.Send(event)
		return true
	})
}
