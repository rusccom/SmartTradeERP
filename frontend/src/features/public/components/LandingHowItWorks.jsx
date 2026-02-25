import { motion } from "framer-motion";
import { UserPlus, PackagePlus, Rocket } from "lucide-react";

const STEPS = [
  {
    icon: UserPlus,
    num: "01",
    title: "Create your workspace",
    text: "Register in 30 seconds. Your isolated tenant is provisioned instantly with role-based access for owners, managers, and cashiers.",
  },
  {
    icon: PackagePlus,
    num: "02",
    title: "Set up your catalog",
    text: "Add products, define variants and composite bundles, assign warehouses. Import existing data or start fresh.",
  },
  {
    icon: Rocket,
    num: "03",
    title: "Start trading",
    text: "Process sales, receipts, and transfers. Cost calculations, ledger postings, and profit tracking happen automatically.",
  },
];

function LandingHowItWorks() {
  return (
    <section id="how-it-works" className="py-24 px-6">
      <div className="max-w-3xl mx-auto">
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
            Up and running in minutes
          </h2>
          <p className="text-text-secondary text-lg">
            No complex setup, no consultants.
          </p>
        </motion.div>

        <div className="relative">
          <div
            className="absolute left-5 md:left-6 top-0 bottom-0 w-px
              bg-gradient-to-b from-accent via-cyan to-transparent"
          />

          <div className="flex flex-col gap-12">
            {STEPS.map((s, i) => (
              <motion.div
                key={s.num}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.15 }}
                className="relative flex gap-6 items-start"
              >
                <div
                  className="relative z-10 flex-shrink-0 w-10 h-10
                    md:w-12 md:h-12 rounded-xl bg-bg-card border
                    border-border flex items-center justify-center"
                >
                  <s.icon size={20} className="text-accent" />
                </div>

                <div className="pt-1">
                  <span
                    className="text-xs font-bold text-accent
                      tracking-widest"
                  >
                    STEP {s.num}
                  </span>
                  <h3
                    className="text-lg font-bold text-text-primary mt-1
                      mb-2"
                  >
                    {s.title}
                  </h3>
                  <p
                    className="text-sm text-text-secondary
                      leading-relaxed"
                  >
                    {s.text}
                  </p>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

export default LandingHowItWorks;
