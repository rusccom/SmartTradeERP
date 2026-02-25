const ROWS = [
  {
    title: "Bundle products with automatic cost distribution",
    text: "Create composite products \u2014 bundles, kits, or recipes \u2014 and sell them as a single item. SmartTrade automatically distributes revenue proportionally across components based on their weighted average cost, giving you accurate per-component profit without manual math.",
    visual: "composite",
  },
  {
    title: "Edit posted documents without breaking your books",
    text: "Made a mistake in last week\u2019s receipt? Edit it directly. SmartTrade recalculates every affected ledger entry, cost figure, and downstream transaction automatically within a single atomic operation. Your books stay accurate without manual corrections or reversal documents.",
    visual: "retroactive",
  },
  {
    title: "Manage stock across every location from one screen",
    text: "Track inventory in unlimited warehouses. Transfer stock between locations with full audit trails. Run physical inventory counts per warehouse and reconcile discrepancies instantly. See consolidated or per-location stock levels and average costs in real time.",
    visual: "warehouse",
  },
];

function ShowcaseVisual({ type }) {
  if (type === "composite") return <CompositeVisual />;
  if (type === "retroactive") return <RetroactiveVisual />;
  return <WarehouseVisual />;
}

function CompositeVisual() {
  return (
    <div className="showcase-visual showcase-composite">
      <div className="sc-parent">Bundle</div>
      <div className="sc-arrows">
        <span className="sc-line" />
        <span className="sc-line" />
        <span className="sc-line" />
      </div>
      <div className="sc-children">
        <span className="sc-child">Item A<em>40%</em></span>
        <span className="sc-child">Item B<em>35%</em></span>
        <span className="sc-child">Item C<em>25%</em></span>
      </div>
    </div>
  );
}

function RetroactiveVisual() {
  return (
    <div className="showcase-visual showcase-retro">
      <span className="sr-node sr-doc">DOC</span>
      <span className="sr-arrow">\u21BA</span>
      <span className="sr-node sr-edit">EDIT</span>
      <span className="sr-arrow">\u2192</span>
      <span className="sr-node sr-ledger">RECALC</span>
    </div>
  );
}

function WarehouseVisual() {
  return (
    <div className="showcase-visual showcase-wh">
      <span className="sw-box">WH-1<em>124</em></span>
      <span className="sw-link">\u21C4</span>
      <span className="sw-box">WH-2<em>87</em></span>
      <span className="sw-link">\u21C4</span>
      <span className="sw-box">WH-3<em>53</em></span>
    </div>
  );
}

function LandingShowcase() {
  return (
    <section className="landing-section landing-showcase">
      <div className="landing-section-inner">
        <h2 className="landing-section-title">Built for real trade complexity</h2>
        <p className="landing-section-subtitle">
          Features designed for how retail and distribution actually works.
        </p>
        <div className="showcase-rows">
          {ROWS.map((row, i) => (
            <div
              key={row.visual}
              className={`showcase-row ${i % 2 === 1 ? "showcase-row--reversed" : ""}`}
            >
              <div className="showcase-text">
                <h3>{row.title}</h3>
                <p>{row.text}</p>
              </div>
              <ShowcaseVisual type={row.visual} />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

export default LandingShowcase;
