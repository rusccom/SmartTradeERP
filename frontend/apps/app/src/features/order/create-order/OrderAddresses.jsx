import { useI18n } from "@smarterp/i18n/useI18n";
import BoundField from "./BoundField";
import BillingCustomer from "./BillingCustomer";
import { useCountryOptions } from "@smarterp/model/countryOptions";

const address = ["name", "company", "street", "postal_code", "city", "region", "country"];
const pickup = ["name", "point_id", "street", "postal_code", "city", "country"];

export default function OrderAddresses({ editor }) {
  const { t, locale } = useI18n();
  const dictionary = useCountryOptions(locale);
  const countries = [["", "—"], ...dictionary.options];
  const groups = [["shipping", address], ["billing", [...address, "tax_id"]]];
  if (editor.form.details.delivery_kind === "pickup") groups.push(["pickup", pickup]);
  return <div className="order-grid">{dictionary.error && <p role="alert">{dictionary.error} <button type="button" onClick={dictionary.retry}>{t("entry.retry")}</button></p>}{groups.map(([section, fields]) =>
    <section className="order-card" key={section} aria-label={t(`entry.${section}`)}><h2>{t(`entry.${section}`)}</h2>
      {fields.map((name) => <BoundField key={name} editor={editor} path={`details.${section}.${name}`}
        label={t(section === "pickup" && name === "name" ? "entry.pointName" : `entry.${name}`)}
        options={name === "country" ? countries : undefined} required={name === "point_id"} />)}
      {section === "billing" && <BillingCustomer editor={editor} />}
    </section>)}</div>;
}
