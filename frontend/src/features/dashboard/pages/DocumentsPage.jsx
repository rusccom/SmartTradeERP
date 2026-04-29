import { useMemo } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { createDocumentsTablePreset } from "./table/documentsTablePreset";

function DocumentsPage() {
  const { t } = useI18n();
  const preset = useMemo(() => createDocumentsTablePreset(t), [t]);
  return <ServerListTable preset={preset} />;
}

export default DocumentsPage;
