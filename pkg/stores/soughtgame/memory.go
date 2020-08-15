package soughtgame

import (
	"context"
	"errors"
	"sync"

	"github.com/domino14/liwords/pkg/entity"
	"github.com/rs/zerolog/log"
)

var (
	errNoID     = errors.New("sought game ID was not defined")
	errNotFound = errors.New("sought game ID was not found")
)

// MemoryStore is a purely in-memory store of a sought game. In the real final
// implementation, we will probably use a Postgres-backed store, or perhaps
// something in Redis.
type MemoryStore struct {
	sync.Mutex

	soughtGames             map[string]*entity.SoughtGame
	soughtGamesByUser       map[string]*entity.SoughtGame
	matchRequestsByReceiver map[string][]*entity.SoughtGame
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		soughtGames:             make(map[string]*entity.SoughtGame),
		soughtGamesByUser:       make(map[string]*entity.SoughtGame),
		matchRequestsByReceiver: make(map[string][]*entity.SoughtGame),
	}
}

// Get gets the game with the given ID.
func (m *MemoryStore) Get(ctx context.Context, id string) (*entity.SoughtGame, error) {
	g, ok := m.soughtGames[id]
	if !ok {
		return nil, errNotFound
	}
	return g, nil
}

// Set sets the game in the store.
func (m *MemoryStore) Set(ctx context.Context, game *entity.SoughtGame) error {
	m.Lock()
	defer m.Unlock()
	m.soughtGames[game.ID()] = game
	m.soughtGamesByUser[game.Seeker()] = game
	if game.Type() == entity.TypeMatch {
		if m.matchRequestsByReceiver[game.MatchRequest.ReceivingUser.UserId] == nil {
			m.matchRequestsByReceiver[game.MatchRequest.ReceivingUser.UserId] = []*entity.SoughtGame{}
		}

		m.matchRequestsByReceiver[game.MatchRequest.ReceivingUser.UserId] = append(
			m.matchRequestsByReceiver[game.MatchRequest.ReceivingUser.UserId], game)
	}

	log.Debug().Interface("by-user", m.soughtGamesByUser).Msg("set-sought-game")
	return nil
}

// Delete deletes the game by game ID.
func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()

	g, ok := m.soughtGames[id]
	if !ok {
		log.Warn().Str("game-id", id).Msg("tried-to-delete-nonexistent-game-id")
		return nil
	}

	userID := g.Seeker()
	delete(m.soughtGames, id)
	delete(m.soughtGamesByUser, userID)
	if g.Type() == entity.TypeMatch {
		m.deleteFromReqsByReceiver(g)
	}
	return nil
}

func (m *MemoryStore) deleteFromReqsByReceiver(g *entity.SoughtGame) {
	if len(m.matchRequestsByReceiver[g.MatchRequest.ReceivingUser.UserId]) == 1 {
		delete(m.matchRequestsByReceiver, g.MatchRequest.ReceivingUser.UserId)
		return
	}

	arrcpy := []*entity.SoughtGame{}
	for _, sg := range m.matchRequestsByReceiver[g.MatchRequest.ReceivingUser.UserId] {
		if sg.ID() != g.ID() {
			arrcpy = append(arrcpy, sg)
		}
	}
	m.matchRequestsByReceiver[g.MatchRequest.ReceivingUser.UserId] = arrcpy
}

// DeleteForUser deletes the game by seeker ID.
func (m *MemoryStore) DeleteForUser(ctx context.Context, userID string) (string, error) {
	game, ok := m.soughtGamesByUser[userID]
	if !ok {
		// Do nothing, game never existed
		return "", nil
	}
	m.Lock()
	defer m.Unlock()
	delete(m.soughtGamesByUser, userID)
	delete(m.soughtGames, game.ID())
	if game.Type() == entity.TypeMatch {
		m.deleteFromReqsByReceiver(game)
	}
	return game.ID(), nil
}

// ListOpenSeeks lists all open seek requests
func (m *MemoryStore) ListOpenSeeks(ctx context.Context) ([]*entity.SoughtGame, error) {
	ret := []*entity.SoughtGame{}
	for _, v := range m.soughtGames {
		if v.Type() == entity.TypeSeek {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (m *MemoryStore) ListOpenMatches(ctx context.Context, receiver string) ([]*entity.SoughtGame, error) {
	ret := []*entity.SoughtGame{}
	for _, v := range m.matchRequestsByReceiver[receiver] {
		if v.Type() != entity.TypeMatch {
			return nil, errors.New("unexpected type")
		}
		ret = append(ret, v)
	}
	return ret, nil
}

// ExistsForUser returns true if the user already has an outstanding seek request.
func (m *MemoryStore) ExistsForUser(ctx context.Context, userID string) (bool, error) {
	_, ok := m.soughtGamesByUser[userID]
	return ok, nil
}

// UserMatchedBy returns true if there is an open seek request from matcher for user
func (m *MemoryStore) UserMatchedBy(ctx context.Context, userID, matcher string) (bool, error) {
	matches, ok := m.matchRequestsByReceiver[userID]
	if !ok {
		return false, nil
	}
	for _, m := range matches {
		if m.MatchRequest.User.UserId == matcher {
			return true, nil
		}
	}
	return false, nil
}
