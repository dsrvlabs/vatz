package service

import (
	"time"
)

type service struct {
	CheckInterval time.Duration
	workerChannel chan bool
}

func (s *service) Init() error {
	return nil
}

func (s *service) Execute() error {
	return nil
}

func (s *service) End() error {
	return nil
}

func (s *service) Verify() error {
	return nil
}

func (s *service) UpdateConfig() error {
	return nil
}
