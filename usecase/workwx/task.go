package workwx

//type param struct {
//	unfinishedCount int // 待完成数总和
//	failCount       int // 错误数总和
//	timeout         *time.Timer
//}
//type Task struct {
//	sync.Mutex
//	failTH             int // 失败次数阈值
//	closeFlag          bool
//	wg                 *sync.WaitGroup
//	params             map[int64]*param
//	cbSuccess          map[int64]func() // 成功回调
//	cbFail             map[int64]func() // 失败回调
//	signSuccess        chan int64       // 成功信号量
//	signFail           chan int64       // 失败信号量
//	signRetry          chan func() error
//	signclose, cbclose chan struct{}
//	logger             *log.Logger
//}
//
//func NewTask(failTH int, logger *log.Logger) {
//	task := Task{
//		Mutex:       sync.Mutex{},
//		failTH:      failTH,
//		wg:          &sync.WaitGroup{},
//		params:      make(map[int64]*param),
//		cbSuccess:   make(map[int64]func()),
//		cbFail:      make(map[int64]func()),
//		signSuccess: make(chan int64, 10),
//		signFail:    make(chan int64, 10),
//		signRetry:   make(chan func(), 100),
//		signclose:   make(chan struct{}, 1),
//		cbclose:     make(chan struct{}),
//		logger:      logger,
//	}
//	task.wg.Add(1)
//}
//
//func (t *Task) Retry(g *sync.WaitGroup) {
//	defer g.Done()
//	for {
//		t.logger.Info("开始服务")
//		select {
//		case <-t.cbclose:
//			t.logger.Info("收到回调通道关闭信号")
//			return
//		case f, ok := <-t.signRetry:
//			if !ok {
//				t.logger.Info("Retry通道已关闭")
//			}
//			err := f()
//			if err != nil {
//
//			}
//		}
//	}
//}
//
//func (t *Task) notifyFail(id int64) (retry bool) {
//
//}
