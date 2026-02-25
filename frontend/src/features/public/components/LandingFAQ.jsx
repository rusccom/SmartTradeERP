import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { ChevronDown } from "lucide-react";

const ITEMS = [
  {
    q: "Can I migrate data from my current system?",
    a: "Yes. You can import products, warehouses, and opening balances via CSV or use our API for automated migration from any existing system.",
  },
  {
    q: "Is there an API for integrations?",
    a: "The Pro and Enterprise plans include full REST API access. You can integrate SmartTrade with your POS, e-commerce platform, or accounting software.",
  },
  {
    q: "How is my data secured?",
    a: "Each tenant gets isolated data storage with encryption at rest and in transit. We run on enterprise-grade cloud infrastructure with automated backups.",
  },
  {
    q: "What happens when I exceed the free plan limits?",
    a: "You'll be notified before reaching limits. Upgrading is instant — no data migration needed, just unlock more capacity.",
  },
  {
    q: "Can multiple team members work simultaneously?",
    a: "Absolutely. Role-based access lets owners, managers, and cashiers work in parallel with appropriate permissions for each role.",
  },
];

function FAQItem({ item }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="border-b border-border last:border-b-0">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center justify-between py-5
          text-left bg-transparent border-none cursor-pointer
          text-text-primary font-semibold text-base
          hover:text-accent transition-colors"
      >
        {item.q}
        <ChevronDown
          size={18}
          className={`text-text-muted flex-shrink-0 ml-4
            transition-transform duration-200
            ${open ? "rotate-180" : ""}`}
        />
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
            <p
              className="pb-5 text-sm text-text-secondary
                leading-relaxed"
            >
              {item.a}
            </p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function LandingFAQ() {
  return (
    <section className="py-24 px-6">
      <div className="max-w-2xl mx-auto">
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
            Frequently asked questions
          </h2>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="rounded-2xl border border-border bg-bg-card
            px-6"
        >
          {ITEMS.map((item) => (
            <FAQItem key={item.q} item={item} />
          ))}
        </motion.div>
      </div>
    </section>
  );
}

export default LandingFAQ;
