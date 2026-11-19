import { DataTable } from "@smarterp/ui/tables/list-table";

export default function StatusDraftTable({ statuses, onEdit }) {
  const columns = [
    { accessorKey: "name", header: "Назва статусу" },
    { id: "actions", header: "Дії", size: 116, cell: (_, status) => editButton(status, onEdit) },
  ];
  return (
    <DataTable data={statuses} columns={columns} getRowId={(status) => status.id}
      searchable={false} components={{ pagination: false }} emptyText="У чернетці ще немає статусів." />
  );
}

function editButton(status, onEdit) {
  return (
    <button className="order-status-edit" type="button" aria-label={`Змінити назву статусу «${status.name}»`}
      onClick={() => onEdit(status)}>Змінити</button>
  );
}
