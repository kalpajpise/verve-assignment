package job

import (
	"fmt"
	"log"
	"time"

	"github.com/kalpaj/verve/pkg/aws"
	"github.com/kalpaj/verve/pkg/constant"
	"github.com/kalpaj/verve/pkg/db/redis"
)

type Job struct {
	redis   *redis.Redis
	kinesis *aws.Kinesis
	stream  string
}

func New(r *redis.Redis, k *aws.Kinesis, s string) *Job {
	return &Job{
		redis:   r,
		kinesis: k,
		stream:  s,
	}
}

// This job will send the uniques id to a distributed queue every clock minute.
// If app starts at 04:10:22, then the first event of "count" will be sent at
// 4:11:00 and there after every subsequent minute.
func (j *Job) Start() {
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	duration := nextMinute.Sub(now)

	// Sleep until the next minute and publish the first minute count
	time.Sleep(duration)
	log.Printf("[INFO] Publishing data to kinesis and resetting current minute count\n")
	j.doPublishAndReset()

	// Publish every minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("[INFO] Publishing data to kinesis and resetting current minute count\n")
		j.doPublishAndReset()
	}
}

func (j *Job) doPublishAndReset() {
	count, _ := j.redis.SetLength(constant.UniqueIDSet)

	// Send to kafka
	if err := j.kinesis.PublishMessage(j.stream, fmt.Sprint(count)); err != nil {
		log.Printf("[Error] Failed to publish data on kinesis: %s\n", err.Error())
	}

	// Reset the count
	j.redis.Delete(constant.UniqueIDSet)
}
