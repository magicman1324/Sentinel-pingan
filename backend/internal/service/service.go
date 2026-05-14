package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/pingan/monitor-backend/internal/model"
	"github.com/pingan/monitor-backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	ruleRepo  *repository.RuleRepository
	alertRepo *repository.AlertRepository
	rdb       *redis.Client
	mu        sync.Mutex // guards SyncRulesToRedis
}

func NewService(rr *repository.RuleRepository, ar *repository.AlertRepository, rdb *redis.Client) *Service {
	return &Service{ruleRepo: rr, alertRepo: ar, rdb: rdb}
}

func (s *Service) GetRules(ctx context.Context) ([]model.Rule, error) {
	return s.ruleRepo.ListAll()
}

func (s *Service) GetEnabledRules(ctx context.Context) ([]model.Rule, error) {
	return s.ruleRepo.ListEnabled()
}

func (s *Service) GetRuleByID(ctx context.Context, id int64) (*model.Rule, error) {
	return s.ruleRepo.GetByID(id)
}

func (s *Service) CreateRule(ctx context.Context, rule *model.Rule) error {
	if err := s.ruleRepo.Create(rule); err != nil {
		return err
	}
	return s.SyncEnabledRulesToRedis(ctx)
}

func (s *Service) UpdateRule(ctx context.Context, rule *model.Rule) error {
	_, err := s.ruleRepo.GetByID(rule.ID)
	if err != nil {
		return ErrNotFound
	}
	if err := s.ruleRepo.Update(rule); err != nil {
		return err
	}
	return s.SyncEnabledRulesToRedis(ctx)
}

func (s *Service) DeleteRule(ctx context.Context, id int64) error {
	if err := s.ruleRepo.Delete(id); err != nil {
		return err
	}
	return s.SyncEnabledRulesToRedis(ctx)
}

func (s *Service) SyncEnabledRulesToRedis(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules, err := s.ruleRepo.ListEnabled()
	if err != nil {
		return err
	}
	data, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, "monitor:rules", data, 0).Err()
}

func (s *Service) GetAlerts(ctx context.Context, page, size int) ([]model.Alert, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.alertRepo.List(size, offset)
}

func (s *Service) ResolveAlert(ctx context.Context, id int64) error {
	return s.alertRepo.Resolve(id)
}
