import { useState } from "react";
import { DataTable } from "@smarterp/ui/tables/list-table";
import { orderColumns } from "../table/orderColumns";

const emptyOrders = [];
const selectionLabels = {
  all: "Вибрати всі замовлення на сторінці",
  row: (order) => `Вибрати замовлення ${order.number}`,
};

export default function OrdersPage({ orders }) {
  const [selection, setSelection] = useState({});
  const emptyText = orders ? "Замовлень поки немає." : "Список замовлень ще не підключено.";
  return (
    <section className="orders-list" aria-labelledby="orders-title">
      <h1 id="orders-title">Моє замовлення</h1>
      <DataTable
        columns={orderColumns} data={orders ?? emptyOrders} getRowId={(order) => order.id}
        selectable selectionLabels={selectionLabels} showHeaderWhenEmpty
        rowSelection={selection} onRowSelectionChange={setSelection}
        searchable={false} components={{ pagination: false }}
        emptyText={<span className="orders-empty">{emptyText}</span>}
      />
    </section>
  );
}
