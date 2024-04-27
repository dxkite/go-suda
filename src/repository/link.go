package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Link interface {
	Link(ctx context.Context, direct string, sourceId, linkedId uint64) error
	BatchLink(ctx context.Context, direct string, sourceId uint64, linkedIds []uint64) error
	LinkOnce(ctx context.Context, direct string, sourceId, linkedId uint64) error
	LinkOf(ctx context.Context, direct string, sourceId uint64) ([]*entity.Link, error)
	DeleteLink(ctx context.Context, direct string, sourceId, linkedId uint64) error
}

func NewLink(db *gorm.DB) Link {
	return &link{db: db}
}

type link struct {
	db *gorm.DB
}

func (r *link) Link(ctx context.Context, direct string, sourceId, linkedId uint64) error {
	link := entity.Link{}
	link.Direct = direct
	link.SourceId = sourceId
	link.LinkedId = linkedId
	return r.db.Create(&link).Error
}

func (r *link) BatchLink(ctx context.Context, direct string, sourceId uint64, linkedIds []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range linkedIds {
			if err := r.Link(ctx, direct, sourceId, v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *link) LinkOnce(ctx context.Context, direct string, sourceId, linkedId uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(entity.Link{Direct: direct, SourceId: sourceId}).Error; err != nil {
			return err
		}
		link := entity.Link{}
		link.Direct = direct
		link.SourceId = sourceId
		link.LinkedId = linkedId
		return r.db.Create(&link).Error
	})
}

func (r *link) LinkOf(ctx context.Context, direct string, sourceId uint64) ([]*entity.Link, error) {
	links := []*entity.Link{}
	if err := r.db.Model(entity.Link{}).Where(entity.Link{SourceId: sourceId}).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *link) DeleteLink(ctx context.Context, direct string, sourceId, linkedId uint64) error {
	if err := r.db.Delete(entity.Link{Direct: direct, SourceId: sourceId, LinkedId: linkedId}).Error; err != nil {
		return err
	}
	return nil
}
