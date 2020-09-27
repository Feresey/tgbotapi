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

type ConversationState int

type Choice struct {
	// nil === true
	Accept func(*Message) bool
	Apply  func(context.Context, *Message) (ConversationState, error)
}

type Conversation struct {
	// read only
	states map[ConversationState][]Choice
	cache  *ttlcache.Cache
}

func NewConversation(cache *ttlcache.Cache) *Conversation {
	res := &Conversation{
		cache:  cache,
		states: make(map[ConversationState][]Choice),
	}
	return res
}

func (c *Conversation) Stop() {
	c.cache.Close()
}

// AddChoices to conversation list with given state.
// unsafe to use after starting a conversation.
func (c *Conversation) AddChoices(state ConversationState, choices ...Choice) {
	c.states[state] = append(c.states[state], choices...)
}

func keyFromUserID(userID int64) string {
	return strconv.FormatInt(userID, 10)
}

func (c *Conversation) AddUser(userID int64, state ConversationState) {
	c.cache.Set(keyFromUserID(userID), state)
}

func (c *Conversation) GetUserState(userID int64) (ConversationState, bool) {
	state, ok := c.cache.Get(keyFromUserID(userID))
	if !ok {
		return 0, false
	}
	return state.(ConversationState), true
}

func (c *Conversation) CheckUser(userID int64) bool {
	_, ok := c.cache.Get(keyFromUserID(userID))
	return ok
}

func (c *Conversation) RemoveUser(userID int64) {
	c.cache.Remove(keyFromUserID(userID))
}

func (c *Conversation) Handle(ctx context.Context, msg *Message) (ConversationState, error) {
	userID := msg.GetFrom().GetID()
	key := keyFromUserID(userID)
	stateI, ok := c.cache.Get(key)
	if !ok {
		return 0, ErrNoSuchConversation
	}
	state := stateI.(ConversationState)

	for _, choice := range c.states[state] {
		ok := choice.Accept == nil || choice.Accept(msg)
		if ok {
			next, err := choice.Apply(ctx, msg)
			if err == nil {
				c.cache.Set(key, next)
				return next, nil
			}
			return state, err
		}
	}

	return state, ErrNoSuchChoice
}
