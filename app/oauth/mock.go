package oauth

import (
	"context"

	"github.com/go-session/session"
)

type MemoryStore struct {
	CheckFn   func(ctx context.Context, sid string) (bool, error)
	CreateFn  func(ctx context.Context, sid string, expired int64) (session.Store, error)
	UpdateFn  func(ctx context.Context, sid string, expired int64) (session.Store, error)
	DeleteFn  func(ctx context.Context, sid string) error
	RefreshFn func(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error)
	CloseFn   func() error
}

func (s *MemoryStore) Check(ctx context.Context, sid string) (bool, error) {
	return s.CheckFn(ctx, sid)
}

func (s *MemoryStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return s.CreateFn(ctx, sid, expired)
}

func (s *MemoryStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return s.UpdateFn(ctx, sid, expired)
}

func (s *MemoryStore) Delete(ctx context.Context, sid string) error {
	return s.DeleteFn(ctx, sid)
}

func (s *MemoryStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	return s.RefreshFn(ctx, oldsid, sid, expired)
}

func (s *MemoryStore) Close() error {
	return s.CloseFn()
}

type KeyValueStore struct {
	ContextFn   func() context.Context
	SessionIDFn func() string
	SetFn       func(key string, value interface{})
	GetFn       func(key string) (interface{}, bool)
	DeleteFn    func(key string) interface{}
	SaveFn      func() error
	FlushFn     func() error
}

func (s *KeyValueStore) Context() context.Context {
	return s.ContextFn()
}

func (s *KeyValueStore) SessionID() string {
	return s.SessionIDFn()
}

func (s *KeyValueStore) Set(key string, value interface{}) {
	s.SetFn(key, value)
}

func (s *KeyValueStore) Get(key string) (interface{}, bool) {
	return s.GetFn(key)
}

func (s *KeyValueStore) Delete(key string) interface{} {
	return s.DeleteFn(key)
}

func (s *KeyValueStore) Save() error {
	return s.SaveFn()
}

func (s *KeyValueStore) Flush() error {
	return s.FlushFn()
}
