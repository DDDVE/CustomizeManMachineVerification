package model

import (
	"errors"
	"time"

	"pkg/conf"
	"rpc/internal/svc"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type (
	EmployeeModel interface {
		Insert(e *Employee) error
	}
	Employee struct {
		ID                uint64
		MobileNum         string
		EmployeeLevel     uint8
		ContributionScore uint32
		AuditScore        uint32
		RegistrationTime  time.Time
	}
	DefaultEmployeeModel struct {
		db *gorm.DB
		rp *redis.Pool
	}
)

func NewDefaultEmployeeModel(sc *svc.ServiceContext) EmployeeModel {
	return &DefaultEmployeeModel{db: sc.GormDB, rp: sc.RedisPool}
}

func (d *DefaultEmployeeModel) Insert(e *Employee) error {
	count := 0
	d.db.Model(&Employee{}).Where("mobile_num = ?", e.MobileNum).Limit(1).Count(&count)
	if count == 0 {
		if !d.db.NewRecord(*e) {
			return errors.New(conf.GlobalError[conf.INSERT_DATA_ERROR])
		}
	}
	return nil
}
