const STEPS = [
  {
    num: "01",
    title: "Create your workspace",
    text: "Register in 30 seconds. Your isolated tenant is provisioned instantly with role-based access for owners, managers, and cashiers.",
  },
  {
    num: "02",
    title: "Set up your catalog",
    text: "Add products, define variants and composite bundles, assign warehouses. Import existing data or start fresh \u2014 your choice.",
  },
  {
    num: "03",
    title: "Start trading",
    text: "Process sales, receipts, and transfers. The system handles cost calculations, ledger postings, and profit tracking automatically.",
  },
];

function LandingHowItWorks() {
  return (
    <section className="landing-section landing-how">
      <h2 className="landing-section-title">
        Up and running in three steps
      </h2>
      <p className="landing-section-subtitle">
        No complex setup, no consultants. Get your team productive in minutes.
      </p>
      <div className="landing-how-grid">
        {STEPS.map((s) => (
          <article key={s.num} className="landing-step">
            <span className="landing-step-num">{s.num}</span>
            <h3>{s.title}</h3>
            <p>{s.text}</p>
          </article>
        ))}
      </div>
    </section>
  );
}

export default LandingHowItWorks;
