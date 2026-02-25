const BARS = [
  { size: "86%", delay: "0s" },
  { size: "72%", delay: "0.12s" },
  { size: "64%", delay: "0.24s" },
];

const KPI = [
  { label: "Order throughput", value: "x2.7" },
  { label: "Inventory accuracy", value: "99.4%" },
  { label: "Document flow uptime", value: "24/7" },
];

function LandingPreview() {
  return (
    <aside className="landing-preview" aria-hidden="true">
      <div className="landing-preview-header">
        <span className="landing-live-dot" />
        <p className="landing-preview-title">Operational pulse</p>
      </div>
      <div className="landing-bars">
        {BARS.map((bar) => (
          <span
            key={bar.size}
            className="landing-bar"
            style={{ "--bar-size": bar.size, "--bar-delay": bar.delay }}
          />
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
