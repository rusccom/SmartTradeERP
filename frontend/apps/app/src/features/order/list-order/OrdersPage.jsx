import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";
import { createApiTablePreset } from "@smarterp/model/data-table/createApiTablePreset";
import { ServerListTable } from "@smarterp/ui/tables/list-table";
import { orderPath } from "../menu";
import { createOrderColumns } from "./orderColumns";
import "./list-order.css";

export default function OrdersPage() {
  const { t, locale } = useI18n();
  const navigate = useNavigate();
  const preset = useMemo(() => ordersPreset(t, locale), [t, locale]);
  return <section className="orders-list" aria-labelledby="orders-title">
    <h1 id="orders-title">{t("orders.list.title")}</h1>
    <ServerListTable preset={preset} selectable showHeaderWhenEmpty emptyText={t("orders.list.empty")}
      selectionLabels={{ all: t("orders.list.selectAll"), row: (r) => t("orders.list.selectRow", { number: r.number }) }}
      onRowOpen={(r) => navigate(`${orderPath}/${r.id}`)}
      primaryAction={{ label: t("entry.new"), onClick: () => navigate(`${orderPath}/new`) }} />
  </section>;
}

function ordersPreset(t, locale) {
  const columns = createOrderColumns(t, locale);
  return createApiTablePreset({
    id: "orders", path: apiPaths.orders, columns, rowId: (r) => r.id,
    capabilities: { sorting: true, search: true }, sorting: columns.map((c) => c.accessorKey),
    defaultState: { sorting: [{ id: "created_at", desc: true }] },
  });
}
