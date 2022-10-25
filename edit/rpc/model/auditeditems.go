package model

import (
	"log"
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
	log.Printf("item.ItemID//////////////////////%+v", item)
	if item.ItemID == 0 {
		var ai = &AuditedItem{}
		if err := d.db.Select("item_id").Order("item_id DESC").Limit(1).Find(ai).Error; err != nil {
			return err
		}
		item.ItemID = ai.ItemID + 1
		log.Println("item.ItemID2//////////////////////", item.ItemID)
	}
	item.CreateTime = time.Now()
	d.db.Create(item)

	var e = &Employee{}
	if err := d.db.Select("contribution_score").Where("mobile_num = ?", item.Producer).Limit(1).Find(e).Error; err != nil {
		log.Println("1111111111111111111111111111111111111111111111111")
		return err
	}
	if err := d.db.Model(&Employee{}).Where("mobile_num = ?", item.Producer).Update("contribution_score", e.ContributionScore+10).Error; err != nil {
		log.Println("2222222222222222222222222222222222222222222222222")
		return err
	}

	log.Printf("插入条目表成功：%v", item)

	return nil
}
