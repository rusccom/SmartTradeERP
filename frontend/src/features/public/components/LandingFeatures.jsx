import { motion } from "framer-motion";
import { Warehouse, FileText, TrendingUp, BarChart3 } from "lucide-react";

const FEATURES = [
  {
    icon: Warehouse,
    title: "Inventory & Warehousing",
    text: "Track stock across multiple warehouses in real time. Process receipts, write-offs, transfers, and full inventory counts — all from one interface.",
    span: "md:col-span-2 md:row-span-2",
  },
  {
    icon: FileText,
    title: "Documents & Workflow",
    text: "Five document types follow a clear lifecycle: Draft, Posted, Cancelled. Every posting updates your cost ledger automatically.",
    span: "",
  },
  {
    icon: TrendingUp,
    title: "Cost Accounting & Profit",
    text: "Automatic weighted average cost on every receipt. See true COGS and real profit margins based on actual numbers.",
    span: "",
  },
  {
    icon: BarChart3,
    title: "Reports & Analytics",
    text: "Profit reports by date range, stock summaries per warehouse, top-selling products ranked by margin, and complete movement history.",
    span: "md:col-span-2",
  },
];

function LandingFeatures() {
  return (
    <section id="features" className="py-24 px-6">
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-12"
        >
          <h2
            className="text-3xl md:text-5xl font-extrabold
              text-text-primary tracking-tight mb-4"
          >
            Everything your business needs
          </h2>
          <p className="text-text-secondary text-lg max-w-xl">
            From receiving stock to calculating profit — every step of your
            daily trade operations, covered.
          </p>
        </motion.div>

        <div className="grid md:grid-cols-4 gap-4">
          {FEATURES.map((f, i) => (
            <motion.article
              key={f.title}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className={`group relative rounded-2xl border border-border
                bg-bg-card p-6 transition-all duration-300
                hover:border-accent/40 hover:shadow-[0_0_30px_rgba(99,102,241,0.08)]
                ${f.span}`}
            >
              <div
                className="w-10 h-10 rounded-xl bg-accent/10
                  flex items-center justify-center mb-4
                  group-hover:bg-accent/20 transition-colors"
              >
                <f.icon size={20} className="text-accent" />
              </div>
              <h3
                className="text-base font-bold text-text-primary mb-2"
              >
                {f.title}
              </h3>
              <p className="text-sm text-text-secondary leading-relaxed">
                {f.text}
              </p>
            </motion.article>
          ))}
        </div>
      </div>
    </section>
  );
}

export default LandingFeatures;
