export function createOrderColumns(t, locale) {
  const date = new Intl.DateTimeFormat(locale, { dateStyle: "short", timeStyle: "short" });
  const amount = new Intl.NumberFormat(locale, { maximumFractionDigits: 4 });
  return [
    { accessorKey: "number", header: t("orders.columns.number"), openOnClick: true },
    { accessorKey: "customer_name", header: t("orders.columns.recipient"), cell: (name, row) => recipientCell(name, row.source, t) },
    { accessorKey: "total_with_delivery", header: t("orders.columns.total"), cell: (value, order) => amountCell(value, order, amount, t) },
    { accessorKey: "status", header: t("orders.columns.status") },
    { accessorKey: "created_at", header: t("orders.columns.received"), cell: (value) => receivedCell(value, date) },
  ];
}

function recipientCell(name, source, t) {
  const label = ["manual", "storefront"].includes(source) ? t(`entry.source.${source}`) : source;
  return <span className="orders-cell"><strong>{name}</strong><small>{label || t("orders.list.unknownSource")}</small></span>;
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
