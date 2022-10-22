package model

import (
	"encoding/json"
	"errors"
	"log"
	"pkg/conf"
	"rpc/internal/svc"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

var count uint = 0

const (
	KEY_PREFIX         = "cache:output:audited_items:id:"
	KEY_COUNT_PREFIX   = "cache:output:audited_items:count:type:"
	EXPIRATION_SECONDS = 3 * 60
)

type (
	AuditedItemsModel interface {
		// FindRandomOne() (*AuditedItem, error)
		FindRandomOneByType(questionType uint32) (*AuditedItem, error)
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

func (d *DefaultAuditedItemsModel) FindRandomOneByType(questionType uint32) (*AuditedItem, error) {

	c := d.rp.Get()
	defer c.Close()

	//查询数据数量
	total := 0
	t, err := redis.String(c.Do("GET", KEY_COUNT_PREFIX+strconv.Itoa(int(questionType))))
	if err != nil {
		//未缓存则读取数据库
		if questionType == 0 {
			err = d.db.Table("audited_items").Count(&total).Error
		} else {
			err = d.db.Table("audited_items").Where("question_type = ?", questionType).Count(&total).Error
		}
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("SETEX", KEY_COUNT_PREFIX+strconv.Itoa(int(questionType)), EXPIRATION_SECONDS, strconv.Itoa(total)); err != nil {
			log.Fatal(err)
		}
	} else {
		total, err = strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
	}
	if total == 0 {
		return nil, errors.New(conf.GlobalError[conf.NO_DATA_WAS_QUERIED])
	}

	//查询条目
	var item AuditedItem
	s, err := redis.String(c.Do("GET", KEY_PREFIX+strconv.Itoa(int(count%uint(total)))))
	if err != nil {
		if questionType == 0 {
			err = d.db.Select("question,answer,disturb_answer").Offset(count % uint(total)).Limit(1).Find(&item).Error
		} else {
			err = d.db.Select("question,answer,disturb_answer").Where("question_type = ?", questionType).Offset(count % uint(total)).Limit(1).Find(&item).Error
		}
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("SETEX", KEY_PREFIX+strconv.Itoa(int(count%uint(total))), EXPIRATION_SECONDS, string(b)); err != nil {
			log.Fatal(err)
		}
	} else {
		err := json.Unmarshal([]byte(s), &item)
		if err != nil {
			return nil, err
		}
	}

	count++
	return &item, nil
}

// func (d *DefaultAuditedItemsModel) FindRandomOne() (*AuditedItem, error) {

// 	c := d.rp.Get()
// 	defer c.Close()

// 	//读取缓存
// 	total := 0
// 	t, err := redis.String(c.Do("GET", KEY_COUNT_PREFIX))
// 	if err != nil {
// 		//未缓存则读取数据库
// 		err := d.db.Table("audited_items").Count(&total).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		if _, err := c.Do("SETEX", KEY_COUNT_PREFIX, EXPIRATION_SECONDS, strconv.Itoa(total)); err != nil {
// 			log.Fatal(err)
// 		}
// 	} else {
// 		total, err = strconv.Atoi(t)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if total == 0 {
// 		return nil, errors.New(conf.GlobalError[conf.NO_DATA_WAS_QUERIED])
// 	}

// 	var item AuditedItem
// 	s, err := redis.String(c.Do("GET", KEY_PREFIX+strconv.Itoa(int(count%uint(total)))))
// 	if err != nil {
// 		err = d.db.Select("question,answer").Offset(count % uint(total)).Limit(1).Find(&item).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		b, err := json.Marshal(item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if _, err := c.Do("SETEX", KEY_PREFIX+strconv.Itoa(int(count%uint(total))), string(b)); err != nil {
// 			log.Fatal(err)
// 		}
// 	} else {
// 		err := json.Unmarshal([]byte(s), &item)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	count++
// 	return &item, nil
// }
