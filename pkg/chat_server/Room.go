package chat_server

import (
	"sync"
)

type Room struct {
	Id    uint64
	Users map[uint64]*User
	mu    sync.Mutex
}

func NewRoom(id uint64) *Room {
	return &Room{
		Id:    id,
		Users: make(map[uint64]*User),
		mu:    sync.Mutex{},
	}
}

func (r *Room) DeleteUser(userId uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Users, userId)
	return nil
}
func (r *Room) AddUser(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[user.GetId()] = user
	return nil
}
