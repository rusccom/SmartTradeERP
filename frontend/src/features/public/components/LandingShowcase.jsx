import { motion } from "framer-motion";
import { Package, PenLine, MapPin } from "lucide-react";

const ROWS = [
  {
    icon: Package,
    title: "Bundle products with automatic cost distribution",
    text: "Create composite products — bundles, kits, or recipes — and sell them as a single item. SmartTrade distributes revenue across components based on weighted average cost automatically.",
    imgHint: "Place a screenshot of the bundle/composite product editor here",
  },
  {
    icon: PenLine,
    title: "Edit posted documents without breaking your books",
    text: "Made a mistake in last week's receipt? Edit it directly. SmartTrade recalculates every affected ledger entry and downstream transaction within a single atomic operation.",
    imgHint: "Place a screenshot of the document editing interface here",
  },
  {
    icon: MapPin,
    title: "Manage stock across every location from one screen",
    text: "Track inventory in unlimited warehouses. Transfer stock with full audit trails. Run physical counts per warehouse and reconcile discrepancies instantly.",
    imgHint: "Place a screenshot of the multi-warehouse stock view here",
  },
];

function LandingShowcase() {
  return (
    <section className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2
            className="text-3xl md:text-5xl font-extrabold
              text-text-primary tracking-tight mb-4"
          >
            Built for real trade complexity
          </h2>
          <p className="text-text-secondary text-lg max-w-xl mx-auto">
            Features designed for how retail and distribution actually
            works.
          </p>
        </motion.div>

        <div className="flex flex-col gap-20">
          {ROWS.map((row, i) => (
            <motion.div
              key={row.title}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className={`grid md:grid-cols-2 gap-10 items-center
                ${i % 2 === 1 ? "md:direction-rtl" : ""}`}
            >
              <div className={i % 2 === 1 ? "md:order-2" : ""}>
                <div
                  className="w-10 h-10 rounded-xl bg-accent/10
                    flex items-center justify-center mb-4"
                >
                  <row.icon size={20} className="text-accent" />
                </div>
                <h3
                  className="text-xl md:text-2xl font-bold
                    text-text-primary mb-3"
                >
                  {row.title}
                </h3>
                <p
                  className="text-text-secondary leading-relaxed"
                >
                  {row.text}
                </p>
              </div>

              {/*
                СКРИНШОТ ФИЧИ — замените этот блок на:
                <img src="/screenshots/feature-name.png"
                  alt={row.title}
                  className="w-full rounded-2xl" />
              */}
              <div
                className={`rounded-2xl border border-border bg-bg-card
                  overflow-hidden ${i % 2 === 1 ? "md:order-1" : ""}`}
              >
                <div
                  className="aspect-[4/3] bg-gradient-to-br
                    from-bg-secondary to-bg-card flex items-center
                    justify-center text-text-muted"
                >
                  <div className="text-center px-6">
                    <row.icon
                      size={32}
                      className="mx-auto mb-3 text-border"
                    />
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

export default LandingShowcase;
