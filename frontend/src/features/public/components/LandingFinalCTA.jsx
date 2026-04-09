import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { ArrowRight } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingFinalCTA() {
  const { t } = useI18n();
  return (
    <section className="py-24 px-6">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="relative max-w-4xl mx-auto rounded-3xl border border-border overflow-hidden"
      >
        <CalloutBackdrop />
        <div className="relative text-center py-16 px-8">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.final.title")}
          </h2>
          <p className="max-w-lg mx-auto text-text-secondary text-lg mb-8">{t("public.final.description")}</p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              to="/register"
              className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl bg-accent hover:bg-accent-hover text-white font-semibold text-base transition-all no-underline shadow-[0_0_24px_rgba(99,102,241,0.3)] hover:shadow-[0_0_32px_rgba(99,102,241,0.5)]"
            >
              {t("public.final.primaryCta")}
              <ArrowRight size={18} />
            </Link>
            <Link
              to="/login"
              className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl border border-border text-text-secondary hover:text-text-primary hover:border-text-muted font-semibold text-base transition-all no-underline"
            >
              {t("public.final.secondaryCta")}
            </Link>
          </div>
        </div>
      </motion.div>
    </section>
  );
}

function CalloutBackdrop() {
  return (
    <div className="absolute inset-0 pointer-events-none">
      <div className="absolute top-[-40%] left-[-10%] w-[500px] h-[500px] rounded-full blur-[120px] bg-[rgba(99,102,241,0.15)]" />
      <div className="absolute bottom-[-40%] right-[-10%] w-[400px] h-[400px] rounded-full blur-[100px] bg-[rgba(6,182,212,0.1)]" />
    </div>
  );
}

export default LandingFinalCTA;
