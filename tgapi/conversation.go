package tgapi

import (
	"context"
	"errors"
	"strconv"

	"github.com/ReneKroon/ttlcache"
)

var (
	ErrNoSuchConversation = errors.New("no such conversation")
	ErrNoSuchChoice       = errors.New("no such choice")
)

type State int

const StateFinished State = 0

type Choice struct {
	// nil === true
	Accept func(*Message) bool
	Apply  func(context.Context, *Message) (State, error)
}

type Conversation struct {
	// read only
	states map[State][]Choice

	api   *API
	cache *ttlcache.Cache
}

func NewConversation(api *API, cache *ttlcache.Cache) *Conversation {
	res := &Conversation{
		api:    api,
		cache:  cache,
		states: make(map[State][]Choice),
	}
	return res
}

func (c *Conversation) Stop() {
	c.cache.Close()
}

// AddChoices to conversation list with given state.
// unsafe to use after starting a conversation.
func (c *Conversation) AddChoices(state State, choices ...Choice) {
	c.states[state] = append(c.states[state], choices...)
}

func keyFromUserID(userID int64) string {
	return strconv.FormatInt(userID, 10)
}

func (c *Conversation) AddUser(userID int64, state State) {
	c.cache.Set(keyFromUserID(userID), state)
}

func (c *Conversation) CheckUser(userID int64) bool {
	_, ok := c.cache.Get(keyFromUserID(userID))
	return ok
}

func (c *Conversation) RemoveUser(userID int64) {
	c.cache.Remove(keyFromUserID(userID))
}

func (c *Conversation) Handle(ctx context.Context, msg *Message) error {
	userID := msg.GetFrom().GetID()
	key := keyFromUserID(userID)
	state, ok := c.cache.Get(key)
	if !ok {
		return ErrNoSuchConversation
	}

	for _, choice := range c.states[state.(State)] {
		ok := choice.Accept == nil || choice.Accept(msg)
		if ok {
			next, err := choice.Apply(ctx, msg)
			if err == nil {
				if next == StateFinished {
					c.cache.Remove(key)
				} else {
					c.cache.Set(key, next)
				}
			}
			return err
		}
	}

	return ErrNoSuchChoice
}
