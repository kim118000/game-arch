package persistence

//数据实体对象持久话
type IDbEntityUpdate interface {
	DoSave(entity any) (int64, error)
	DoUpdate(entity any) (int64, error)
	DoDelete(entity any) (int64, error)
	DoQuery(dest interface{}, query interface{}, args ...interface{}) error
}

type ILifeCycle interface {
	Active() bool
	IsActive() bool
	Destroy()
	IsDestroy() bool
	CheckModifiable() bool
	GetPersistenceObject() IPersistenceObject
}

type ICheckNeedUpdate interface {
	CheckNeedUpdate() bool
}

type IDataObject interface {
	DoUpdated()
	SetDestroy()
	IsInDb() bool
	SetInDb(flag bool)
	GetUpdateVersion() int32
	SetUpdateState(version int32)
	IsNeedUpdate() bool
	GetLifeCycle() ILifeCycle
	SetModified()
}

type IPersistenceObject interface {
	ICheckNeedUpdate
	IDataObject
	GetUniqueId() string
	GetId() uint64
	GetUpdateEntity() UpdateEntity
	GetEntity() any
	FromEntity(entity any)
	GetPoUpdater() IPersistenceUpdater
}

type IPersistenceUpdater interface {
	Save(ipo IPersistenceObject)
	Delete(ipo IPersistenceObject)
}
