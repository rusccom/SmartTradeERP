import { Link, Outlet, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";

function PublicLayout() {
  const [scrolled, setScrolled] = useState(false);
  const location = useLocation();
  const isLanding = location.pathname === "/";

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <div className="landing-zone flex flex-col bg-bg-primary">
      <header
        className={`sticky top-0 z-50 flex items-center justify-between
          px-4 sm:px-6 py-3 sm:py-4 transition-all duration-300 border-b
          ${scrolled
            ? "bg-bg-primary/90 backdrop-blur-xl border-border"
            : "bg-transparent border-transparent"
          }`}
      >
        <Link to="/" className="flex items-center gap-2.5 no-underline">
          <span className="w-8 h-8 rounded-lg bg-gradient-to-br from-accent to-cyan" />
          <span className="text-lg font-bold text-text-primary">
            SmartTrade
          </span>
        </Link>

        {isLanding && (
          <nav className="hidden md:flex items-center gap-6">
            {["Features", "How it works", "Pricing"].map((item) => (
              <a
                key={item}
                href={`#${item.toLowerCase().replace(/\s+/g, "-")}`}
                className="text-sm text-text-secondary hover:text-text-primary
                  transition-colors no-underline"
              >
                {item}
              </a>
            ))}
          </nav>
        )}

        <div className="flex items-center gap-3">
          <Link
            to="/login"
            className="text-sm font-semibold text-text-secondary
              hover:text-text-primary transition-colors no-underline"
          >
            Sign in
          </Link>
          <Link
            to="/register"
            className="text-sm font-semibold px-4 py-2 rounded-lg
              bg-accent hover:bg-accent-hover text-white
              transition-colors no-underline"
          >
            Get started
          </Link>
        </div>
      </header>

      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}

export default PublicLayout;
