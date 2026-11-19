import { useI18n } from "@smarterp/i18n/useI18n";

export default function SidebarHeader({ onClose }) {
  const { t } = useI18n();
  return (
    <div className="sidebar-header">
      <span className="sidebar-title">{t("dashboard.sidebar.menu")}</span>
      <button className="sidebar-close" type="button" onClick={onClose} aria-label={t("dashboard.sidebar.close")}>
        &times;
      </button>
    </div>
  );
}
