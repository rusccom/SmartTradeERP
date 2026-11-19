import { useI18n } from "@smarterp/i18n/useI18n";
import inpost from "../inpost/provider";
import furgonetka from "../furgonetka/provider";
import baseCourier from "../base-courier/provider";
import allegro from "../allegro/provider";
import IntegrationTile from "../ui/IntegrationTile";
import "../integrations.css";

const providers = [inpost, furgonetka, baseCourier, allegro];

export default function IntegrationsPage() {
  const { t } = useI18n();
  return (
    <section className="integrations" aria-labelledby="integrations-title">
      <h1 id="integrations-title">{t("dashboard.menu.integrations")}</h1>
      <ul className="integration-grid">
        {providers.map((provider) => <li key={provider.name}><IntegrationTile provider={provider} /></li>)}
      </ul>
    </section>
  );
}
