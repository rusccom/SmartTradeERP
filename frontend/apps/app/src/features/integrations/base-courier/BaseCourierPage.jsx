import { ArrowLeft } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import BaseCourierAccount from "./BaseCourierAccount";
import BaseCourierSettings from "./BaseCourierSettings";
import "../integrations.css";
import "./settings.css";

export default function BaseCourierPage() {
  const { t } = useI18n();
  return (
    <section className="integrations base-courier-settings" aria-labelledby="base-courier-title">
      <Link className="integration-back" to="..">
        <ArrowLeft size={16} aria-hidden="true" />{t("dashboard.menu.integrations")}
      </Link>
      <h1 id="base-courier-title">{t("integrations.providers.baseCourier")}</h1>
      <BaseCourierSettings />
      <BaseCourierAccount />
    </section>
  );
}
