package themes

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

//go:embed */theme.json
var descriptors embed.FS

type Info struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PreviewURL string `json:"preview_url"`
	Default    bool   `json:"default,omitempty"`
}

type Index struct {
	Items   []Info
	Default string
}

// Load reads only descriptors. The ERP administration never loads storefront templates.
func Load() (*Index, error) {
	entries, err := descriptors.ReadDir(".")
	if err != nil {
		return nil, err
	}
	idx := &Index{}
	for _, entry := range entries {
		if err := idx.add(entry.Name()); err != nil {
			return nil, err
		}
	}
	if idx.Default == "" {
		return nil, errors.New("missing default storefront theme")
	}
	return idx, nil
}

func (idx *Index) add(id string) error {
	raw, err := descriptors.ReadFile(id + "/theme.json")
	if err != nil {
		return err
	}
	var info Info
	if err := json.Unmarshal(raw, &info); err != nil {
		return err
	}
	if info.ID != id || info.Name == "" {
		return fmt.Errorf("invalid theme descriptor %q", id)
	}
	if info.Default && idx.Default != "" {
		return errors.New("multiple default storefront themes")
	}
	if info.Default {
		idx.Default = id
	}
	idx.Items = append(idx.Items, info)
	return nil
}

func (idx *Index) Has(id string) bool {
	return slices.ContainsFunc(idx.Items, func(i Info) bool { return i.ID == id })
}
