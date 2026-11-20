import { useMemo, useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { DataTable } from "@smarterp/ui/tables/list-table";
import { createOrderColumns } from "./orderColumns";
import "./list-order.css";

const emptyOrders = [];

export default function OrdersPage({ orders }) {
  const { t, locale } = useI18n();
  const columns = useMemo(() => createOrderColumns(t, locale), [t, locale]);
  const selectionLabels = useMemo(() => orderSelectionLabels(t), [t]);
  const [selection, setSelection] = useState({});
  const emptyText = t(orders ? "orders.list.empty" : "orders.list.disconnected");
  return (
    <section className="orders-list" aria-labelledby="orders-title">
      <h1 id="orders-title">{t("orders.list.title")}</h1>
      <DataTable
        columns={columns} data={orders ?? emptyOrders} getRowId={(order) => order.id}
        selectable selectionLabels={selectionLabels} showHeaderWhenEmpty
        rowSelection={selection} onRowSelectionChange={setSelection}
        searchable={false} components={{ pagination: false }}
        emptyText={<span className="orders-empty">{emptyText}</span>}
      />
    </section>
  );
}

function orderSelectionLabels(t) {
  return {
    all: t("orders.list.selectAll"),
    row: (order) => t("orders.list.selectRow", { number: order.number }),
  };
}
