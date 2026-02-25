import { Link } from "react-router-dom";

import LandingPreview from "./LandingPreview";

function LandingHero() {
  return (
    <section className="landing-hero">
      <div className="landing-copy">
        <p className="landing-kicker">SmartTrade ERP for retail and distribution</p>
        <h1>One workspace for sales, inventory, documents, and operational reporting.</h1>
        <p className="landing-lead">
          Teams track stock movement in real time, accounting gets clean documents, and management controls margin and
          turnover without manual spreadsheets.
        </p>
        <div className="landing-actions">
          <Link to="/register" className="primary-button">
            Start registration
          </Link>
          <Link to="/login" className="secondary-button">
            Open workspace
          </Link>
        </div>
      </div>
      <LandingPreview />
    </section>
  );
}

export default LandingHero;
