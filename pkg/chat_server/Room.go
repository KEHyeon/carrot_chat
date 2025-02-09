package chat_server

import (
	"sync"
)

type Chat struct {
	Id    uint64
	Users map[uint64]*User
	mu    sync.Mutex
}

func NewRoom(id uint64) *Chat {
	return &Chat{
		Id:    id,
		Users: make(map[uint64]*User),
		mu:    sync.Mutex{},
	}
}

func (r *Chat) DeleteUser(userId uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Users, userId)
	return nil
}
func (r *Chat) AddUser(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[user.GetId()] = user
	return nil
}
