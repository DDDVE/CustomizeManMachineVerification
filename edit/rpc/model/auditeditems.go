package model

import (
	"errors"
	"pkg/conf"
	"rpc/internal/svc"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type (
	AuditedItemsModel interface {
		Insert(item *AuditedItem) error
	}
	AuditedItem struct {
		ID            uint64
		ItemID        uint64
		Producer      string
		QuestionType  uint8
		Question      string
		Answer        string
		DisturbAnswer string
		CreateTime    time.Time
	}
	DefaultAuditedItemsModel struct {
		db *gorm.DB
		rp *redis.Pool
	}
)

func NewDefaultAuditedItemsModel(d *svc.ServiceContext) AuditedItemsModel {
	return &DefaultAuditedItemsModel{db: d.GormDB, rp: d.RedisPool}
}

func (d *DefaultAuditedItemsModel) Insert(item *AuditedItem) error {
	if item.ItemID == 0 {
		var ai = &AuditedItem{}
		if err := d.db.Select("item_id").Order("item_id DESC").Limit(1).Find(ai).Error; err != nil {
			return err
		}
		item.ItemID = ai.ItemID + 1
	}
	if !d.db.NewRecord(*item) {
		return errors.New(conf.GlobalError[conf.INSERT_DATA_ERROR])
	}

	var e = &Employee{}
	if err := d.db.Select("contribution_score").Where("mobile_num = ?", item.Producer).Limit(1).Find(e).Error; err != nil {
		return err
	}
	if err := d.db.Model(&Employee{}).Where("mobile_num = ?", item.Producer).Update("contribution_score", e.ContributionScore+10).Error; err != nil {
		return err
	}

	return nil
}
