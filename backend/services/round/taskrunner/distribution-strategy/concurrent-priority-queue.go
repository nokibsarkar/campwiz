package distributionstrategy

import (
	"container/heap"
	"log"
)

type heapPush struct {
	x any
}
type heapPop struct {
	result chan any
}
type NokibPriorityQueue[T heap.Interface] struct {
	h                 T
	actionQuitChannel chan bool
	heapPushChannel   chan heapPush
	heapPopChannel    chan heapPop
}

func (pq *NokibPriorityQueue[T]) Push(x any) {
	pq.heapPushChannel <- heapPush{x}
}
func (pq *NokibPriorityQueue[T]) Pop() any {
	var result = make(chan any)
	pq.heapPopChannel <- heapPop{result}
	return <-result
}
func (pq *NokibPriorityQueue[T]) StopWatchHeapOps() {
	pq.actionQuitChannel <- true
}
func (pq *NokibPriorityQueue[T]) WatchHeapOps() chan bool {
	go func() {
		for {
			select {
			case <-pq.actionQuitChannel:
				log.Println("Quitting")
				return
			case popMsg := <-pq.heapPopChannel:
				popMsg.result <- heap.Pop(pq.h)
			case pushMsg := <-pq.heapPushChannel:
				heap.Push(pq.h, pushMsg.x)
			}
		}
	}()
	return pq.actionQuitChannel
}
func NewNokibPriorityQueue[T heap.Interface](h T) *NokibPriorityQueue[T] {
	heap.Init(h)
	hp := &NokibPriorityQueue[T]{
		h:                 h,
		actionQuitChannel: make(chan bool),
		heapPushChannel:   make(chan heapPush),
		heapPopChannel:    make(chan heapPop),
	}
	hp.WatchHeapOps()
	return hp
}
