import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Check } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingPricing() {
  const { t } = useI18n();
  const plans = createPlans(t);
  return (
    <section id="pricing" className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-12">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.pricing.title")}
          </h2>
          <p className="text-text-secondary text-lg">{t("public.pricing.subtitle")}</p>
        </motion.div>
        <div className="grid md:grid-cols-3 gap-6">
          {plans.map((plan, index) => (
            <motion.div
              key={plan.name}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className={`rounded-2xl border p-6 flex flex-col ${plan.highlighted ? "border-accent bg-accent/5 shadow-[0_0_40px_rgba(99,102,241,0.1)]" : "border-border bg-bg-card"}`}
            >
              {plan.highlighted && <span className="inline-block self-start text-xs font-bold text-accent bg-accent/10 px-3 py-1 rounded-full mb-4">{t("public.pricing.highlight")}</span>}
              <h3 className="text-lg font-bold text-text-primary">{plan.name}</h3>
              <div className="mt-3 mb-2">
                <span className="text-4xl font-extrabold text-text-primary">{plan.price}</span>
                {plan.period && <span className="text-text-muted text-sm">{plan.period}</span>}
              </div>
              <p className="text-sm text-text-secondary mb-6">{plan.desc}</p>
              <ul className="flex flex-col gap-3 mb-8 flex-1">
                {plan.features.map((feature) => (
                  <li key={feature} className="flex items-center gap-2 text-sm text-text-secondary">
                    <Check size={16} className="text-accent flex-shrink-0" />
                    {feature}
                  </li>
                ))}
              </ul>
              <Link
                to="/register"
                className={`w-full py-3 rounded-xl font-semibold text-sm text-center no-underline transition-all ${plan.highlighted ? "bg-accent hover:bg-accent-hover text-white shadow-[0_0_20px_rgba(99,102,241,0.3)]" : "border border-border text-text-secondary hover:text-text-primary hover:border-text-muted"}`}
              >
                {plan.cta}
              </Link>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}

function createPlans(t) {
  return [
    {
      name: t("public.pricing.starter.name"),
      price: t("public.pricing.starter.price"),
      period: "",
      desc: t("public.pricing.starter.desc"),
      features: [
        t("public.pricing.starter.feature1"),
        t("public.pricing.starter.feature2"),
        t("public.pricing.starter.feature3"),
        t("public.pricing.starter.feature4"),
      ],
      cta: t("public.pricing.starter.cta"),
      highlighted: false,
    },
    {
      name: t("public.pricing.pro.name"),
      price: t("public.pricing.pro.price"),
      period: t("public.pricing.pro.period"),
      desc: t("public.pricing.pro.desc"),
      features: [
        t("public.pricing.pro.feature1"),
        t("public.pricing.pro.feature2"),
        t("public.pricing.pro.feature3"),
        t("public.pricing.pro.feature4"),
        t("public.pricing.pro.feature5"),
        t("public.pricing.pro.feature6"),
      ],
      cta: t("public.pricing.pro.cta"),
      highlighted: true,
    },
    {
      name: t("public.pricing.enterprise.name"),
      price: t("public.pricing.enterprise.price"),
      period: "",
      desc: t("public.pricing.enterprise.desc"),
      features: [
        t("public.pricing.enterprise.feature1"),
        t("public.pricing.enterprise.feature2"),
        t("public.pricing.enterprise.feature3"),
        t("public.pricing.enterprise.feature4"),
        t("public.pricing.enterprise.feature5"),
      ],
      cta: t("public.pricing.enterprise.cta"),
      highlighted: false,
    },
  ];
}

export default LandingPricing;
