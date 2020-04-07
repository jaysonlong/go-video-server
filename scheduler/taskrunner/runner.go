package taskrunner

import "log"

type TaskRunner struct {
	control    chan string //控制信息信道
	err        chan string //错误信息信道
	data       dataChan    //数据信道
	dataSize   int         //数据信道缓冲区大小
	dispatcher operator    //派发器，负责派发任务
	executor   operator    //执行器，负责执行任务
	longLived  bool        //是否常驻
}

func NewTaskRunner(dataSize int, longLived bool, dispatcher operator, executor operator) *TaskRunner {
	runner := &TaskRunner{}
	runner.control = make(chan string, 1)
	runner.err = make(chan string, 1)
	runner.dataSize = dataSize
	runner.data = make(dataChan, dataSize)
	runner.executor = executor
	runner.dispatcher = dispatcher
	runner.longLived = longLived
	return runner
}

func (tr *TaskRunner) startDispatch() {
	defer func() {
		// 执行到此处时必定返回了错误
		if !tr.longLived {
			close(tr.control)
			close(tr.err)
			close(tr.data)
		}
	}()

forloop:
	for {
		select {
		case msg := <-tr.control:
			switch msg {
			case READY_TO_DISPATCH:
				if err := tr.dispatcher(tr.data); err != nil {
					tr.err <- err.Error()
				} else {
					tr.control <- READY_TO_EXECUTE
				}
			case READY_TO_EXECUTE:
				if err := tr.executor(tr.data); err != nil {
					tr.err <- err.Error()
				} else {
					tr.control <- READY_TO_DISPATCH
				}
			}
		case errMsg := <-tr.err:
			log.Printf("%v", errMsg)
			break forloop

		default:
		}
	}
}

func (tr *TaskRunner) start() {
	tr.control <- READY_TO_DISPATCH
	tr.startDispatch()
}
