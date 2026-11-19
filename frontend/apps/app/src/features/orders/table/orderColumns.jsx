const receivedDate = new Intl.DateTimeFormat("uk-UA", { dateStyle: "short", timeStyle: "short" });
const amountFormat = new Intl.NumberFormat("uk-UA", { maximumFractionDigits: 4 });

export const orderColumns = [
  { accessorKey: "number", header: "№ замовлення" },
  { accessorKey: "customer_name", header: "Отримувач / джерело", cell: recipientCell },
  { accessorKey: "total_with_delivery", header: "Сума з доставкою", cell: amountCell },
  { accessorKey: "status", header: "Статус" },
  { accessorKey: "created_at", header: "Дата надходження", cell: receivedCell },
];

function recipientCell(name) {
  return <span className="orders-cell"><strong>{name}</strong><small>Джерело —</small></span>;
}

function amountCell(withDelivery, order) {
  const amount = withDelivery ?? order.total_amount;
  if (amount == null) return "—";
  return (
    <span className="orders-cell">
      <span>{amountFormat.format(amount)} {order.currency_code}</span>
      {withDelivery == null && <small>Без вартості доставки</small>}
    </span>
  );
}

function receivedCell(value) {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return <time dateTime={value}>{receivedDate.format(date)}</time>;
}
