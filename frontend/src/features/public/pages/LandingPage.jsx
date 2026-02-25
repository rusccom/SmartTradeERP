import LandingHero from "../components/LandingHero";
import LandingStats from "../components/LandingStats";
import LandingFeatures from "../components/LandingFeatures";
import LandingHowItWorks from "../components/LandingHowItWorks";
import LandingShowcase from "../components/LandingShowcase";
import LandingValue from "../components/LandingValue";
import LandingFinalCTA from "../components/LandingFinalCTA";

import "../styles/landing-base.css";
import "../styles/landing-hero.css";
import "../styles/landing-stats.css";
import "../styles/landing-features.css";
import "../styles/landing-how.css";
import "../styles/landing-showcase.css";
import "../styles/landing-value.css";
import "../styles/landing-cta.css";

function LandingPage() {
  return (
    <div className="landing-stack">
      <LandingHero />
      <LandingStats />
      <LandingFeatures />
      <LandingHowItWorks />
      <LandingShowcase />
      <LandingValue />
      <LandingFinalCTA />
    </div>
  );
}

export default LandingPage;
