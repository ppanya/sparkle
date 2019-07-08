package endpoints

import (
	"context"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
)

type SpikeEndpoints struct {
}

func NewSpikeEndpoints() *SpikeEndpoints {
	return &SpikeEndpoints{}
}

func (s *SpikeEndpoints) GetUserByAccessToken(context.Context, *svcsv1.GetUserByAccessTokenInput) (*svcsv1.GetUserByAccessTokenOutput, error) {
	panic("implement me")
}
