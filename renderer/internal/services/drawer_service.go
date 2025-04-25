package services

import (
	"xlink/renderer/internal/ports"
	"xlink/renderer/internal/statistics_data"
)

type DrawerService struct {
	drawerRepo *ports.DrawerRepository
}

func NewDrawerService(dr ports.DrawerRepository) *DrawerService {
	return &DrawerService{drawerRepo: &dr}
}

func (s *DrawerService) Generate(input statistics_data.StatisticsData, paramName string) ([]byte, error) {
	return (*s.drawerRepo).Generate(input, paramName)
}
