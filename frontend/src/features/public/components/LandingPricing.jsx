import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Check } from "lucide-react";

const PLANS = [
  {
    name: "Starter",
    price: "Free",
    period: "",
    desc: "For solo entrepreneurs getting started.",
    features: [
      "1 warehouse",
      "Up to 500 products",
      "All 5 document types",
      "Basic reports",
    ],
    cta: "Get started",
    highlighted: false,
  },
  {
    name: "Pro",
    price: "$29",
    period: "/mo",
    desc: "For growing businesses with multiple locations.",
    features: [
      "Unlimited warehouses",
      "Unlimited products",
      "Composite bundles",
      "Advanced analytics",
      "Role-based access",
      "Priority support",
    ],
    cta: "Start free trial",
    highlighted: true,
  },
  {
    name: "Enterprise",
    price: "Custom",
    period: "",
    desc: "For large-scale distribution operations.",
    features: [
      "Everything in Pro",
      "API access",
      "Dedicated account manager",
      "Custom integrations",
      "SLA guarantee",
    ],
    cta: "Contact sales",
    highlighted: false,
  },
];

function LandingPricing() {
  return (
    <section id="pricing" className="py-24 px-6">
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
            Simple, transparent pricing
          </h2>
          <p className="text-text-secondary text-lg">
            Start free. Upgrade when you grow.
          </p>
        </motion.div>

        <div className="grid md:grid-cols-3 gap-6">
          {PLANS.map((plan, i) => (
            <motion.div
              key={plan.name}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className={`rounded-2xl border p-6 flex flex-col
                ${plan.highlighted
                  ? "border-accent bg-accent/5 shadow-[0_0_40px_rgba(99,102,241,0.1)]"
                  : "border-border bg-bg-card"
                }`}
            >
              {plan.highlighted && (
                <span
                  className="inline-block self-start text-xs font-bold
                    text-accent bg-accent/10 px-3 py-1 rounded-full
                    mb-4"
                >
                  Most popular
                </span>
              )}
              <h3 className="text-lg font-bold text-text-primary">
                {plan.name}
              </h3>
              <div className="mt-3 mb-2">
                <span
                  className="text-4xl font-extrabold text-text-primary"
                >
                  {plan.price}
                </span>
                {plan.period && (
                  <span className="text-text-muted text-sm">
                    {plan.period}
                  </span>
                )}
              </div>
              <p className="text-sm text-text-secondary mb-6">
                {plan.desc}
              </p>

              <ul className="flex flex-col gap-3 mb-8 flex-1">
                {plan.features.map((f) => (
                  <li
                    key={f}
                    className="flex items-center gap-2 text-sm
                      text-text-secondary"
                  >
                    <Check size={16} className="text-accent flex-shrink-0" />
                    {f}
                  </li>
                ))}
              </ul>

              <Link
                to="/register"
                className={`w-full py-3 rounded-xl font-semibold
                  text-sm text-center no-underline transition-all
                  ${plan.highlighted
                    ? "bg-accent hover:bg-accent-hover text-white shadow-[0_0_20px_rgba(99,102,241,0.3)]"
                    : "border border-border text-text-secondary hover:text-text-primary hover:border-text-muted"
                  }`}
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

export default LandingPricing;
