package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/event-management/graphql-service/internal"
)

func TestNewService(t *testing.T) {
	scenarios := []struct {
		name        string
		cfg         internal.Config
		expectedErr error
	}{
		{
			name:        "should return error when migration fails",
			cfg:         internal.Config{Env: "test"},
			expectedErr: ErrFailedTORunMigration,
		},
	}
	ctx := context.Background()
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, err := NewService(ctx, scenario.cfg)
			assert.Equal(t, scenario.expectedErr, err)
		})
	}
}
