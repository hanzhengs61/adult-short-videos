package logic

import (
	"context"

	"adult-short-videos/internal/service/ad/dto"
	"adult-short-videos/internal/service/ad/repository"
)

type AdService struct {
	repo repository.AdRepository
}

func NewAdService(repo repository.AdRepository) *AdService {
	return &AdService{repo: repo}
}

func (s *AdService) GetBanners(ctx context.Context) ([]*dto.BannerResp, error) {
	list, err := s.repo.ListBanners(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.BannerResp, 0, len(list))
	for _, b := range list {
		result = append(result, &dto.BannerResp{
			ID:       b.ID,
			Title:    b.Title,
			SubTitle: b.SubTitle,
			BgStyle:  b.BgStyle,
			Tag:      b.Tag,
			TagClass: b.TagClass,
			URL:      b.URL,
		})
	}
	return result, nil
}

func (s *AdService) GetApps(ctx context.Context, category string) ([]*dto.AppResp, error) {
	list, err := s.repo.ListApps(ctx, category)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.AppResp, 0, len(list))
	for _, a := range list {
		result = append(result, &dto.AppResp{
			ID:         a.ID,
			Name:       a.Name,
			Icon:       a.Icon,
			IconText:   a.IconText,
			BgStyle:    a.BgStyle,
			Official:   a.Official,
			Hot:        a.Hot,
			Categories: a.Categories,
			URL:        a.URL,
		})
	}
	return result, nil
}
