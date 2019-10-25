package oauth_router

import (
	"context"

	"github.com/go-session/session"
)

type MockMemoryStore struct {
	CheckFn   func(ctx context.Context, sid string) (bool, error)
	CreateFn  func(ctx context.Context, sid string, expired int64) (session.Store, error)
	UpdateFn  func(ctx context.Context, sid string, expired int64) (session.Store, error)
	DeleteFn  func(ctx context.Context, sid string) error
	RefreshFn func(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error)
	CloseFn   func() error
}

func (s *MockMemoryStore) Check(ctx context.Context, sid string) (bool, error) {
	return s.CheckFn(ctx, sid)
}

func (s *MockMemoryStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return s.CreateFn(ctx, sid, expired)
}

func (s *MockMemoryStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return s.UpdateFn(ctx, sid, expired)
}

func (s *MockMemoryStore) Delete(ctx context.Context, sid string) error {
	return s.DeleteFn(ctx, sid)
}

func (s *MockMemoryStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	return s.RefreshFn(ctx, oldsid, sid, expired)
}

func (s *MockMemoryStore) Close() error {
	return s.CloseFn()
}

type MockKeyValueStore struct {
	ContextFn   func() context.Context
	SessionIDFn func() string
	SetFn       func(key string, value interface{})
	GetFn       func(key string) (interface{}, bool)
	DeleteFn    func(key string) interface{}
	SaveFn      func() error
	FlushFn     func() error
}

func (s *MockKeyValueStore) Context() context.Context {
	return s.ContextFn()
}

func (s *MockKeyValueStore) SessionID() string {
	return s.SessionIDFn()
}

func (s *MockKeyValueStore) Set(key string, value interface{}) {
	s.SetFn(key, value)
}

func (s *MockKeyValueStore) Get(key string) (interface{}, bool) {
	return s.GetFn(key)
}

func (s *MockKeyValueStore) Delete(key string) interface{} {
	return s.DeleteFn(key)
}

func (s *MockKeyValueStore) Save() error {
	return s.SaveFn()
}

func (s *MockKeyValueStore) Flush() error {
	return s.FlushFn()
}
