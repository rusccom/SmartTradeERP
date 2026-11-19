import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import InPostTokenField from "./InPostTokenField";

export default function InPostFields({ connection }) {
  const { t } = useI18n();
  const [env, setEnv] = useState(connection?.environment || "production");
  const [id, setID] = useState(connection?.organization_id || "");
  const reuse = env === connection?.environment && id === connection?.organization_id;
  return (
    <>
      <div className="inpost-account">
        <label>{t("inpost.environment")}<select name="environment" value={env} onChange={e => setEnv(e.target.value)}>
          <option value="production">{t("inpost.production")}</option>
          <option value="sandbox">{t("inpost.sandbox")}</option>
        </select></label>
        <label>{t("inpost.organization")}<input name="organization_id" value={id} onChange={e => setID(e.target.value)}
          inputMode="numeric" pattern="[1-9][0-9]*" maxLength={19} required /></label>
      </div>
      <InPostTokenField key={env} reuse={reuse} />
    </>
  );
}
