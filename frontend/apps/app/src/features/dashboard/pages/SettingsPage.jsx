import { useI18n } from "../../../shared/i18n/useI18n";
import CurrencySettingsPanel from "../../currencies/ui/CurrencySettingsPanel";

function SettingsPage() {
  const { t } = useI18n();
  return (
    <section className="placeholder">
      <h2>{t("dashboard.settings.title")}</h2>
      <CurrencySettingsPanel />
    </section>
  );
}

export default SettingsPage;
