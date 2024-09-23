package pkg

import (
	"github.com/sony/sonyflake"
	"time"
)

var generatorId *sonyflake.Sonyflake

func init() {
	var err error
	generatorId, err = sonyflake.New(sonyflake.Settings{
		StartTime: time.Now(),
	})
	if err != nil {
		panic(err)
	}
}

func GenerateId() (uint64, error) {
	return generatorId.NextID()
}
