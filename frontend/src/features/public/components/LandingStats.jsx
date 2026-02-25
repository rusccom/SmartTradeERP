const STATS = [
  {
    icon: "\u25A3",
    title: "Unlimited warehouses",
    desc: "Real-time stock across all locations",
  },
  {
    icon: "\u25CE",
    title: "Penny-accurate costing",
    desc: "Automatic weighted average cost",
  },
  {
    icon: "\u2630",
    title: "5 document types",
    desc: "Full draft-to-posted lifecycle",
  },
  {
    icon: "\u2713",
    title: "Multi-tenant isolation",
    desc: "Each business gets private data",
  },
];

function LandingStats() {
  return (
    <section className="landing-section landing-stats">
      {STATS.map((s) => (
        <div key={s.title} className="landing-stat-item">
          <span className="landing-stat-icon">{s.icon}</span>
          <strong className="landing-stat-title">{s.title}</strong>
          <span className="landing-stat-desc">{s.desc}</span>
        </div>
      ))}
    </section>
  );
}

export default LandingStats;
