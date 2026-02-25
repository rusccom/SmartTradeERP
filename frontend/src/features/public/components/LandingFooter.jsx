import { Link } from "react-router-dom";

const LINKS = {
  Product: [
    { label: "Features", href: "#features" },
    { label: "Pricing", href: "#pricing" },
    { label: "How it works", href: "#how-it-works" },
  ],
  Company: [
    { label: "About", href: "#" },
    { label: "Blog", href: "#" },
    { label: "Careers", href: "#" },
  ],
  Legal: [
    { label: "Privacy", href: "#" },
    { label: "Terms", href: "#" },
  ],
};

function LandingFooter() {
  return (
    <footer className="border-t border-border py-12 px-6">
      <div
        className="max-w-5xl mx-auto grid grid-cols-2 md:grid-cols-4
          gap-8"
      >
        <div className="col-span-2 md:col-span-1">
          <Link to="/" className="flex items-center gap-2 no-underline mb-4">
            <span
              className="w-7 h-7 rounded-lg bg-gradient-to-br
                from-accent to-cyan"
            />
            <span className="text-base font-bold text-text-primary">
              SmartTrade
            </span>
          </Link>
          <p className="text-sm text-text-muted leading-relaxed">
            Cloud ERP for retail &amp; distribution. Track inventory, post
            documents, know your real profit.
          </p>
        </div>

        {Object.entries(LINKS).map(([group, items]) => (
          <div key={group}>
            <p
              className="text-xs font-bold text-text-muted uppercase
                tracking-widest mb-4"
            >
              {group}
            </p>
            <ul className="list-none p-0 m-0 flex flex-col gap-2.5">
              {items.map((link) => (
                <li key={link.label}>
                  <a
                    href={link.href}
                    className="text-sm text-text-secondary
                      hover:text-text-primary transition-colors
                      no-underline"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      <div
        className="max-w-5xl mx-auto mt-10 pt-6 border-t
          border-border-subtle flex flex-col md:flex-row
          items-center justify-between gap-4"
      >
        <p className="text-xs text-text-muted">
          &copy; {new Date().getFullYear()} SmartTrade ERP. All rights
          reserved.
        </p>
      </div>
    </footer>
  );
}

export default LandingFooter;
