package service

import (
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/stretchr/testify/assert"
	"testing"
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

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, err := NewService(scenario.cfg)
			assert.Equal(t, scenario.expectedErr, err)
		})
	}
}
