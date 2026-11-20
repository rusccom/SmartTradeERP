import { useI18n } from "@smarterp/i18n/useI18n";
import { useBundleSettings } from "./useBundleSettings";
import BundlePricingFields from "./BundlePricingFields";
import BundleSlotRule from "./BundleSlotRule";
import "./smartbundle.css";

export default function SmartBundleSettings({ bundle }) {
  const { t } = useI18n();
  const { editor, feedback, busy, patch, save } = useBundleSettings(bundle);
  const settings = editor?.settings;
  return (
    <details className="smartbundle-settings"><summary>{t("smartBundle.title")}</summary>
      {settings && <form onSubmit={event => { event.preventDefault(); save(); }}>
        <BundlePricingFields settings={settings} patch={patch} />
        {editor.components.map(slot => <BundleSlotRule key={slot.key} slot={slot} rule={settings.slots[slot.key]} mode={settings.mode}
          onChange={rule => patch({ slots: { ...settings.slots, [slot.key]: rule } })} />)}
        <button disabled={busy}>{t("smartBundle.save")}</button>
      </form>}
      <p role="status">{feedback?.saved ? t("smartBundle.saved") : feedback?.error}</p>
    </details>
  );
}
