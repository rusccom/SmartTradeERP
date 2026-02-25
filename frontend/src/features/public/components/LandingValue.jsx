const POINTS = [
  {
    title: "No per-warehouse fees",
    desc: "Scale to as many locations as you need without extra charges.",
  },
  {
    title: "Role-based access",
    desc: "Owner, manager, and cashier permissions are built in from day one.",
  },
  {
    title: "Cloud-native",
    desc: "Access from anywhere. Your data is backed up and isolated automatically.",
  },
];

function LandingValue() {
  return (
    <section className="landing-section landing-value">
      <div className="landing-section-inner">
        <h2 className="landing-section-title">
          Built for small businesses that think big
        </h2>
        <p className="landing-section-subtitle">
          SmartTrade ERP gives you the inventory control, cost accuracy, and
          financial visibility that used to require enterprise software &mdash; at
          a fraction of the complexity.
        </p>
        <div className="landing-value-grid">
          {POINTS.map((p) => (
            <div key={p.title} className="landing-value-point">
              <span className="landing-value-check" aria-hidden="true" />
              <div>
                <strong>{p.title}</strong>
                <span>{p.desc}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

export default LandingValue;
