import { Link } from "react-router-dom";

import { buildLandingPath } from "../../../shared/i18n/localeConfig";
import { useI18n } from "../../../shared/i18n/useI18n";

function LandingFooter() {
  const { locale, t } = useI18n();
  const landingPath = buildLandingPath(locale);
  const links = createLinks(t);
  return (
    <footer className="border-t border-border py-12 px-6">
      <div className="max-w-5xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-8">
        <div className="col-span-2 md:col-span-1">
          <Link to={landingPath} className="flex items-center gap-2 no-underline mb-4">
            <span className="w-7 h-7 rounded-lg bg-gradient-to-br from-accent to-cyan" />
            <span className="text-base font-bold text-text-primary">{t("brand.name")}</span>
          </Link>
          <p className="text-sm text-text-muted leading-relaxed">{t("public.footer.description")}</p>
        </div>
        {links.map((group) => (
          <div key={group.label}>
            <p className="text-xs font-bold text-text-muted uppercase tracking-widest mb-4">{group.label}</p>
            <ul className="list-none p-0 m-0 flex flex-col gap-2.5">
              {group.items.map((item) => (
                <li key={item.label}>
                  <a href={item.href} className="text-sm text-text-secondary hover:text-text-primary transition-colors no-underline">
                    {item.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
      <div className="max-w-5xl mx-auto mt-10 pt-6 border-t border-border-subtle flex flex-col md:flex-row items-center justify-between gap-4">
        <p className="text-xs text-text-muted">{t("public.footer.rights", { year: new Date().getFullYear() })}</p>
      </div>
    </footer>
  );
}

function createLinks(t) {
  return [
    {
      label: t("public.footer.product"),
      items: [
        { label: t("public.nav.features"), href: "#features" },
        { label: t("public.footer.pricing"), href: "#pricing" },
        { label: t("public.footer.howItWorks"), href: "#how-it-works" },
      ],
    },
    {
      label: t("public.footer.company"),
      items: [
        { label: t("public.footer.about"), href: "#" },
        { label: t("public.footer.blog"), href: "#" },
        { label: t("public.footer.careers"), href: "#" },
      ],
    },
    {
      label: t("public.footer.legal"),
      items: [
        { label: t("public.footer.privacy"), href: "#" },
        { label: t("public.footer.terms"), href: "#" },
      ],
    },
  ];
}

export default LandingFooter;
