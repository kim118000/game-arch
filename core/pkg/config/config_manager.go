package config

type ConfigManager struct {
	loader ILoader
	configs map[string]ITemplate
}

func NewConfigManager(loader ILoader) *ConfigManager {
	return &ConfigManager{
		loader: loader,
	}
}

func (cm *ConfigManager) RegTemplate(config ITemplate)  {
	if cm.configs == nil {
		cm.configs = make(map[string]ITemplate)
	}
	cm.configs[config.GetName()] = config
}

func (cm *ConfigManager) LoadTemplate() {
	for _, t := range cm.configs {
		t.Load(cm.loader)
	}
}

func (cm *ConfigManager) UpdateTemplate() {
	for _, t := range cm.configs {
		t.Update()
	}
}

func (cm *ConfigManager) CheckTemplate() {
	for _, t := range cm.configs {
		t.Check()
	}
}