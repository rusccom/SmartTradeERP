import { motion } from "framer-motion";

import { useI18n } from "../../../shared/i18n/useI18n";

const LOGOS = ["Partner 1", "Partner 2", "Partner 3", "Partner 4", "Partner 5"];

function LandingStats() {
  const { t } = useI18n();
  const stats = createStats(t);
  return (
    <section className="py-16 px-6 border-y border-border">
      <div className="max-w-5xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
          {stats.map((item, index) => (
            <motion.div
              key={item.label}
              initial={{ opacity: 0, y: 16 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className="text-center"
            >
              <p className="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-accent to-cyan bg-clip-text text-transparent">
                {item.value}
              </p>
              <p className="text-sm text-text-muted mt-1">{item.label}</p>
            </motion.div>
          ))}
        </div>

        <div className="flex flex-wrap items-center justify-center gap-8">
          {LOGOS.map((name) => (
            <span key={name} className="px-4 py-2 rounded-lg border border-border-subtle bg-bg-card/50 text-text-muted text-sm font-medium">
              {name}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}

function createStats(t) {
  return [
    { value: "500+", label: t("public.stats.item1") },
    { value: "99.9%", label: t("public.stats.item2") },
    { value: "<50ms", label: t("public.stats.item3") },
    { value: "24/7", label: t("public.stats.item4") },
  ];
}

export default LandingStats;
