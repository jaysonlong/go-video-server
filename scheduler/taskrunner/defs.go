package taskrunner

const (
	READY_TO_DISPATCH     = "r2d"
	READY_TO_EXECUTE      = "r2e"
	BATCH_DELETE_CNT      = 3 // 每次批量删除的记录数
	BATCH_DELETE_INTERVAL = 3 // 批量删除任务的执行间隔

	VIDEO_PATH = "./videos/"
)

// 数据信道，用于消费者和生产者交换信息
type dataChan chan interface{}

// 数据操作函数
type operator func(dc dataChan) error
