package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/alphadose/haxmap"
	"github.com/google/uuid"
)

type Queue struct {
	Chan    chan []byte
	ctx     context.Context
	cancel  context.CancelFunc
	Members *haxmap.Map[string, *Actor]
}

var queues *haxmap.Map[string, *Queue]

type Actor struct {
	UUID   string
	IP     string
	Conn   net.Conn
	Queues []string
}

func NewActor(ip string, conn net.Conn) *Actor {
	return &Actor{
		UUID:   uuid.NewString(),
		IP:     ip,
		Conn:   conn,
		Queues: []string{},
	}
}

func (a *Actor) AddQueue(item string) error {
	if a.Queues == nil {
		return errors.New("queues is nil")
	}
	queue, ok := queues.Get(item)
	if !ok {
		// Assign a queue and start a goroutine surrounding it's c ahnenl
		ctx, cancel := context.WithCancel(context.Background())
		queue = &Queue{
			Chan:    make(chan []byte, 1),
			ctx:     ctx,
			cancel:  cancel,
			Members: haxmap.New[string, *Actor](),
		}
		queues.Set(item, queue)
		go func(q *Queue) {
			for {
				select {
				case <-queue.ctx.Done():
					return
				case msg := <-queue.Chan:
					log.Println("Message received", string(msg))
					log.Println(q.Members.Len())
					q.Members.ForEach(func(s string, a *Actor) bool {
						if a.Conn != nil {
							a.Conn.Write(msg)
						}
						return true
					})
				}
			}
		}(queue)

	}

	queue.Members.Set(a.UUID, a)
	a.Queues = append(a.Queues, item)
	return nil

}

func (a *Actor) RemoveQueue(item string) error {
	if a.Queues == nil {
		return errors.New("queues is nil")
	}
	queue, ok := queues.Get("item")
	if !ok {
		return errors.New("Queue not found")
	}
	queue.Members.Del(a.UUID)
	a.Queues = append(a.Queues, item)
	// if members is 0 dissipate the queue
	if queue.Members.Len() == 0 {
		queue.cancel()
		queues.Del(item)
	}
	return nil
}

func DistributeMessage(intake string, message []byte) {
	queue, ok := queues.Get(intake)
	if !ok {
		return
	}
	queue.Chan <- message

}

func handleConnection(c net.Conn, actor *Actor) {
	defer func() {
		c.Close()
		for _, queue_ := range actor.Queues {
			queue, ok := queues.Get(queue_)
			if !ok {
				continue
			}
			queue.Members.Del(actor.UUID)
			if queue.Members.Len() == 0 {
				queue.cancel()
				queues.Del(queue_)
			}
		}
	}()
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			return
		}
		// Json decode
		var msg map[string]interface{}
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			continue
		}
		if msg["type"] == nil {
			continue
		}
		switch msg["type"] {
		case "join":
			if msg["queue"] == nil {
				continue
			}
			actor.AddQueue(msg["queue"].(string))
		case "leave":
			if msg["queue"] == nil {
				continue
			}
			actor.RemoveQueue(msg["queue"].(string))
		case "message":
			if msg["queue"] == nil {
				continue
			}
			if msg["message"] == nil {
				continue
			}
			DistributeMessage(msg["queue"].(string), []byte(msg["message"].(string)))
		default:
			continue
		}
	}
}

func main() {
	queues = haxmap.New[string, *Queue]()

	// Listen on port 2000
	ln, err := net.Listen("tcp", ":2000")
	if err != nil {
		panic(err)
	}

	log.Println("Server started on port 2000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		actor := NewActor(conn.RemoteAddr().String(), conn)
		go handleConnection(conn, actor)
	}
}
