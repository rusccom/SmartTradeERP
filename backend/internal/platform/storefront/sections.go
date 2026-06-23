package storefront

import "encoding/json"

// homeSectionOrder is the default order of home-page sections. The tenant can
// toggle and reorder them; unknown keys are ignored.
var homeSectionOrder = []string{"hero", "products"}

type sectionConfig struct {
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

// HomeSectionKeys is the catalog of toggleable home sections, exported for the
// admin so it validates against the same set the renderer knows.
func HomeSectionKeys() []string {
	return append([]string(nil), homeSectionOrder...)
}

// homeSections returns the ordered, enabled home-section keys from the stored
// config, falling back to the default order when nothing valid is configured.
func homeSections(raw []byte) []string {
	configs := parseSections(raw)
	if len(configs) == 0 {
		return HomeSectionKeys()
	}
	keys := enabledKeys(configs)
	if len(keys) == 0 {
		return keys
	}
	return keys
}

func parseSections(raw []byte) []sectionConfig {
	if len(raw) == 0 {
		return nil
	}
	configs := []sectionConfig{}
	if err := json.Unmarshal(raw, &configs); err != nil {
		return nil
	}
	return configs
}

func enabledKeys(configs []sectionConfig) []string {
	keys := make([]string, 0, len(configs))
	seen := map[string]bool{}
	for _, cfg := range configs {
		if cfg.Enabled && isKnownSection(cfg.Key) && !seen[cfg.Key] {
			keys = append(keys, cfg.Key)
			seen[cfg.Key] = true
		}
	}
	return keys
}

func isKnownSection(key string) bool {
	for _, known := range homeSectionOrder {
		if known == key {
			return true
		}
	}
	return false
}
