package graph

import (
	"context"
	"database/sql"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/riyadennis/sigist/graphql-service/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var (
	errFailedDBOperation = errors.New("failed to perform db operation")
	logger               = otelzap.New(zap.NewNop())
)

type mockDB struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

type mockProducer struct {
}

func (m *mockProducer) Produce(_ *kafka.Message, _ chan kafka.Event) error {
	return nil
}

func TestMutationResolverSaveUser(t *testing.T) {
	scenarios := []struct {
		name        string
		in          *model.UserFeedbackInput
		out         *model.UserFeedback
		mockDB      *mockDB
		expectedErr error
	}{
		{
			name: "db prepare error",
			in: &model.UserFeedbackInput{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@doe.com",
				Feedback:  "This is a feedback",
			},
			mockDB:      mockUserSavePrepareError(t),
			expectedErr: errFailedDBOperation,
		},
		{
			name: "db exec error",
			in: &model.UserFeedbackInput{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@doe.com",
				Feedback:  "This is a feedback",
			},
			mockDB:      mockUserSaveStatementError(t),
			expectedErr: errFailedDBOperation,
		},
		{
			name: "db exec success",
			in: func() *model.UserFeedbackInput {
				jobTitle := "Quality Engineer"
				return &model.UserFeedbackInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@doe.com",
					Feedback:  "This is a feedback",
					JobTitle:  &jobTitle,
				}
			}(),
			mockDB: mockUserSaveStatementSuccess(t),
			out: func() *model.UserFeedback {
				ID := "123"
				firstName := "John"
				lastName := "Doe"
				email := "john@doe.com"
				jobTitle := "Quality Engineer"
				feedback := "This is a feedback"
				createAt := time.Now().Format(time.RFC3339)
				return &model.UserFeedback{
					ID:        &ID,
					FirstName: &firstName,
					LastName:  &lastName,
					Email:     &email,
					JobTitle:  &jobTitle,
					Feedback:  &feedback,
					CreateAt:  &createAt,
				}
			}(),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			resolver := &mutationResolver{
				Resolver: &Resolver{
					logger: logger,
					db:     scenario.mockDB.db,
					KafkaConfig: &KafkaConfig{
						Topic:    "test",
						Producer: &mockProducer{},
					},
				},
			}
			user, err := resolver.SaveUserFeedback(context.Background(), *scenario.in)
			assert.Equal(t, scenario.expectedErr, err)
			if user != nil {
				assert.Equal(t, *scenario.out.Feedback, *user.Feedback)
				assert.Equal(t, *scenario.out.CreateAt, *user.CreateAt)
			}
			err = scenario.mockDB.mock.ExpectationsWereMet()
			assert.NoError(t, err)

		})
	}
}

func TestQueryResolverGetUser(t *testing.T) {
	scenarios := []struct {
		name        string
		in          *model.FilterInput
		out         []*model.User
		mockDB      *mockDB
		expectedErr error
	}{
		{
			name: "db prepare error",
			in:   &model.FilterInput{},
			mockDB: func() *mockDB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnError(errFailedDBOperation)
				return &mockDB{
					db:   db,
					mock: mock,
				}
			}(),
			expectedErr: errFailedDBOperation,
		},
		{
			name: "db select error",
			in:   &model.FilterInput{},
			mockDB: func() *mockDB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"id", "first_name",
					"last_name", "email",
					"job_title", "created_at"}).AddRow(
					"INVALID", "John",
					"Doe", "john.doe@gmail.com",
					"Quality Engineer", time.Now().Format(time.RFC3339))
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
				return &mockDB{
					db:   db,
					mock: mock,
				}
			}(),
			expectedErr: errors.New("sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"INVALID\") to a int64: invalid syntax"),
		},
		{
			name: "db select success",
			in:   &model.FilterInput{},
			mockDB: func() *mockDB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"id", "first_name",
					"last_name", "email",
					"job_title", "created_at"}).AddRow(
					"1", "John",
					"Doe", "john.doe@gmail.com",
					"Quality Engineer", time.Now().Format(time.RFC3339))
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
				return &mockDB{
					db:   db,
					mock: mock,
				}
			}(),
			out: func() []*model.User {
				return []*model.User{
					{
						ID:        "1",
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john.doe@gmail.com",
						JobTitle:  "Quality Engineer",
						CreateAt:  time.Now().Format(time.RFC3339),
					},
				}
			}(),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			resolver := &queryResolver{
				Resolver: &Resolver{
					logger: logger,
					db:     scenario.mockDB.db,
				},
			}
			users, err := resolver.GetUserFeedback(context.Background(), *scenario.in)
			if err != nil {
				assert.Equal(t, scenario.expectedErr.Error(), err.Error())
			}

			assert.Equal(t, scenario.out, users)
			err = scenario.mockDB.mock.ExpectationsWereMet()
			assert.NoError(t, err)

		})
	}
}

func mockUserSavePrepareError(t *testing.T) *mockDB {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectPrepare("INSERT INTO").
		WillReturnError(errFailedDBOperation)

	return &mockDB{
		db:   db,
		mock: mock,
	}
}

func mockUserSaveStatementError(t *testing.T) *mockDB {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectPrepare("INSERT INTO").WillBeClosed()
	mock.ExpectExec("INSERT INTO").
		WillReturnError(errFailedDBOperation)

	return &mockDB{
		db:   db,
		mock: mock,
	}
}

func mockUserSaveStatementSuccess(t *testing.T) *mockDB {
	db, mock, err := sqlmock.New()
	result := sqlmock.NewResult(1, 1)
	assert.NoError(t, err)
	mock.ExpectPrepare("INSERT INTO").WillBeClosed()
	mock.ExpectExec("INSERT INTO").
		WillReturnResult(result)

	return &mockDB{
		db:   db,
		mock: mock,
	}
}
