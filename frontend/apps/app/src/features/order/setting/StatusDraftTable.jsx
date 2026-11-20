import { useI18n } from "@smarterp/i18n/useI18n";
import { DataTable } from "@smarterp/ui/tables/list-table";

export default function StatusDraftTable({ statuses, onEdit }) {
  const { t } = useI18n();
  const columns = [
    { accessorKey: "name", header: t("orders.statuses.name") },
    { id: "actions", header: t("orders.statuses.actions"), size: 116, cell: (_, status) => editButton(status, onEdit, t) },
  ];
  return (
    <DataTable data={statuses} columns={columns} getRowId={(status) => status.id}
      searchable={false} components={{ pagination: false }} emptyText={t("orders.statuses.empty")} />
  );
}

function editButton(status, onEdit, t) {
  return (
    <button className="order-status-edit" type="button" aria-label={t("orders.statuses.editName", { name: status.name })}
      onClick={() => onEdit(status)}>{t("orders.statuses.edit")}</button>
  );
}
