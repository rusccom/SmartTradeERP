import { motion } from "framer-motion";

/*
  ОТЗЫВЫ — замените данные ниже на реальные.
  Для аватаров: поместите файлы в /public/avatars/ и укажите путь в поле avatar.
  Пример: avatar: "/avatars/john.jpg"
*/
const REVIEWS = [
  {
    name: "Alex K.",
    role: "Owner, RetailPro",
    text: "SmartTrade replaced three different tools we used. Inventory accuracy went from ~85% to 99.4% in the first month.",
    avatar: null,
  },
  {
    name: "Maria S.",
    role: "Operations Manager",
    text: "The automatic cost recalculation alone saves us 10+ hours per week. No more manual spreadsheet reconciliations.",
    avatar: null,
  },
  {
    name: "David L.",
    role: "CFO, TradeGroup",
    text: "Finally, real profit margins we can trust. The weighted average costing is exactly what we needed for multi-warehouse operations.",
    avatar: null,
  },
];

function LandingTestimonials() {
  return (
    <section className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <h2
            className="text-3xl md:text-5xl font-extrabold
              text-text-primary tracking-tight mb-4"
          >
            Trusted by real businesses
          </h2>
          <p className="text-text-secondary text-lg">
            See what our customers say about SmartTrade ERP.
          </p>
        </motion.div>

        <div className="grid md:grid-cols-3 gap-6">
          {REVIEWS.map((r, i) => (
            <motion.article
              key={r.name}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="rounded-2xl border border-border bg-bg-card
                p-6 flex flex-col"
            >
              <p
                className="text-text-secondary leading-relaxed
                  flex-1 mb-6"
              >
                &ldquo;{r.text}&rdquo;
              </p>
              <div className="flex items-center gap-3">
                {/*
                  АВАТАР — замените на:
                  <img src={r.avatar} alt={r.name}
                    className="w-10 h-10 rounded-full object-cover" />
                */}
                <div
                  className="w-10 h-10 rounded-full bg-accent/20
                    flex items-center justify-center text-accent
                    font-bold text-sm"
                >
                  {r.name.charAt(0)}
                </div>
                <div>
                  <p className="text-sm font-semibold text-text-primary">
                    {r.name}
                  </p>
                  <p className="text-xs text-text-muted">{r.role}</p>
                </div>
              </div>
            </motion.article>
          ))}
        </div>
      </div>
    </section>
  );
}

export default LandingTestimonials;
