import { lazy, Suspense } from "react";
import { useLocation } from "react-router-dom";
import AppRoutes from "./AppRoutes";

const ConsoleRoutes = lazy(() => import("../../console/src/AppRoutes"));

export default function WorkspaceRoutes() {
  const { pathname } = useLocation();
  if (pathname === "/admin" || pathname.startsWith("/admin/")) {
    return <Suspense fallback={null}><ConsoleRoutes /></Suspense>;
  }
  return <AppRoutes />;
}
