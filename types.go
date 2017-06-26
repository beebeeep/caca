package main

type CacusInstance struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
	Default bool   `yaml:"default"`
}

type CacaConfig struct {
	Instances map[string]CacusInstance `yaml:'instances'`
}

type CacusStatus struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
}

type DistroInfo struct {
	CacusStatus
	Result []struct {
		Distro     string   `json:"distro"`
		Components []string `json:"components"`
	} `json:"result"`
}
