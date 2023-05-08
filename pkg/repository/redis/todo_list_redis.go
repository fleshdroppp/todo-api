package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type TodoListRedis struct {
	ctx    context.Context
	client *redis.Client
}

func NewTodoListRedis(ctx context.Context, client *redis.Client) *TodoListRedis {
	return &TodoListRedis{ctx: ctx, client: client}
}

func (r *TodoListRedis) HSet(userId int, data string) error {
	pipeline := r.client.Pipeline()
	pipeline.HSetNX(r.ctx, fmt.Sprintf("user:%d", userId), "lists", data)
	pipeline.Expire(r.ctx, fmt.Sprintf("user:%d", userId), duration)
	_, err := pipeline.Exec(r.ctx)
	return err
}

func (r *TodoListRedis) HSetById(userId int, listId int, data string) error {
	pipeline := r.client.Pipeline()
	pipeline.HSetNX(r.ctx, fmt.Sprintf("user:%d", userId), fmt.Sprintf("list:%d", listId), data)
	pipeline.Expire(r.ctx, fmt.Sprintf("user:%d", userId), duration)
	_, err := pipeline.Exec(r.ctx)
	return err
}

func (r *TodoListRedis) HGet(userId int) (string, error) {
	result, err := r.client.HGet(r.ctx, fmt.Sprintf("user:%d", userId), "lists").Result()
	return result, err
}

func (r *TodoListRedis) HGetById(userId int, listId int) (string, error) {
	result, err := r.client.HGet(r.ctx, fmt.Sprintf("user:%d", userId), fmt.Sprintf("list:%d", listId)).Result()
	return result, err
}
