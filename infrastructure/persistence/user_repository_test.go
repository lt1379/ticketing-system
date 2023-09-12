package persistence

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"my-project/domain/model"
	"my-project/domain/repository"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository repository.IUser
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.repository = NewUserRepository(db)
}

func (s *Suite) TestUserRepository_GetByIdSQLMock() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	createdAtTime, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")
	updatedAtTime, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")
	var (
		ID        = 1
		Name      = "Lambok Tulus Simamora"
		UserName  = "lamboktulus1379"
		Password  = "a252f77af72638ea5a0f9e5fbe5f2b2e"
		CreatedAt = createdAtTime.In(loc)
		UpdatedAt = updatedAtTime.In(loc)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.name, u.user_name, u.password, u.created_at, u.updated_at 
	FROM public.user AS u 
	WHERE u.id = $1`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_name", "password", "created_at", "updated_at"}).
			AddRow(ID, Name, UserName, Password, CreatedAt, UpdatedAt))

	res, err := s.repository.GetById(context.Background(), 1)
	exp := model.User{
		ID:        1,
		Name:      "Lambok Tulus Simamora",
		UserName:  "lamboktulus1379",
		Password:  "a252f77af72638ea5a0f9e5fbe5f2b2e",
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}

	require.NoError(s.T(), err)
	require.Equal(s.T(), exp, res)

}

func (s *Suite) TestUserRepository_GetByIdSQLMockErr() {
	s.mock.ExpectPrepare(`SELECT u.id, u.name, u.user_name, u.password, u.created_at, u.updated_at 
	FROM public.user AS u 
	WHERE u.id = $1`).
		WillReturnError(fmt.Errorf("error statement"))

	res, err := s.repository.GetById(context.Background(), 1)
	exp := model.User{}

	require.Error(s.T(), err)
	require.Equal(s.T(), exp, res)

}

func TestUserRepository_GetById(t *testing.T) {
	DB, err := NewPostgreSQLDb()
	if err != nil {
		t.Fatalf("Error instantiate database. %v", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("Error while load location. %v", err)
	}
	createdAt, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")
	updatedAt, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")

	createdAt = createdAt.In(loc)
	updatedAt = updatedAt.In(loc)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser model.User
		wantErr  bool
	}{
		{
			name: "Test Success #1",
			fields: fields{
				DB: DB,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantUser: model.User{
				ID:        1,
				Name:      "Lambok Tulus Simamora",
				UserName:  "lamboktulus1379",
				Password:  "a252f77af72638ea5a0f9e5fbe5f2b2e",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				sqlDB: tt.fields.DB,
			}
			gotUser, err := repo.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserRepository.GetById() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, _ := NewPostgreSQLDb()
	pw := []byte("MyPassword_123")
	type fields struct {
		sqlDB *sql.DB
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Create User #1",
			fields: fields{
				sqlDB: db,
			},
			args: args{
				ctx: context.Background(),
				user: model.User{
					Name:      "Lambok Tulus Simamora",
					UserName:  "lamboktulus1379",
					Email:     "lamboktulus1379@gmail.com",
					Password:  fmt.Sprintf("%x", md5.Sum(pw)),
					CreatedBy: 0,
					UpdatedBy: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &UserRepository{
				sqlDB: tt.fields.sqlDB,
			}
			if err := userRepository.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_GetByUserName(t *testing.T) {
	DB, err := NewPostgreSQLDb()
	if err != nil {
		t.Fatalf("Error instantiate database. %v", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("Error while load location. %v", err)
	}
	createdAt, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")
	updatedAt, _ := time.Parse("2006-01-02 15:04:05.999999999+07 MST", "2023-09-04 01:02:10.911651+07 WIB")

	createdAt = createdAt.In(loc)
	updatedAt = updatedAt.In(loc)
	if err != nil {
		t.Fatalf("Error connection to database. %v", err)
	}
	type fields struct {
		sqlDB *sql.DB
	}
	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.User
		wantErr bool
	}{
		{
			name: "Test Success #1",
			fields: fields{
				sqlDB: DB,
			},
			args: args{
				ctx:      context.Background(),
				userName: "lamboktulus1379",
			},
			want: model.User{
				ID:        1,
				Name:      "Lambok Tulus Simamora",
				UserName:  "lamboktulus1379",
				Password:  "a252f77af72638ea5a0f9e5fbe5f2b2e",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &UserRepository{
				sqlDB: tt.fields.sqlDB,
			}
			got, err := userRepository.GetByUserName(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetByUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetByUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
