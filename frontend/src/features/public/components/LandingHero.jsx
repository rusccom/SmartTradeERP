import { Link } from "react-router-dom";

import LandingPreview from "./LandingPreview";

function LandingHero() {
  return (
    <section className="landing-section landing-hero">
      <div className="landing-section-inner landing-hero-inner">
      <div className="landing-copy">
        <p className="landing-kicker">Cloud ERP for Retail &amp; Distribution</p>
        <h1>Track every item. Post every document. Know your real profit.</h1>
        <p className="landing-lead">
          Manage inventory across multiple warehouses, process sales and receipts
          with automatic cost accounting, and see real profit margins — all from
          one workspace. No spreadsheets, no guesswork.
        </p>
        <div className="landing-actions">
          <Link to="/register" className="primary-button">
            Start free
          </Link>
          <Link to="/login" className="secondary-button">
            Open workspace
          </Link>
        </div>
      </div>
      <LandingPreview />
      </div>
    </section>
  );
}

export default LandingHero;
