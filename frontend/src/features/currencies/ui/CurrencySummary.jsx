import { CheckCircle2 } from "lucide-react";

function CurrencySummary({ currency, emptyLabel, formatMoney, loadingLabel }) {
  if (!currency) {
    return <div className="currency-summary currency-summary--empty">{loadingLabel || emptyLabel}</div>;
  }
  return (
    <div className="currency-summary">
      <div className="currency-summary-main">
        <span className="currency-summary-code">{currency.code}</span>
        <span className="currency-summary-name">{currency.name}</span>
      </div>
      <div className="currency-summary-meta">
        <span>{formatMoney(100)}</span>
        <CheckCircle2 size={16} />
      </div>
    </div>
  );
}

export default CurrencySummary;
