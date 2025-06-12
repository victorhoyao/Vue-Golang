package WorkPool

import "sync"

type WorkerPool struct {
	wg      sync.WaitGroup
	jobChan chan func()
}

// 初始化协程池
func NewWorkerPool(maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobChan: make(chan func(), maxWorkers*2), // 缓冲大小
	}
	// 启动固定数量的 worker
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for job := range pool.jobChan {
				job()
				pool.wg.Done()
			}
		}()
	}
	return pool
}

// 提交任务到协程池
func (p *WorkerPool) Submit(job func()) {
	p.wg.Add(1)
	p.jobChan <- job
}

// 等待所有任务完成
func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.jobChan)
}
