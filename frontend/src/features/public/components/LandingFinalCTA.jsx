import { Link } from "react-router-dom";

function LandingFinalCTA() {
  return (
    <section className="landing-section landing-final-cta">
      <h2>Ready to take control of your trade operations?</h2>
      <p>
        Join businesses that track inventory accurately, post documents
        instantly, and know their real margins &mdash; all in one workspace.
      </p>
      <div className="landing-final-cta-actions">
        <Link to="/register" className="primary-button lfc-primary">
          Start registration
        </Link>
        <Link to="/login" className="secondary-button lfc-secondary">
          Open workspace
        </Link>
      </div>
    </section>
  );
}

export default LandingFinalCTA;
