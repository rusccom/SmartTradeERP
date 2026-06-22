import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronDown } from "lucide-react";

import { useI18n } from "../../../shared/i18n/useI18n";

function FAQItem({ item }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="border-b border-border last:border-b-0">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center justify-between py-5 text-left bg-transparent border-none cursor-pointer text-text-primary font-semibold text-base hover:text-accent transition-colors"
      >
        {item.q}
        <ChevronDown size={18} className={`text-text-muted flex-shrink-0 ml-4 transition-transform duration-200 ${open ? "rotate-180" : ""}`} />
      </button>
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <p className="pb-5 text-sm text-text-secondary leading-relaxed">{item.a}</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function LandingFAQ() {
  const { t } = useI18n();
  const items = createItems(t);
  return (
    <section className="py-24 px-6">
      <div className="max-w-2xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-12">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.faq.title")}
          </h2>
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="rounded-2xl border border-border bg-bg-card px-6">
          {items.map((item) => (
            <FAQItem key={item.q} item={item} />
          ))}
        </motion.div>
      </div>
    </section>
  );
}

function createItems(t) {
  return [
    { q: t("public.faq.item1.question"), a: t("public.faq.item1.answer") },
    { q: t("public.faq.item2.question"), a: t("public.faq.item2.answer") },
    { q: t("public.faq.item3.question"), a: t("public.faq.item3.answer") },
    { q: t("public.faq.item4.question"), a: t("public.faq.item4.answer") },
    { q: t("public.faq.item5.question"), a: t("public.faq.item5.answer") },
  ];
}

export default LandingFAQ;
