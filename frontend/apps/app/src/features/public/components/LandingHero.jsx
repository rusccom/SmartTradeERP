import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { ArrowRight, BarChart3 } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingHero() {
  const { t } = useI18n();
  return (
    <section className="relative overflow-hidden pt-20 pb-24 px-6">
      <Backdrop />
      <div className="relative max-w-4xl mx-auto text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="inline-flex items-center gap-2 px-4 py-1.5 mb-6 rounded-full border border-border bg-bg-card/50 text-sm text-text-secondary"
        >
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
          {t("public.hero.badge")}
        </motion.div>
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="text-4xl sm:text-5xl md:text-7xl font-extrabold leading-[1.05] tracking-tight mb-6"
        >
          <span className="text-text-primary">{t("public.hero.titleLine1")}</span>
          <br />
          <span className="bg-gradient-to-r from-accent to-cyan bg-clip-text text-transparent">
            {t("public.hero.titleLine2")}
          </span>
        </motion.h1>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="max-w-2xl mx-auto text-lg text-text-secondary leading-relaxed mb-10"
        >
          {t("public.hero.description")}
        </motion.p>
        <HeroActions />
      </div>
      <HeroScreenshot />
    </section>
  );
}

function Backdrop() {
  return (
    <div className="absolute inset-0 pointer-events-none">
      <div className="absolute top-[-20%] left-1/2 -translate-x-1/2 w-[800px] h-[600px] rounded-full blur-[120px] bg-[rgba(99,102,241,0.15)]" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[400px] h-[400px] rounded-full blur-[100px] bg-[rgba(6,182,212,0.1)]" />
    </div>
  );
}

function HeroActions() {
  const { t } = useI18n();
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay: 0.3 }}
      className="flex flex-col sm:flex-row items-center justify-center gap-4"
    >
      <Link to="/register" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl bg-accent hover:bg-accent-hover text-white font-semibold text-base transition-all no-underline shadow-[0_0_24px_rgba(99,102,241,0.3)] hover:shadow-[0_0_32px_rgba(99,102,241,0.5)]">
        {t("public.hero.primaryCta")}
        <ArrowRight size={18} />
      </Link>
      <a href="#features" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl border border-border text-text-secondary hover:text-text-primary hover:border-text-muted font-semibold text-base transition-all no-underline">
        {t("public.hero.secondaryCta")}
      </a>
    </motion.div>
  );
}

function HeroScreenshot() {
  const { t } = useI18n();
  return (
    <motion.div
      initial={{ opacity: 0, y: 40 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, delay: 0.5 }}
      className="relative max-w-5xl mx-auto mt-16"
    >
      <div className="absolute inset-0 bg-gradient-to-t from-bg-primary via-transparent to-transparent z-10 pointer-events-none" />
      <div className="relative rounded-2xl border border-border bg-bg-card overflow-hidden shadow-2xl shadow-accent/5">
        <div className="w-full aspect-[16/9] bg-gradient-to-br from-bg-secondary to-bg-card flex items-center justify-center text-text-muted">
          <div className="text-center">
            <div className="w-16 h-16 rounded-2xl bg-bg-secondary border border-border mx-auto mb-4 flex items-center justify-center text-2xl">
              <BarChart3 size={28} />
            </div>
            <p className="text-sm font-medium">{t("public.hero.screenshotTitle")}</p>
            <p className="text-xs mt-1">{t("public.hero.screenshotText")}</p>
          </div>
        </div>
      </div>
    </motion.div>
  );
}

export default LandingHero;
