export function createOrderColumns(t, locale) {
  const date = new Intl.DateTimeFormat(locale, { dateStyle: "short", timeStyle: "short" });
  const amount = new Intl.NumberFormat(locale, { maximumFractionDigits: 4 });
  return [
    { accessorKey: "number", header: t("orders.columns.number") },
    { accessorKey: "customer_name", header: t("orders.columns.recipient"), cell: (name) => recipientCell(name, t) },
    { accessorKey: "total_with_delivery", header: t("orders.columns.total"), cell: (value, order) => amountCell(value, order, amount, t) },
    { accessorKey: "status", header: t("orders.columns.status") },
    { accessorKey: "created_at", header: t("orders.columns.received"), cell: (value) => receivedCell(value, date) },
  ];
}

function recipientCell(name, t) {
  return <span className="orders-cell"><strong>{name}</strong><small>{t("orders.list.unknownSource")}</small></span>;
}

function amountCell(withDelivery, order, amountFormat, t) {
  const amount = withDelivery ?? order.total_amount;
  if (amount == null) return "—";
  return (
    <span className="orders-cell">
      <span>{amountFormat.format(amount)} {order.currency_code}</span>
      {withDelivery == null && <small>{t("orders.list.excludesDelivery")}</small>}
    </span>
  );
}

function receivedCell(value, receivedDate) {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return <time dateTime={value}>{receivedDate.format(date)}</time>;
}
