import { Boxes, Package, Store, Truck } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";
import IntegrationTile from "../ui/IntegrationTile";
import "../integrations.css";

const providers = [
  { name: "InPost", icon: Package, tone: "inpost", path: "inpost" },
  { name: "Furgonetka", icon: Truck, tone: "furgonetka" },
  { name: "Base Courier", icon: Boxes, tone: "base" },
  { name: "Allegro", icon: Store, tone: "allegro" },
];

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
