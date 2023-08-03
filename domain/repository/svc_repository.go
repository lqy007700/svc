package repository

import (
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/model"
	"gorm.io/gorm"
)

type ISvcRepository interface {
	InitTable() error
	FindSvcById(id int64) (*model.Svc, error)
	CreateSvc(svc *model.Svc) (int64, error)
	DeleteSvcById(id int64) error
	Update(svc *model.Svc) error
	FindAll() ([]model.Svc, error)
}

func NewSvcRepository(db *gorm.DB) ISvcRepository {
	return &SvcRepository{
		db: db,
	}
}

type SvcRepository struct {
	db *gorm.DB
}

func (p *SvcRepository) InitTable() error {
	common.Info("Init table 11")
	return p.db.AutoMigrate(&model.Svc{}, &model.SvcPort{})
}

func (p *SvcRepository) FindSvcById(id int64) (*model.Svc, error) {
	svc := &model.Svc{}

	err := p.db.First(svc, id).Error
	return svc, err
}

func (p *SvcRepository) CreateSvc(svc *model.Svc) (int64, error) {
	err := p.db.Create(svc).Error
	return svc.ID, err
}

func (p *SvcRepository) DeleteSvcById(id int64) error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Where("id = ?", id).Delete(&model.Svc{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("svc_id = ?", id).Delete(&model.SvcPort{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *SvcRepository) Update(svc *model.Svc) error {
	return p.db.Model(svc).Updates(svc).Error
}

func (p *SvcRepository) FindAll() ([]model.Svc, error) {
	svcs := make([]model.Svc, 0)
	err := p.db.Find(&svcs).Error
	return svcs, err
}
