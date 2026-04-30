# Table Filters

This folder contains generic table filter builders only.

Use these rules when adding or changing table filters:

- Keep generic filter infrastructure here.
- Keep feature-specific filter instances inside the feature folder that owns the table preset.
- A table preset connects a filter object through `search` or future filter config fields.
- The table filter object may define `id`, `enabled`, `queryKey`, `placeholderKey`, `debounceMs`, and `serialize`.
- Do not put product, document, bundle, customer, or backend SQL rules in this folder.
- Do not make the universal table aware of feature fields like SKU, barcode, variant name, document number, or customer email.
- Add new generic filter behavior through small factory functions, then import those factories from feature presets.
