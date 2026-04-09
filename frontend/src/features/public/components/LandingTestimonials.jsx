import { motion } from "framer-motion";

import { useI18n } from "../../../shared/i18n/useI18n";

function LandingTestimonials() {
  const { t } = useI18n();
  const reviews = createReviews(t);
  return (
    <section className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-12">
          <h2 className="text-3xl md:text-5xl font-extrabold text-text-primary tracking-tight mb-4">
            {t("public.testimonials.title")}
          </h2>
          <p className="text-text-secondary text-lg">{t("public.testimonials.subtitle")}</p>
        </motion.div>
        <div className="grid md:grid-cols-3 gap-6">
          {reviews.map((review, index) => (
            <motion.article
              key={review.name}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className="rounded-2xl border border-border bg-bg-card p-6 flex flex-col"
            >
              <p className="text-text-secondary leading-relaxed flex-1 mb-6">&ldquo;{review.text}&rdquo;</p>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-accent/20 flex items-center justify-center text-accent font-bold text-sm">
                  {review.name.charAt(0)}
                </div>
                <div>
                  <p className="text-sm font-semibold text-text-primary">{review.name}</p>
                  <p className="text-xs text-text-muted">{review.role}</p>
                </div>
              </div>
            </motion.article>
          ))}
        </div>
      </div>
    </section>
  );
}

function createReviews(t) {
  return [
    { name: "Alex K.", role: t("public.testimonials.review1.role"), text: t("public.testimonials.review1.text") },
    { name: "Maria S.", role: t("public.testimonials.review2.role"), text: t("public.testimonials.review2.text") },
    { name: "David L.", role: t("public.testimonials.review3.role"), text: t("public.testimonials.review3.text") },
  ];
}

export default LandingTestimonials;
