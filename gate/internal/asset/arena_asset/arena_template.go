package arena_asset

import (
	"encoding/json"
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/gate/internal/service"
)

var ArenaTemplate = new(arenaTemplate)

type arenaTemplate struct {
	VoList *ArenaVoList
}

func (at *arenaTemplate) GetName() string {
	return "arena_season"
}

func (at *arenaTemplate) GetFileName() string {
	return "arena_season.json"
}

func (at *arenaTemplate) Load(loader config.ILoader) {
	content, err := loader.Load(at.GetFileName())
	if err != nil {
		service.GateLog.Errorf("arena template load err %s", err)
		return
	}

	at.VoList = new(ArenaVoList)
	errJson := json.Unmarshal(content, at.VoList)
	if errJson != nil {
		service.GateLog.Errorf("arena template unmarshal json err %s", err)
	}

	at.VoList.Init()
}

func (at *arenaTemplate) Update() bool {
	return at.VoList.Update()
}

func (at *arenaTemplate) Check() bool {
	return at.VoList.Check()
}
