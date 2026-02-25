import { useEffect } from "react";
import { useLocation } from "react-router-dom";

const LANDING_META = {
  title: "SmartTrade ERP | Sales, inventory and documents in one ERP",
  description: "SmartTrade ERP helps teams manage sales, inventory, documents and reports from one workspace.",
  robots: "index,follow",
};

const PRIVATE_META = {
  title: "SmartTrade ERP | Secure Workspace",
  description: "SmartTrade ERP authentication and private workspace access.",
  robots: "noindex,nofollow",
};

function RouteSeo() {
  const { pathname } = useLocation();

  useEffect(() => {
    const meta = readMeta(pathname);
    document.title = meta.title;
    upsertMeta("description", meta.description);
    upsertMeta("robots", meta.robots);
    upsertMeta("googlebot", meta.robots);
    syncCanonical(pathname);
  }, [pathname]);

  return null;
}

function readMeta(pathname) {
  if (pathname === "/") {
    return LANDING_META;
  }
  return PRIVATE_META;
}

function upsertMeta(name, content) {
  const selector = `meta[name="${name}"]`;
  const element = document.head.querySelector(selector) ?? createMeta(name);
  element.setAttribute("content", content);
}

function createMeta(name) {
  const element = document.createElement("meta");
  element.setAttribute("name", name);
  document.head.appendChild(element);
  return element;
}

function syncCanonical(pathname) {
  const existing = document.head.querySelector('link[rel="canonical"]');
  if (pathname !== "/") {
    existing?.remove();
    return;
  }
  const canonical = existing ?? createCanonical();
  canonical.setAttribute("href", `${window.location.origin}/`);
}

function createCanonical() {
  const element = document.createElement("link");
  element.setAttribute("rel", "canonical");
  document.head.appendChild(element);
  return element;
}

export default RouteSeo;
