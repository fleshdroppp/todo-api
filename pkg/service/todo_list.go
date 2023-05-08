package service

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	todo "todo-api"
	"todo-api/pkg/repository"
)

type TodoListService struct {
	repo  repository.TodoList
	cache repository.TodoListCache
}

func NewTodoListService(repo repository.TodoList, cache repository.TodoListCache) *TodoListService {
	return &TodoListService{repo: repo, cache: cache}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	lists := make([]todo.TodoList, 0)
	cache, err := s.cache.HGet(userId)
	if err == redis.Nil {
		logrus.Print("cache miss")
		lists, err = s.repo.GetAll(userId)
		if err != nil {
			return lists, err
		}
		data, err := json.Marshal(lists)
		err = s.cache.HSet(userId, string(data))
		return lists, err
	} else if err != nil {
		return lists, err
	}
	logrus.Print("cache hit")
	err = json.Unmarshal([]byte(cache), &lists)
	return lists, err
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
