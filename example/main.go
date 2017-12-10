package main

import (
	"fmt"
	"github.com/anjieych/timewheel"
	"time"
)

func main() {
	tw := timewheel.NewTimewheel("tw-example", time.Second)
	tw.Start()
	tick := time.NewTicker(3 * time.Second)
	for {
		d := &Data{
			eid:  time.Now().UnixNano(),
			data: <-tick.C,
		}
		d.SetSlotId(tw.Add(d, 5*time.Second))
	}
}

// Data must implements timewheel.Entity
type Data struct {
	eid    int64
	slotId int
	data   interface{}
}

func (d *Data) SetEId(eId int64) {
	d.eid = eId
}
func (d *Data) GetEId() (eId int64) {
	return d.eid
}
func (d *Data) SetSlotId(slotId int) {
	d.slotId = slotId
}
func (d *Data) GetSlotId() (slotId int) {
	return d.slotId
}
func (d *Data) OnExpired() {
	fmt.Printf("%s\t OnExpired :{slotId: %d\t,eid: %d\t,data: %s}\n", time.Now(), d.GetSlotId(), d.GetEId(), d.data)
}
