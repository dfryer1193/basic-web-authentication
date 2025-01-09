package storage

import (
	"github.com/dfryer1193/basic-web-authentication/models"
	"sync"
)

type InMemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{users: make(map[string]models.User)}
}

func (store *InMemoryUserStore) Get(username string) (models.User, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	user, exists := store.users[username]
	return user, exists
}

func (store *InMemoryUserStore) Set(username string, user models.User) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.users[username] = user
}
