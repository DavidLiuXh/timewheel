package timewheel

import (
	"container/list"
	"container/ring"
	"time"
)

type TimeWheelCallback func([]interface{})

type TimeWheel struct {
	interval       time.Duration
	ticker         *time.Ticker
	callback       TimeWheelCallback
	buckets        *ring.Ring
	currentPos     int
	requestChannel chan interface{}
	quitChannel    chan interface{}
}

type BucketItem struct {
	interval time.Duration
	value    interface{}
}

func NewTimeWhell(interval time.Duration, bucketCount int, callbackFunc TimeWheelCallback) *TimeWheel {
	timeWheel := &TimeWheel{
		buckets:        ring.New(bucketCount),
		interval:       interval,
		callback:       callbackFunc,
		currentPos:     0,
		requestChannel: make(chan interface{}),
		quitChannel:    make(chan interface{}),
	}

	for i := 1; i <= timeWheel.buckets.Len(); i++ {
		timeWheel.buckets.Value = list.New()
		timeWheel.buckets = timeWheel.buckets.Next()
	}

	return timeWheel
}

func (timeWheel *TimeWheel) Start() {
	timeWheel.ticker = time.NewTicker(timeWheel.interval * time.Second)
	go timeWheel.run()
}

func (timeWheel *TimeWheel) Add(interval time.Duration, item ...interface{}) {
	timeWheel.requestChannel <- &BucketItem{
		interval: interval,
		value:    item,
	}
}

func (timeWheel *TimeWheel) Stop() {
	close(timeWheel.quitChannel)
}

func (timeWheel *TimeWheel) run() {
	for {
		select {
		case <-timeWheel.quitChannel:
			timeWheel.ticker.Stop()
			break
		case <-timeWheel.ticker.C:
			if nil != timeWheel.callback {
				buckets := timeWheel.buckets.Value.(*list.List)
				userDatas := make([]interface{}, 0, buckets.Len())

				var n *list.Element
				for v := buckets.Front(); v != nil; v = n {
					if bucketItem := v.Value.(*BucketItem); nil != bucketItem {
						insertBucket := timeWheel.buckets.Move(int(bucketItem.interval.Nanoseconds()) / int(timeWheel.interval))
						itemList, _ := insertBucket.Value.(*list.List)
						itemList.PushBack(bucketItem)

						userDatas = append(userDatas, bucketItem.value)

						n = v.Next()
						buckets.Remove(v)
					}
				}

				if len(userDatas) > 0 {
					timeWheel.callback(userDatas)
				}
			}

			timeWheel.buckets = timeWheel.buckets.Next()

		case item := <-timeWheel.requestChannel:
			if bucketItem, _ := item.(*BucketItem); nil != bucketItem {
				insertBucket := timeWheel.buckets.Move(int(bucketItem.interval.Nanoseconds()) / int(timeWheel.interval))
				if itemList, _ := insertBucket.Value.(*list.List); nil != itemList {
					itemList.PushBack(bucketItem)
				}
			}
		}
	}
}
