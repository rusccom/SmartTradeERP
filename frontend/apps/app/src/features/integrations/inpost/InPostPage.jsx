import { ArrowLeft, MapPin, PackagePlus, Printer, Truck } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import InPostToolLink from "./InPostToolLink";
import "../integrations.css";

const manager = "https://manager.paczkomaty.pl";
const tools = [
  { key: "create", icon: PackagePlus, url: `${manager}/shipments/send/simple` },
  { key: "manage", icon: Printer, url: `${manager}/shipments/list` },
  { key: "pickups", icon: Truck, url: `${manager}/shipments/list/dispatch-orders` },
  { key: "track", icon: MapPin, url: "https://inpost.pl/sledzenie-przesylek" },
];

export default function InPostPage() {
  const { t } = useI18n();
  return (
    <section className="integrations" aria-labelledby="inpost-title">
      <Link className="integration-back" to="..">
        <ArrowLeft size={16} aria-hidden="true" />{t("dashboard.menu.integrations")}
      </Link>
      <h1 id="inpost-title">InPost</h1>
      <p className="integration-intro">{t("integrations.inpost.intro")}</p>
      <ul className="inpost-tools">
        {tools.map((tool) => <li key={tool.key}><InPostToolLink tool={tool} /></li>)}
      </ul>
    </section>
  );
}
