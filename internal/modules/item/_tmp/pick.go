package service

import (
	"context"
	"fmt"
)

func (s *Service) Pick(ctx context.Context, weights []int64, n int) ([]int, error) {
	totalWeight := int64(0)

	for i, weight := range weights {
		if weight < 0 {
			return nil, fmt.Errorf("weight %d has negative value: %d", i, weight)
		}
		totalWeight += weight
	}

	if totalWeight <= 0 {
		return nil, fmt.Errorf("non-positive weight: %d", totalWeight)
	}

	results := make([]int, 0, n)
	for range n {
		result, err := s.singlePick(weights, totalWeight)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) singlePick(weights []int64, totalWeight int64) (int, error) {
	roll := s.rng.Int64() % totalWeight

	weightAccumulation := int64(0)
	for i, weight := range weights {
		weightAccumulation += weight
		if roll < weightAccumulation {
			return i, nil
		}
	}
	return 0, fmt.Errorf("unable to pick result for roll %v, total weight %v", roll, totalWeight)
}
