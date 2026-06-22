import { Route, Routes } from "react-router-dom";

import RouteSeo from "@smarterp/ui/seo/RouteSeo";
import RequireNoSession from "@smarterp/auth/guards/RequireNoSession";
import RequireAdminAuth from "@smarterp/auth/guards/RequireAdminAuth";
import FallbackRoute from "@smarterp/auth/guards/FallbackRoute";
import AdminLayout from "./features/admin/layout/AdminLayout";
import AdminLoginPage from "./features/admin/pages/AdminLoginPage";
import AdminDashboardPage from "./features/admin/pages/AdminDashboardPage";
import AdminTenantsPage from "./features/admin/pages/AdminTenantsPage";

function AppRoutes() {
  return (
    <>
      <RouteSeo />
      <Routes>
        <Route element={<RequireNoSession><AdminLayout /></RequireNoSession>}>
          <Route path="/admin" element={<AdminLoginPage />} />
        </Route>
        <Route element={<RequireAdminAuth><AdminLayout /></RequireAdminAuth>}>
          <Route path="/admin/dashboard" element={<AdminDashboardPage />} />
          <Route path="/admin/tenants" element={<AdminTenantsPage />} />
        </Route>
        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </>
  );
}

export default AppRoutes;
