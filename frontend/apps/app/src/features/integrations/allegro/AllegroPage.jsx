import { ArrowLeft } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import AllegroReturnNotice from "./oauth/AllegroReturnNotice";
import AllegroSettings from "./AllegroSettings";
import "../integrations.css";
import "./settings.css";

export default function AllegroPage() {
  const { t } = useI18n();
  return (
    <section className="integrations allegro-settings" aria-labelledby="allegro-title">
      <Link className="integration-back" to="..">
        <ArrowLeft size={16} aria-hidden="true" />{t("dashboard.menu.integrations")}
      </Link>
      <h1 id="allegro-title">{t("integrations.providers.allegro")}</h1>
      <AllegroReturnNotice />
      <AllegroSettings />
    </section>
  );
}
