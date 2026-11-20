export const catalogSettingsPath = "/dashboard/catalog/settings";

export function catalogSettingsMenu(t) {
  return { key: "catalogSettings", path: catalogSettingsPath, title: t("dashboard.menu.settings") };
}
