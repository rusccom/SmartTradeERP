const FEATURES = [
  {
    icon: "\u2302",
    title: "Inventory & Warehousing",
    text: "Track stock across multiple warehouses in real time. Process receipts from suppliers, record write-offs, move inventory between locations with inter-warehouse transfers, and run full physical inventory counts \u2014 all from one interface.",
  },
  {
    icon: "\u2637",
    title: "Documents & Workflow",
    text: "Five document types \u2014 Receipt, Sale, WriteOff, Transfer, and Inventory \u2014 follow a clear lifecycle: Draft, Posted, Cancelled. Every posting updates your cost ledger automatically. No manual journal entries needed.",
  },
  {
    icon: "\u2197",
    title: "Cost Accounting & Profit",
    text: "Automatic weighted average cost (AVCost) calculation on every receipt. See true cost of goods sold per transaction and real profit margins based on actual numbers, not estimates or assumptions.",
  },
  {
    icon: "\u2261",
    title: "Reports & Analytics",
    text: "Profit reports by date range, stock summaries per warehouse, top-selling products ranked by margin, and complete movement history for any item. Filter, drill down, and export the data you need.",
  },
];

function LandingFeatures() {
  return (
    <section className="landing-section landing-features">
      <h2 className="landing-section-title">
        Everything your business needs to trade smarter
      </h2>
      <p className="landing-section-subtitle">
        From receiving stock to calculating profit, SmartTrade ERP covers every
        step of your daily trade operations.
      </p>
      <div className="landing-features-grid">
        {FEATURES.map((f) => (
          <article key={f.title} className="landing-feature-card">
            <span className="landing-feature-icon">{f.icon}</span>
            <h3>{f.title}</h3>
            <p>{f.text}</p>
          </article>
        ))}
      </div>
    </section>
  );
}

export default LandingFeatures;
