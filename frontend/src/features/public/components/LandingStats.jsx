import { motion } from "framer-motion";

const STATS = [
  { value: "500+", label: "Active businesses" },
  { value: "99.9%", label: "Uptime SLA" },
  { value: "<50ms", label: "Avg response" },
  { value: "24/7", label: "Cloud access" },
];

/*
  ЛОГОТИПЫ КЛИЕНТОВ — замените заглушки ниже на реальные логотипы:
  <img src="/logos/client1.svg" alt="Client name" className="h-8 opacity-50" />
*/
const LOGOS = [
  "Partner 1",
  "Partner 2",
  "Partner 3",
  "Partner 4",
  "Partner 5",
];

function LandingStats() {
  return (
    <section className="py-16 px-6 border-y border-border">
      <div className="max-w-5xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
          {STATS.map((s, i) => (
            <motion.div
              key={s.label}
              initial={{ opacity: 0, y: 16 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="text-center"
            >
              <p
                className="text-3xl md:text-4xl font-extrabold
                  bg-gradient-to-r from-accent to-cyan bg-clip-text
                  text-transparent"
              >
                {s.value}
              </p>
              <p className="text-sm text-text-muted mt-1">{s.label}</p>
            </motion.div>
          ))}
        </div>

        <div className="flex flex-wrap items-center justify-center gap-8">
          {LOGOS.map((name) => (
            <span
              key={name}
              className="px-4 py-2 rounded-lg border border-border-subtle
                bg-bg-card/50 text-text-muted text-sm font-medium"
            >
              {name}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}

export default LandingStats;
