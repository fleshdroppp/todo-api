package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	redis2 "github.com/redis/go-redis/v9"
	todo "todo-api"
	"todo-api/pkg/repository/postgres"
	"todo-api/pkg/repository/redis"
)

// postgres

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

// redis

type TodoListCache interface {
	HSet(userId int, data string) error
	HSetById(userId int, listId int, data string) error
	HGet(userId int) (string, error)
	HGetById(userId int, listId int) (string, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoListCache
	TodoItem
}

func NewRepository(ctx context.Context, db *sqlx.DB, client *redis2.Client) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
		TodoList:      postgres.NewTodoListPostgres(db),
		TodoListCache: redis.NewTodoListRedis(ctx, client),
		TodoItem:      postgres.NewTodoItemPostgres(db),
	}
}
