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
	errFailedDBOperation      = errors.New("failed to perform db operation")
	errFailedToPublishToKafka = errors.New("failed to produce kafka message")
	logger                    = otelzap.New(zap.NewNop())
	id                        = "123"
	createdAt                 = time.Now().Format(time.RFC3339)
	firstName                 = "John"
	lastName                  = "Doe"
	email                     = "john@test.com"
	feedback                  = "This is a feedback"
	jobTitle                  = "Software Engineer"
)

type mockDB struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

type mockProducer struct {
	err error
}

func (m *mockProducer) Produce(_ *kafka.Message, _ chan kafka.Event) error {
	return m.err
}

func TestMutationResolverSaveUserFeedback(t *testing.T) {
	scenarios := []struct {
		name        string
		in          *model.UserFeedbackInput
		out         *model.UserFeedback
		mockDB      *mockDB
		mockKafka   *mockProducer
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
			name: "db exec success kafka error",
			in: func() *model.UserFeedbackInput {
				return &model.UserFeedbackInput{
					FirstName: firstName,
					LastName:  lastName,
					Email:     email,
					Feedback:  feedback,
					JobTitle:  &jobTitle,
				}
			}(),
			mockDB: mockUserSaveStatementSuccess(t),
			mockKafka: &mockProducer{
				err: errFailedToPublishToKafka,
			},
			out: func() *model.UserFeedback {
				return &model.UserFeedback{
					ID:        &id,
					FirstName: &firstName,
					LastName:  &lastName,
					Email:     &email,
					JobTitle:  &jobTitle,
					Feedback:  &feedback,
					CreateAt:  &createdAt,
				}
			}(),
			expectedErr: errFailedToPublishToKafka,
		},
		{
			name: "success",
			in: func() *model.UserFeedbackInput {
				return &model.UserFeedbackInput{
					FirstName: firstName,
					LastName:  lastName,
					Email:     email,
					Feedback:  feedback,
					JobTitle:  &jobTitle,
				}
			}(),
			mockDB:    mockUserSaveStatementSuccess(t),
			mockKafka: &mockProducer{},
			out: func() *model.UserFeedback {
				return &model.UserFeedback{
					ID:        &id,
					FirstName: &firstName,
					LastName:  &lastName,
					Email:     &email,
					JobTitle:  &jobTitle,
					Feedback:  &feedback,
					CreateAt:  &createdAt,
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
						Producer: scenario.mockKafka,
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
		out         []*model.UserFeedback
		mockDB      *mockDB
		expectedErr error
	}{
		{
			name: "db prepare error",
			in:   &model.FilterInput{},
			mockDB: func() *mockDB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(queryGetAllUsers).
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
				mock.ExpectQuery(queryGetAllUsers).WillReturnError(errFailedDBOperation)
				return &mockDB{
					db:   db,
					mock: mock,
				}
			}(),
			expectedErr: errors.New("failed to perform db operation"),
		},
		{
			name: "db select success",
			in:   &model.FilterInput{},
			mockDB: func() *mockDB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"id", "first_name",
					"last_name", "email",
					"job_title", "feedback", "created_at"}).AddRow(
					"1", "John",
					"Doe", "john.doe@gmail.com",
					"Quality Engineer", "loved it", time.Now().Format(time.RFC3339))
				mock.ExpectQuery(queryGetAllUsers).
					WillReturnRows(rows)
				return &mockDB{
					db:   db,
					mock: mock,
				}
			}(),
			out: func() []*model.UserFeedback {
				return []*model.UserFeedback{
					{
						ID:        &id,
						FirstName: &firstName,
						LastName:  &lastName,
						Email:     &email,
						JobTitle:  &jobTitle,
						CreateAt:  &createdAt,
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
			if len(users) > 0 {
				assert.Equal(t, *scenario.out[0].FirstName, *users[0].FirstName)
			}
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
