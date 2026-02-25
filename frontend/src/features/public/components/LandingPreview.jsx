const BARS = ["86%", "72%", "64%"];

const KPI = [
  { label: "Order throughput", value: "x2.7" },
  { label: "Inventory accuracy", value: "99.4%" },
  { label: "Document flow uptime", value: "24/7" },
];

function LandingPreview() {
  return (
    <aside className="landing-preview" aria-hidden="true">
      <p className="landing-preview-title">Operational pulse</p>
      <div className="landing-bars">
        {BARS.map((size) => (
          <span key={size} className="landing-bar" style={{ "--bar-size": size }} />
        ))}
      </div>
      <div className="landing-kpi">
        {KPI.map((item) => (
          <article key={item.label} className="landing-kpi-item">
            <span>{item.label}</span>
            <strong>{item.value}</strong>
          </article>
        ))}
      </div>
    </aside>
  );
}

export default LandingPreview;
