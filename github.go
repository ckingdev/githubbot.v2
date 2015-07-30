package githubbot

import (
	"fmt"

	"github.com/cpalone/gobot"
	"github.com/cpalone/gohook"
)

func (c *CIListener) GithubListener(r *gobot.Room, port int, secret string, returnChan chan string) {
	gServer := gohook.NewServer(port, secret, "/postreceive")
	gServer.GoListenAndServe()
	for {
		et := <-gServer.EventAndTypes
		var msg string
		switch et.Type {
		case gohook.PingEventType:
			continue
		case gohook.CommitCommentEventType:
			payload, ok := et.Event.(*gohook.CommitCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			msg = fmt.Sprintf("[ %s ] Comment on commit: %s (%s)",
				payload.Repository.Name,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
		case gohook.CreateEventType:
			payload, ok := et.Event.(*gohook.CreateEvent)
			if !ok {
				panic("Malformed *CreateEvent.")
			}
			msg = fmt.Sprintf("[ %s | Branch/Tag: %s] Created.",
				payload.Repository.Name,
				payload.RefType,
			)
		case gohook.DeleteEventType:
			payload, ok := et.Event.(*gohook.DeleteEvent)
			if !ok {
				panic("Malformed *DeleteEvent.")
			}
			msg = fmt.Sprintf("[ %s | Branch/Tag: %s] Deleted.",
				payload.Repository,
				payload.RefType,
			)
		case gohook.IssueCommentEventType:
			payload, ok := et.Event.(*gohook.IssueCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			msg = fmt.Sprintf("[ %s | Issue: %s ] Comment: %s (%s)",
				payload.Repository.Name,
				payload.Issue.Title,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
		case gohook.IssuesEventType:
			payload, ok := et.Event.(*gohook.IssuesEvent)
			if !ok {
				panic("Malformed *IssuesEvent.")
			}
			msg = fmt.Sprintf("[ %s | Issue: %s ] Action: %s. (%s)",
				payload.Repository.Name,
				payload.Issue.Title,
				payload.Action,
				payload.Issue.HTMLURL,
			)
		case gohook.PullRequestEventType:
			payload, ok := et.Event.(*gohook.PullRequestEvent)
			if !ok {
				panic("Malformed *PullRequestEvent.")
			}
			action := payload.Action
			if action == "synced" {
				action = "New commits made to synced branch."
			}
			msg = fmt.Sprintf(":pencil: [ %s | PR: %s ] %s (%s)",
				payload.Repository.Name,
				payload.PullRequest.Title,
				action,
				payload.PullRequest.HTMLURL,
			)
		case gohook.PullRequestReviewCommentEventType:
			payload, ok := et.Event.(*gohook.PullRequestReviewCommentEvent)
			if !ok {
				panic("Malformed *PullRequestReviewCommentEvent.")
			}
			msg = fmt.Sprintf(":speech_balloon: [ %s | PR: %s ] Comment: %s: %s (%s)",
				payload.Repository.Name,
				payload.PullRequest.Title,
				payload.Sender.Login,
				payload.Comment.Body,
				payload.PullRequest.HTMLURL,
			)
		case gohook.RepositoryEventType:
			payload, ok := et.Event.(*gohook.RepositoryEvent)
			if !ok {
				panic("Malformed *RepositoryEvent.")
			}
			msg = fmt.Sprintf("[ Repository: %s ] Action: created. (%s) ",
				payload.Repository.Name,
				payload.Repository.HTMLURL,
			)
		case gohook.PushEventType:
			payload, ok := et.Event.(*gohook.PushEvent)
			if !ok {
				panic("Malformed *PushEvent.")
			}
			if payload.HeadCommit.Message == "" {
				return
			}
			if len(payload.Commits) > 1 {
				msg = fmt.Sprintf(":repeat: [ %s | Branch: %s ] %v Commits: %s (%s)",
					payload.Repository.Name,
					payload.Ref[11:], // this discards "refs/heads/"
					len(payload.Commits),
					payload.HeadCommit.Message,
					payload.Compare,
				)
			} else {
				msg = fmt.Sprintf(":repeat: [ %s | Branch: %s ] Commit: %s (%s)",
					payload.Repository.Name,
					payload.Ref[11:], // this discards "refs/heads/"
					payload.HeadCommit.Message,
					payload.HeadCommit.URL,
				)
			}
			msgID, err := r.SendText("", msg)
			if err != nil {
				continue
			}
			c.commitToMsgID[payload.HeadCommit.ID] = msgID
			continue
		}
		returnChan <- msg
	}
}

// func (c *CIListener) HandleIncoming(r *Room, p *proto.Packet) (*proto.Packet, error) {
// 	payload, err := p.Payload()
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch msg := payload.(type) {
// 	case proto.SendReply:
// 		msg.ID
// 	}
// }
