import { motion } from "framer-motion";
import { BarChart3, FileText, TrendingUp, Warehouse } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingFeatures() {
  const { t } = useI18n();
  const features = createFeatures(t);
  return (
    <section id="features" className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="mb-12">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.features.title")}
          </h2>
          <p className="text-text-secondary text-lg max-w-xl">{t("public.features.subtitle")}</p>
        </motion.div>
        <div className="grid md:grid-cols-4 gap-4">
          {features.map((item, index) => (
            <motion.article
              key={item.title}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className={`group relative rounded-2xl border border-border bg-bg-card p-6 transition-all duration-300 hover:border-accent/40 hover:shadow-[0_0_30px_rgba(99,102,241,0.08)] ${item.span}`}
            >
              <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center mb-4 group-hover:bg-accent/20 transition-colors">
                <item.icon size={20} className="text-accent" />
              </div>
              <h3 className="text-base font-bold text-text-primary mb-2">{item.title}</h3>
              <p className="text-sm text-text-secondary leading-relaxed">{item.text}</p>
            </motion.article>
          ))}
        </div>
      </div>
    </section>
  );
}

function createFeatures(t) {
  return [
    { icon: Warehouse, title: t("public.features.card1.title"), text: t("public.features.card1.text"), span: "md:col-span-2 md:row-span-2" },
    { icon: FileText, title: t("public.features.card2.title"), text: t("public.features.card2.text"), span: "" },
    { icon: TrendingUp, title: t("public.features.card3.title"), text: t("public.features.card3.text"), span: "" },
    { icon: BarChart3, title: t("public.features.card4.title"), text: t("public.features.card4.text"), span: "md:col-span-2" },
  ];
}

export default LandingFeatures;
