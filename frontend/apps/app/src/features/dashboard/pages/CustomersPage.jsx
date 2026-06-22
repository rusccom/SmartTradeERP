import { useMemo } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { createCustomersTablePreset } from "./table/customersTablePreset";

function CustomersPage() {
  const { t } = useI18n();
  const preset = useMemo(() => createCustomersTablePreset(t), [t]);
  return <ServerListTable preset={preset} />;
}

export default CustomersPage;
