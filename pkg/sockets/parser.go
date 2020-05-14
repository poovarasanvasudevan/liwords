package sockets

import (
	"context"
	"errors"

	macondopb "github.com/domino14/macondo/gen/api/proto/macondo"

	"github.com/domino14/crosswords/pkg/entity"
	"github.com/domino14/crosswords/pkg/game"
	pb "github.com/domino14/crosswords/rpc/api/proto"
)

func (h *Hub) parseAndExecuteMessage(msg []byte, sender string) error {
	// All socket messages are encoded entity.Events.
	// (or they better be)

	ew, err := entity.EventFromByteArray(msg)
	if err != nil {
		return err
	}
	switch ew.Name {

	case "crosswords.GameAcceptedEvent":
		evt, ok := ew.Event.(*pb.GameAcceptedEvent)
		if !ok {
			return errors.New("unexpected typing error")
		}
		// XXX: We're going to want to fetch these player IDs from the database later.
		players := []*macondopb.PlayerInfo{
			&macondopb.PlayerInfo{Nickname: evt.Acceptor, RealName: evt.Acceptor},
			&macondopb.PlayerInfo{Nickname: evt.Requester, RealName: evt.Requester},
		}

		g, err := game.InstantiateNewGame(context.Background(), h.gameStore, h.config,
			players, evt.GameRequest)
		if err != nil {
			return err
		}
		// Create a "realm" for these users, and add both of them to the realm.
		// Use the id of the game as the id of the realm.
		realm := h.addNewRealm(g.GameID())
		h.addToRealm(realm, evt.Acceptor)
		h.addToRealm(realm, evt.Requester)
		// Now, we start the timer and register the event hook.
		err = game.StartGameInstance(g, h.eventChan)
		if err != nil {
			return err
		}

	case "crosswords.UserGameplayEvent":
		evt, ok := ew.Event.(*pb.UserGameplayEvent)
		if !ok {
			// This really shouldn't happen
			return errors.New("unexpected typing error")
		}
		err := game.PlayMove(context.Background(), h.gameStore, sender, evt)
		if err != nil {
			return err
		}

	default:
		return errors.New("evt " + ew.Name + " not yet handled")

	}

	return nil
}