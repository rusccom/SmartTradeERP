import { motion } from "framer-motion";
import { PackagePlus, Rocket, UserPlus } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingHowItWorks() {
  const { t } = useI18n();
  const steps = createSteps(t);
  return (
    <section id="how-it-works" className="py-24 px-6">
      <div className="max-w-3xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-16">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.how.title")}
          </h2>
          <p className="text-text-secondary text-lg">{t("public.how.subtitle")}</p>
        </motion.div>
        <div className="relative">
          <div className="absolute left-5 md:left-6 top-0 bottom-0 w-px bg-gradient-to-b from-accent via-cyan to-transparent" />
          <div className="flex flex-col gap-12">
            {steps.map((step, index) => (
              <motion.div
                key={step.num}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.15 }}
                className="relative flex gap-6 items-start"
              >
                <div className="relative z-10 flex-shrink-0 w-10 h-10 md:w-12 md:h-12 rounded-xl bg-bg-card border border-border flex items-center justify-center">
                  <step.icon size={20} className="text-accent" />
                </div>
                <div className="pt-1">
                  <span className="text-xs font-bold text-accent tracking-widest">
                    {t("public.how.step")} {step.num}
                  </span>
                  <h3 className="text-lg font-bold text-text-primary mt-1 mb-2">{step.title}</h3>
                  <p className="text-sm text-text-secondary leading-relaxed">{step.text}</p>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

function createSteps(t) {
  return [
    { icon: UserPlus, num: "01", title: t("public.how.step1.title"), text: t("public.how.step1.text") },
    { icon: PackagePlus, num: "02", title: t("public.how.step2.title"), text: t("public.how.step2.text") },
    { icon: Rocket, num: "03", title: t("public.how.step3.title"), text: t("public.how.step3.text") },
  ];
}

export default LandingHowItWorks;
