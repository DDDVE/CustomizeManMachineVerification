package model

import (
	"time"

	"rpc/internal/svc"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type (
	EmployeeModel interface {
		SelectByMobile(mobile string) (*Employee, error)
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

func (d *DefaultEmployeeModel) SelectByMobile(mobile string) (*Employee, error) {
	var e Employee
	err := d.db.Select("question,answer,disturb_answer").Where("mobile_num = ?", mobile).Limit(1).Find(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}
