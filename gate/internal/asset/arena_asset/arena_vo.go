package arena_asset

import "github.com/kim118000/core/pkg/config"

var _ config.IVoList = (*ArenaVoList)(nil)

type ArenaVo struct {
	RankScore int `json:"rank_score"`
	SeasonID  int `json:"season_id"`
}

type ArenaVoList struct {
	config.TemplateBase
	VoList  []*ArenaVo `json:"data"`
	hashmap map[int]*ArenaVo
}

func (avl *ArenaVoList) Init() {
	if avl.VoList != nil && len(avl.VoList) > 0 {
		if avl.hashmap == nil {
			avl.hashmap = make(map[int]*ArenaVo)
		}
		for _, vo := range avl.VoList {
			avl.hashmap[vo.SeasonID] = vo
		}
	}
}

//自己单独在其他文件中实现
func (avl *ArenaVoList) Update() bool {
	return true
}

func (avl *ArenaVoList) Check() bool {
	return true
}

func (avl *ArenaVoList) GetVo(key int) *ArenaVo {
	v, ok := avl.hashmap[key]
	if ok {
		return v
	}
	return nil
}
