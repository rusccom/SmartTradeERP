import { motion } from "framer-motion";
import { MapPin, Package, PenLine } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingShowcase() {
  const { t } = useI18n();
  const rows = createRows(t);
  return (
    <section className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-16">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.showcase.title")}
          </h2>
          <p className="text-text-secondary text-lg max-w-xl mx-auto">{t("public.showcase.subtitle")}</p>
        </motion.div>
        <div className="flex flex-col gap-20">
          {rows.map((row, index) => (
            <motion.div
              key={row.title}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="grid md:grid-cols-2 gap-10 items-center"
            >
              <div className={index % 2 === 1 ? "md:order-2" : ""}>
                <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center mb-4">
                  <row.icon size={20} className="text-accent" />
                </div>
                <h3 className="text-xl md:text-2xl font-bold text-text-primary mb-3">{row.title}</h3>
                <p className="text-text-secondary leading-relaxed">{row.text}</p>
              </div>
              <div className={`rounded-2xl border border-border bg-bg-card overflow-hidden ${index % 2 === 1 ? "md:order-1" : ""}`}>
                <div className="aspect-[4/3] bg-gradient-to-br from-bg-secondary to-bg-card flex items-center justify-center text-text-muted">
                  <div className="text-center px-6">
                    <row.icon size={32} className="mx-auto mb-3 text-border" />
                    <p className="text-xs">{row.imgHint}</p>
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}

function createRows(t) {
  return [
    { icon: Package, title: t("public.showcase.row1.title"), text: t("public.showcase.row1.text"), imgHint: t("public.showcase.row1.hint") },
    { icon: PenLine, title: t("public.showcase.row2.title"), text: t("public.showcase.row2.text"), imgHint: t("public.showcase.row2.hint") },
    { icon: MapPin, title: t("public.showcase.row3.title"), text: t("public.showcase.row3.text"), imgHint: t("public.showcase.row3.hint") },
  ];
}

export default LandingShowcase;
