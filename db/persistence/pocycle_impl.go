package persistence

type PoCycleImpl struct {
	state CycleState
	po    IPersistenceObject
}

func (pci *PoCycleImpl) Active() bool {
	if pci.state == Destroy {
		return false
	}
	return true
}

func (pci *PoCycleImpl) IsActive() bool {
	return pci.state == Active
}

func (pci *PoCycleImpl) Destroy() {
	pci.state = Destroy
}

func (pci *PoCycleImpl) IsDestroy() bool {
	return pci.state == Destroy
}

func (pci *PoCycleImpl) CheckModifiable() bool {
	return pci.state != Destroy
}
func (pci *PoCycleImpl) GetPersistenceObject() IPersistenceObject {
	return pci.po
}
