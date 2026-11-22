import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";
import BaseCourierFields from "./BaseCourierFields";

export default function BaseCourierForm({ connection, action }) {
  const { t } = useI18n();
  return (
    <form method="post" autoComplete="off" onSubmit={event => submit(event, connection, action)}>
      <fieldset disabled={action.busy}>
        <legend>{t(connection ? "baseCourier.replaceTitle" : "baseCourier.connectTitle")}</legend>
        <BaseCourierFields email={connection?.email} />
        <button className="base-courier-primary" type="submit">
          {t(connection ? "baseCourier.replace" : "baseCourier.connect")}
        </button>
      </fieldset>
    </form>
  );
}

function submit(event, connection, action) {
  event.preventDefault();
  const fields = new FormData(event.currentTarget);
  const input = { email: fields.get("email"), api_key: fields.get("api_key"), version: connection?.version };
  action.run(() => putJSON(apiPaths.baseCourierSettings, input), connection ? "baseCourier.replaced" : "baseCourier.connected");
}
