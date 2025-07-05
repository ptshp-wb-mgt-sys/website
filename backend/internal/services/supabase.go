// Package services/supabase.go contains the supabase service
// NOTE: Should this be here or in store internal/store/?
// NOTE: Should this even be called service or just db?
package services

import (
	"pet-mgt/backend/internal/config"

	"github.com/supabase-community/supabase-go"
)

type SupabaseService struct {
	config *config.Config
	client *supabase.Client
}

func NewSupabaseService(cfg *config.Config) (*SupabaseService, error) {
	client, err := supabase.NewClient(
		cfg.SupabaseURL,
		cfg.SupabaseKey,
		&supabase.ClientOptions{},
	)
	if err != nil {
		return nil, err
	}

	return &SupabaseService{
		config: cfg,
		client: client,
	}, nil
}
