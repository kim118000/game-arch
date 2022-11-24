package persistence

import "sync/atomic"

type DataObject struct {
	isInDb     bool
	lifecycle  ILifeCycle
	version    int32
	needUpdate int32
}

func (d *DataObject) IsInDb() bool {
	return d.isInDb
}

func (d *DataObject) SetInDb(flag bool) {
	d.isInDb = flag
}

func (d *DataObject) GetUpdateVersion() int32 {
	return atomic.LoadInt32(&d.version)
}

func (d *DataObject) addVersion() {
	atomic.AddInt32(&d.version, 1)
}

func (d *DataObject) SetUpdateState(version int32) {
	if d.GetUpdateVersion() == version {
		return
	}
	d.setUpdate()
}

func (d *DataObject) IsNeedUpdate() bool {
	n := atomic.LoadInt32(&d.needUpdate)
	return n == 1
}

func (d *DataObject) setUpdate() {
	atomic.StoreInt32(&d.needUpdate, 1)
}

func (d *DataObject) GetLifeCycle() ILifeCycle {
	return d.lifecycle
}

func (d *DataObject) DoUpdated() {

}

func (d *DataObject) SetDestroy() {
	d.lifecycle.Destroy()
	d.addVersion()
	d.setUpdate()
	d.DoUpdated()
}

func (d *DataObject) SetModified() {
	if !d.lifecycle.CheckModifiable() {
		return
	}
	d.addVersion()
	d.setUpdate()
	d.DoUpdated()
}

func (d *DataObject) CheckNeedUpdate() bool {
	if d.IsNeedUpdate() {
		d.DoUpdated()
	}
	return d.IsNeedUpdate()
}