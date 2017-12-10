# timewheel
timewheel for golang.

It's not only an traditional timewheel (http://blog.csdn.net/mindfloating/article/details/8033340),
but also a supper set of timewheel.If you want ,you can put one or more serval kind of struct Data
which implement Entity into timewheel,each kind struct Data will fired OnExpired when expired  at
d duration. It's also means that you can store more than one kind of struct Data into one instance
of timewheel.

# usage



```
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
```
