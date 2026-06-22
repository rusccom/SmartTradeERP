import LandingHero from "../components/LandingHero";
import LandingStats from "../components/LandingStats";
import LandingFeatures from "../components/LandingFeatures";
import LandingHowItWorks from "../components/LandingHowItWorks";
import LandingShowcase from "../components/LandingShowcase";
import LandingTestimonials from "../components/LandingTestimonials";
import LandingPricing from "../components/LandingPricing";
import LandingFAQ from "../components/LandingFAQ";
import LandingFinalCTA from "../components/LandingFinalCTA";
import LandingFooter from "../components/LandingFooter";

function LandingPage() {
  return (
    <div>
      <LandingHero />
      <LandingStats />
      <LandingFeatures />
      <LandingHowItWorks />
      <LandingShowcase />
      <LandingTestimonials />
      <LandingPricing />
      <LandingFAQ />
      <LandingFinalCTA />
      <LandingFooter />
    </div>
  );
}

export default LandingPage;
