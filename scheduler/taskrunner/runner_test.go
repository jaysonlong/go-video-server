package taskrunner

import (
	"errors"
	"log"
	"testing"
)

func TestRunner(t *testing.T) {
	producer := func(data dataChan) error {
		for i := 0; i < 5; i++ {
			data <- i
			log.Printf("produce: %v", i)
		}
		return nil
	}
	consumer := func(data dataChan) error {
	forloop:
		for {
			select {
			case v := <-data:
				log.Printf("consume: %v", v)
			default:
				break forloop
			}
		}
		// 若想停止本轮任务，只需返回error即可
		return errors.New("unknown error")
	}

	runner := NewTaskRunner(5, false, producer, consumer)
	runner.start()
}
