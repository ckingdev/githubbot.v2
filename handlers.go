package githubbot

import (
	"euphoria.io/heim/proto"
	"github.com/cpalone/gobot"
)

type CIListener struct {
	packetIDtoCommit map[string]string
	commitToMsgID    map[string]string
	port             int
	secret           string
}

func (c *CIListener) HandleIncoming(r *gobot.Room, p *proto.Packet) (*proto.Packet, error) {
	payload, err := p.Payload()
	if err != nil {
		return nil, err
	}
	switch msg := payload.(type) {
	case proto.SendReply:
		commit, ok := c.packetIDtoCommit[p.ID]
		if !ok {
			return nil, nil
		}
		c.commitToMsgID[commit] = msg.ID.String()
	default:
		return nil, nil
	}
	return nil, nil
}

func (c *CIListener) Run(r *gobot.Room) {
	ghReturn := make(chan string)
	go c.GithubListener(r, c.port, c.secret, ghReturn)
	for {
		select {
		case msg <- ghReturn:
			_, err := r.SendText("", msg)
			if err != nil {
				r.Logger.Warning(err)
			}
		}
	}
}
