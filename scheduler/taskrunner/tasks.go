package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/midmis/go-video-server/scheduler/dbops"
)

var waitGroup = sync.WaitGroup{}

func VideoClearDispatcher(data dataChan) error {
	ids, err := dbops.ReadVideoDeletionRecord(BATCH_DELETE_CNT)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		// 所有任务完成，终止本轮任务的执行
		// 由于是longLived模式，runner不会关闭，而会在下一个ticker触发时重新启动
		return errors.New("All tasks finished")
	}

	for _, id := range ids {
		data <- id
	}
	return nil
}

func VideoClearExecutor(data dataChan) error {
	waitGroup.Add(len(data))

	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case id := <-data:
			// 使用协程将会直接返回，导致多个协程同时删除同一文件
			go func(videoId interface{}) {
				if err := deleteVideo(videoId.(string)); err != nil {
					errMap.Store(videoId, err)
					waitGroup.Done()
					return
				}

				if err := dbops.DeleteVideoDeletionRecord(videoId.(string)); err != nil {
					errMap.Store(videoId, err)
				}
				waitGroup.Done()
			}(id)

			// if err := deleteVideo(id.(string)); err != nil {
			// 	errMap.Store(id, err)
			// } else {
			// 	if err := dbops.DeleteVideoDeletionRecord(id.(string)); err != nil {
			// 		errMap.Store(id, err)
			// 	}
			// }
		default:
			break forloop
		}
	}

	// 等待所有协程退出
	waitGroup.Wait()

	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}

func deleteVideo(videoId string) error {
	err := os.Remove(VIDEO_PATH + videoId)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	log.Printf("Delete video success: %v", videoId)
	return nil
}
