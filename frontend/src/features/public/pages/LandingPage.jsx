import LandingHero from "../components/LandingHero";
import LandingCapabilities from "../components/LandingCapabilities";
import "../styles/landing-page.css";

function LandingPage() {
  return (
    <div className="landing-stack">
      <LandingHero />
      <LandingCapabilities />
    </div>
  );
}

export default LandingPage;

