package model

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

// 问题类型_ 0-随机,1-简单常识,2-简单计算,3-逻辑推理,4-科学常识,5-常见诗词成语,6-简单历史,7-娱乐题型 8-脑筋急转弯……大于1000-企业定制
// INSERT INTO audited_items (item_id,question_type,question,answer) VALUES (1,5,"床前？？光","明月")

func TestMysql(t *testing.T) {
	d, _ := gorm.Open("mysql", "root:123456@tcp(43.143.208.232:3306)/cmmvplat?charset=utf8&parseTime=true&loc=Local")
	defer d.Close()
	var g DefaultAuditedItemsModel
	g.db = d
	g.rp = &redis.Pool{
		MaxIdle:     2,
		MaxActive:   0,
		IdleTimeout: 5 * 60 * 1000 * 1000,
		Dial: func() (redis.Conn, error) {
			setdb := redis.DialDatabase(0)      // 设置db
			setPasswd := redis.DialPassword("") // 设置redis连接密码
			c, err := redis.Dial("tcp", "43.143.208.232:6379", setdb, setPasswd)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	var ai = &AuditedItem{}
	if err := d.Select("item_id").Order("item_id DESC").Limit(1).Find(ai).Error; err != nil {
		t.Log("err///", err)
	}
	t.Log("result:", ai.ItemID)
}

func TestRedis(t *testing.T) {
	redisPool := &redis.Pool{
		MaxIdle:     2,
		MaxActive:   0,
		IdleTimeout: 5 * 60 * 1000 * 1000,
		Dial: func() (redis.Conn, error) {
			setdb := redis.DialDatabase(0)      // 设置db
			setPasswd := redis.DialPassword("") // 设置redis连接密码
			c, err := redis.Dial("tcp", "43.143.208.232:6379", setdb, setPasswd)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	rc := redisPool.Get()
	defer rc.Close()
	key := "cache:output:audited_items:id:2"
	_, err := rc.Do("SETEX", key, 3, 548142)
	val, err := rc.Do("Get", key)
	if err != nil {
		t.Log("err///", err)
	}
	t.Log("result///", val.([]uint8))
}
