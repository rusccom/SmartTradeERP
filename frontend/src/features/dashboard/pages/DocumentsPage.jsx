import { useMemo } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { useCurrencies } from "../../currencies/model/useCurrencies";
import { createDocumentsTablePreset } from "./table/documentsTablePreset";

function DocumentsPage() {
  const { t } = useI18n();
  const { formatMoney } = useCurrencies();
  const preset = useMemo(() => createDocumentsTablePreset(t, formatMoney), [formatMoney, t]);
  return <ServerListTable preset={preset} />;
}

export default DocumentsPage;
