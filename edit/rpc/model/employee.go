package model

import (
	"log"
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
	err := d.db.Select("employee_level, contribution_score, audit_score, registration_time").Where("mobile_num = ?", mobile).Limit(1).Find(&e).Error
	if err != nil {
		return nil, err
	}
	log.Printf("查询员工表成功：%v", e)

	return &e, nil
}
