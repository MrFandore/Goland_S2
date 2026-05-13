package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"example.com/Prac9-redis-cache/internal/cache"
	"example.com/Prac9-redis-cache/internal/config"
	"example.com/Prac9-redis-cache/internal/task"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	repo  *task.Repo
	redis *redis.Client
	cfg   config.Config
}

func NewTaskService(repo *task.Repo, redisClient *redis.Client, cfg config.Config) *TaskService {
	return &TaskService{
		repo:  repo,
		redis: redisClient,
		cfg:   cfg,
	}
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (task.Task, error) {
	key := cache.TaskByIDKey(id)

	// 1. Попытка чтения из кэша
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var t task.Task
			if err := json.Unmarshal([]byte(cached), &t); err == nil {
				log.Println("cache hit:", key)
				return t, nil
			}
			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	// 2. Чтение из основного хранилища (БД)
	t, err := s.repo.GetByID(id)
	if err != nil {
		return task.Task{}, err
	}

	// 3. Запись в кэш (если Redis доступен)
	if s.redis != nil {
		bytes, err := json.Marshal(t)
		if err != nil {
			log.Println("cache encode error:", err)
			return t, nil
		}
		ttl := cache.TTLWithJitter(s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
		if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
			log.Println("redis write error:", err)
		}
	}
	return t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t task.Task) error {
	if err := s.repo.Update(t); err != nil {
		return err
	}
	// Инвалидация кэша после обновления
	if s.redis != nil {
		key := cache.TaskByIDKey(t.ID)
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			log.Println("redis delete error:", err)
		}
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	// Инвалидация кэша после удаления
	if s.redis != nil {
		key := cache.TaskByIDKey(id)
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			log.Println("redis delete error:", err)
		}
	}
	return nil
}
