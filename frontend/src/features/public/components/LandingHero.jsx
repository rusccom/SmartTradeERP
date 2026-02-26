import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { ArrowRight } from "lucide-react";

function LandingHero() {
  return (
    <section className="relative overflow-hidden pt-20 pb-24 px-6">
      <div className="absolute inset-0 pointer-events-none">
        <div
          className="absolute top-[-20%] left-1/2 -translate-x-1/2
            w-[800px] h-[600px] rounded-full blur-[120px]
            bg-[rgba(99,102,241,0.15)]"
        />
        <div
          className="absolute bottom-[-10%] right-[-10%]
            w-[400px] h-[400px] rounded-full blur-[100px]
            bg-[rgba(6,182,212,0.1)]"
        />
      </div>

      <div className="relative max-w-4xl mx-auto text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="inline-flex items-center gap-2 px-4 py-1.5 mb-6
            rounded-full border border-border bg-bg-card/50 text-sm
            text-text-secondary"
        >
          <span className="w-2 h-2 rounded-full bg-emerald-500
            animate-pulse" />
          Trusted by 500+ businesses
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="text-4xl sm:text-5xl md:text-7xl font-extrabold
            leading-[1.05] tracking-tight mb-6"
        >
          <span className="text-text-primary">Track every item.</span>
          <br />
          <span
            className="bg-gradient-to-r from-accent to-cyan
              bg-clip-text text-transparent"
          >
            Know your real profit.
          </span>
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="max-w-2xl mx-auto text-lg text-text-secondary
            leading-relaxed mb-10"
        >
          Manage inventory across multiple warehouses, process sales and
          receipts with automatic cost accounting, and see real profit
          margins — all from one workspace.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
          className="flex flex-col sm:flex-row items-center
            justify-center gap-4"
        >
          <Link
            to="/register"
            className="inline-flex items-center gap-2 px-7 py-3.5
              rounded-xl bg-accent hover:bg-accent-hover text-white
              font-semibold text-base transition-all no-underline
              shadow-[0_0_24px_rgba(99,102,241,0.3)]
              hover:shadow-[0_0_32px_rgba(99,102,241,0.5)]"
          >
            Start free — no card required
            <ArrowRight size={18} />
          </Link>
          <a
            href="#features"
            className="inline-flex items-center gap-2 px-7 py-3.5
              rounded-xl border border-border text-text-secondary
              hover:text-text-primary hover:border-text-muted
              font-semibold text-base transition-all no-underline"
          >
            See features
          </a>
        </motion.div>
      </div>

      {/*
        СКРИНШОТ ПРОДУКТА — замените содержимое div ниже на:
        <img src="/screenshots/dashboard.png" alt="SmartTrade Dashboard"
          className="w-full" />
      */}
      <motion.div
        initial={{ opacity: 0, y: 40 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.7, delay: 0.5 }}
        className="relative max-w-5xl mx-auto mt-16"
      >
        <div
          className="absolute inset-0 bg-gradient-to-t from-bg-primary
            via-transparent to-transparent z-10 pointer-events-none"
        />
        <div
          className="relative rounded-2xl border border-border bg-bg-card
            overflow-hidden shadow-2xl shadow-accent/5"
        >
          <div
            className="w-full aspect-[16/9] bg-gradient-to-br
              from-bg-secondary to-bg-card flex items-center
              justify-center text-text-muted"
          >
            <div className="text-center">
              <div
                className="w-16 h-16 rounded-2xl bg-bg-secondary
                  border border-border mx-auto mb-4 flex items-center
                  justify-center text-2xl"
              >
                📊
              </div>
              <p className="text-sm font-medium">Product Screenshot</p>
              <p className="text-xs mt-1">
                Place your dashboard image here
              </p>
            </div>
          </div>
        </div>
      </motion.div>
    </section>
  );
}

export default LandingHero;
