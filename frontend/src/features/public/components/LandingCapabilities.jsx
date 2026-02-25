const CAPABILITIES = [
  {
    title: "Inventory and warehousing",
    text: "Receipts, write-offs, transfers, and reservations are tracked in one flow with full history.",
  },
  {
    title: "Documents and accounting",
    text: "The system links documents to each operation and prepares data for postings and period closing.",
  },
  {
    title: "Analytics for growth",
    text: "Margin, turnover, and sales trend reports help teams make faster operational decisions.",
  },
];

function LandingCapabilities() {
  return (
    <section className="landing-capabilities">
      <h2>What your business gets after launch</h2>
      <div className="landing-grid">
        {CAPABILITIES.map((item) => (
          <article key={item.title} className="landing-card">
            <h3>{item.title}</h3>
            <p>{item.text}</p>
          </article>
        ))}
      </div>
    </section>
  );
}

export default LandingCapabilities;
