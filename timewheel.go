package timewheel

import (
	//"fmt"
	"sync"
	"time"
)

// Timewheel
type Timewheel struct {
	name     string        //name of Timewheel
	interval time.Duration //指针每隔多久向前移动一格（槽)
	ticker   *time.Ticker  //定时器
	slots    *sync.Map     //Slots时间轮槽 ,key:序号[0,1,2,3.....],value:Slot sync.Map
	stop     chan bool     // stop signal
	pos      int64         // current position of slots which was fired and the pos+1 will be fired at next tick
}

// 时间轮槽集合 ,key:序号[0,1,2,3.....],value:Slot
type Slots sync.Map

// 时间轮槽 ,key:序号[0,1,2,3.....],value:Entity
type Slot sync.Map

// Entity the data entity which stored into timewheel and be expired.
type Entity interface {
	OnExpired()       // fired when expired
	SetEId(eId int64) // must be unique ,and not equal zero
	GetEId() (eId int64)
	SetSlotId(slotId int64) // the slotId is id of slot which stored the entity
	GetSlotId() (slotId int64)
}

// NewTimewheel new an instance of Timewheel
func NewTimewheel(name string, interval time.Duration) *Timewheel {

	return &Timewheel{
		name:     name,
		interval: interval,
		slots:    &sync.Map{},
		stop:     make(chan bool, 1),
		pos:      0,
	}
}

// Start starts the Timewheel tw
func (tw *Timewheel) Start() {
	if tw.ticker == nil {
		tw.ticker = time.NewTicker(tw.interval)
		go tw.start()
	}
}

func (tw *Timewheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.pos++
			tw.tickHandler(tw.pos)
		case <-tw.stop:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *Timewheel) tickHandler(pos int64) {
	defer tw.slots.Delete(pos)
	//fmt.Printf("%s\t Timewheel[%s]\t expired at slot: %d \n", time.Now().Local(), tw.name, pos)
	slot, ok := tw.slots.Load(pos)
	if !ok || slot == nil {
		return
	}
	entitys, ok := slot.(*sync.Map)
	if !ok {
		return
	}
	go entitys.Range(func(key, value interface{}) bool {
		if entity, ok := value.(Entity); ok {
			entity.OnExpired()
		}
		return true
	})

}

// Stop stops the Timewheel tw
func (tw *Timewheel) Stop() {
	tw.stop <- true
}

// Add adds Entity e to Timewheel tw and return slotId
func (tw *Timewheel) Add(e Entity, delay time.Duration) (slotId int64) {
	if delay <= 0 {
		return -1
	}
	slotId = tw.pos + int64(delay.Seconds()/tw.interval.Seconds())
	slot, _ := tw.slots.LoadOrStore(slotId, &sync.Map{})

	slotMap, _ := slot.(*sync.Map)
	if eId := e.GetEId(); eId <= 0 {
		e.SetEId(time.Now().UnixNano())
	}

	slotMap.Store(e.GetEId(), e)
	e.SetSlotId(slotId)
	return slotId

}

// Remove remove Entity e from Timewheel tw
func (tw *Timewheel) Remove(e Entity) {
	if e.GetSlotId() < tw.pos {
		//tw.slots.Delete(e.GetSlotId())
		return
	}
	if slot, ok := tw.slots.Load(e.GetSlotId()); ok {
		if entitys, ok := slot.(*sync.Map); ok {
			entitys.Delete(e.GetEId())
		}
	}
}
