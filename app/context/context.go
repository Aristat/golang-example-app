package context

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

const prefix = "app.context"

// Manager
type Manager struct {
	mapping *Mapping
	ctx     context.Context
}

func (m *Manager) ToMapping() *Mapping {
	return &Mapping{
		Subject:     m.mapping.Subject,
		ServiceId:   m.mapping.ServiceId,
		ServiceName: m.mapping.ServiceName,
	}
}

// NewManager
func NewManager(ctx context.Context) (*Manager, error) {
	m := &Manager{ctx: ctx}
	mapping := &Mapping{}

	config := &mapstructure.DecoderConfig{TagName: "json", Result: &mapping}
	decoder, _ := mapstructure.NewDecoder(config)
	err := decoder.Decode(ctx.Value(prefix))
	if err != nil {
		return nil, err
	}

	m.mapping = mapping
	return m, nil
}

// Mapping
type Mapping struct {
	Subject     string
	ServiceId   uint64 `json:"service_id"`
	ServiceName string `json:"service_name"`
}

// NewContext
func NewContext(ctx context.Context, m Mapping) context.Context {
	return context.WithValue(ctx, prefix, map[string]interface{}{
		"subject":      m.Subject,
		"service_id":   m.ServiceId,
		"service_name": m.ServiceName,
	})
}
