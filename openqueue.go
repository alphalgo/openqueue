package openqueue

//
// OPENQUEUE JUST AN EXPANDED QUEUE DATA STRUCTURE SIMULATION WITHOUT FULLY TEST,
// IN SOME WAY OPENSTACK IS A SUPPOSE WITH MY IDEAS AND UNVALIDATE IN REAL WORK.
//

import (
	log "github.com/sirupsen/logrus"
)

// Elem defines the element of oqueue
type Elem struct {
	position  int
	allocated bool
}

// Oqueue defines the queue object
type Oqueue struct {
	entranced bool
	exited    bool
	fenced    []bool
	size      int
	begin     int
	bottom    []*Elem
	top       []*Elem
	empty     bool
	mapped    bool
	_map      map[int]*Elem
	__map     map[*Elem][]interface{}
}

// Openqueue defines the oqueue interface
type Openqueue interface {
	init()
	Size() int
	List() []*Elem
	GetBottom(index int) *Elem
	GetTop(index int) *Elem
	AddElem(e *Elem) bool
	RemoveElem(e *Elem) bool
	IsEmpty() bool
	Check(err error)
	IsExist(e *Elem) bool
	SetMap(e *Elem, v ...interface{}) bool
	GetMap(index int) (*Elem, []interface{})
	Delete(index int)
	Destroy(index int)
}

func (o *Oqueue) init() {
	o.entranced = true
	o.exited = true
	for i := 0; i < o.size; i++ {
		o.fenced[i] = false
		o.bottom[i] = nil
		o.top[i] = nil
	}
	o.size = 0
	o.begin = 0
	o.empty = true
	o.mapped = false
}

// Size gets the size of oqueue
func (o *Oqueue) Size() int {
	if o.empty {
		return 0
	}
	return o.size
}

// List shows the element of oqueue
func (o *Oqueue) List() []*Elem {
	store := make([]*Elem, o.size)
	if !o.IsEmpty() {
		for key, value := range o._map {
			store[key] = value
		}
	}
	return store
}

// GetBottom gets the bottom of oqueue
func (o *Oqueue) GetBottom(index int) *Elem {
	if index < 0 {
		o.exited = false
		return o.bottom[o.size+index]
	} else if index == 0 {
		log.Warnf("Index cannot be zero!")
	}
	o.entranced = false
	return o.bottom[index]
}

// GetTop gets the top of oqueue
func (o *Oqueue) GetTop(index int) *Elem {
	if index < 0 {
		o.exited = false
		return o.top[o.size+index]
	} else if index == 0 {
		log.Warnf("Index cannot be zero!")
	}
	o.entranced = false
	return o.top[index]
}

// AddElem adds the element to oqueue
func (o *Oqueue) AddElem(e *Elem) (added bool) {
	if e.position < 0 || e.position > o.size {
		added = false
		log.Warnf("Cannot add element cause of invalid index.")
	}

	if o.empty {
		o._map[0] = e
	}

	if !e.allocated {
		o._map[o.begin] = e
		o.begin++
	}
	added = true
	return
}

// RemoveElem removes the element from oqueue
func (o *Oqueue) RemoveElem(e *Elem) (removed bool) {
	if o.empty {
		removed = false
	}

	if o._map[e.position] != nil {
		o.Destroy(e.position)
	}
	removed = true
	return
}

// isEmpty checks if the openqueue is empty.
func (o *Oqueue) IsEmpty() bool {
	if o.empty {
		return true
		log.Warnf("Openqueue is empty, please add some elements.")
	}
	return false
}

// Check checks if err
func (o *Oqueue) Check(err error) {
	if err != nil {
		panic(err)
	}
}

// isExist checks if the element e exist.
func (o *Oqueue) IsExist(e *Elem) bool {
	if o.__map[e] != nil && e != nil {
		return true
	} else if o.__map[e] == nil && e != nil {
		return true
	}
	return false
}

// When openqueue's fence status is true, we turn down the entrance
// and start to mapping elements.
func (o *Oqueue) SetMap(e *Elem, v ...interface{}) (mapped []bool) {
	for k := 0; k < o.size; k++ {
		if !o.fenced[k] {
			mapped[k] = false
		}
		o.entranced = true
		mapped[k] = true
	}
	// Here means element e mapped to any value but with slice, for
	// now, this map indicate the related resource of element e.
	// But I don't know whether the related resource means some
	// requests or connections or other stuffs which we can processed.
	// So, in here, we use null interface slice, refactor later.
	if len(v) > 0 {
		o.__map[e] = v[0].([]interface{})
	}
	o.__map[e] = v

	return
}

// getMap gets the value mapped by index and return it's mapping value.
func (o *Oqueue) GetMap(index int) (e *Elem, v []interface{}) {
	if index < 0 || index > o.size {
		log.Warnf("Invalid index.")
	}
	e = o._map[index]
	v = o.__map[e]

	return
}

// In common queue, we just delete that element then move forward.
// But in openqueue, we just delete that element and remove its 'grid'
// and add into tail. But we have follows cases:
//		1) the element's memory cannot occupied before the occupied one
//		2) the map full of elements, we just delete the index mapped to
//		3) we can delete element from begin or end point
func (o *Oqueue) Destroy(index int) {
	if !o._map[index].allocated {
		log.Warnf("This position have no element, you needn't to delete it!")
	}
	if !o.IsExist(o._map[index-1]) && o.IsExist(o._map[index]) {
		log.Warnf("Invalid request, impossible!")
	}
	delete(o._map, index)
	delete(o.__map, o._map[index])
}

// Delete deletes the element by given index, if index < 0, we delete
// element from backend, else we delete element from front.
func (o *Oqueue) Delete(index int) {
	if index < 0 {
		o.entranced = false
		o.Destroy(index)
	}
	o.exited = false
	o.Destroy(index)
}
