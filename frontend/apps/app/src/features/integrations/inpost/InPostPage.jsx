import { ArrowLeft } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import InPostSettings from "./InPostSettings";
import "../integrations.css";
import "./settings.css";

export default function InPostPage() {
  const { t } = useI18n();
  return (
    <section className="integrations inpost-settings" aria-labelledby="inpost-title">
      <Link className="integration-back" to="..">
        <ArrowLeft size={16} aria-hidden="true" />{t("dashboard.menu.integrations")}
      </Link>
      <h1 id="inpost-title">{t("inpost.title")}</h1>
      <p className="integration-intro">{t("inpost.intro")}</p>
      <InPostSettings />
    </section>
  );
}
