import { useI18n } from "@smarterp/i18n/useI18n";
import PluginsPanel from "@smarterp/ui/plugins/PluginsPanel";
import CurrencySettingsPanel from "../currencies/ui/CurrencySettingsPanel";

function SettingsPage() {
  const { t } = useI18n();
  return (
    <section className="placeholder">
      <h2>{t("dashboard.settings.title")}</h2>
      <CurrencySettingsPanel />
      <PluginsPanel />
    </section>
  );
}

export default SettingsPage;
