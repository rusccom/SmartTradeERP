import { Route, Routes } from "react-router-dom";

import RouteSeo from "../seo/RouteSeo";
import LandingPage from "../../features/public/pages/LandingPage";
import RegisterPage from "../../features/public/pages/RegisterPage";
import LoginPage from "../../features/public/pages/LoginPage";
import PublicLayout from "../../features/public/layout/PublicLayout";
import AdminLoginPage from "../../features/admin/pages/AdminLoginPage";
import AdminDashboardPage from "../../features/admin/pages/AdminDashboardPage";
import AdminTenantsPage from "../../features/admin/pages/AdminTenantsPage";
import AdminLayout from "../../features/admin/layout/AdminLayout";
import DashboardPage from "../../features/dashboard/pages/DashboardPage";
import RegistryPage from "../../features/dashboard/pages/RegistryPage";
import ClientLayout from "../../features/dashboard/layout/ClientLayout";
import RequireAdminAuth from "./RequireAdminAuth";
import RequireClientAuth from "./RequireClientAuth";
import RequireNoSession from "./RequireNoSession";
import FallbackRoute from "./FallbackRoute";

function AppRoutes() {
  return (
    <>
      <RouteSeo />
      <Routes>
        <Route element={<RequireNoSession><PublicLayout /></RequireNoSession>}>
          <Route path="/" element={<LandingPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
        </Route>
        <Route element={<RequireNoSession><AdminLayout /></RequireNoSession>}>
          <Route path="/admin" element={<AdminLoginPage />} />
        </Route>
        <Route element={<RequireAdminAuth><AdminLayout /></RequireAdminAuth>}>
          <Route path="/admin/dashboard" element={<AdminDashboardPage />} />
          <Route path="/admin/tenants" element={<AdminTenantsPage />} />
        </Route>
        <Route element={<RequireClientAuth><ClientLayout /></RequireClientAuth>}>
          <Route path="/dashboard" element={<DashboardPage />} />
          {/* Catalog */}
          <Route path="/dashboard/products" element={<RegistryPage pageKey="products" />} />
          <Route path="/dashboard/groups" element={<RegistryPage pageKey="groups" />} />
          <Route path="/dashboard/customers" element={<RegistryPage pageKey="customers" />} />
          {/* Documents */}
          <Route path="/dashboard/docs/income" element={<RegistryPage pageKey="income" />} />
          <Route path="/dashboard/docs/expense" element={<RegistryPage pageKey="expense" />} />
          <Route path="/dashboard/docs/transfer" element={<RegistryPage pageKey="transfer" />} />
          <Route path="/dashboard/docs/inventory" element={<RegistryPage pageKey="inventory" />} />
          <Route path="/dashboard/docs/receipt" element={<RegistryPage pageKey="receipt" />} />
          <Route path="/dashboard/docs/writeoff" element={<RegistryPage pageKey="writeoff" />} />
        </Route>
        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </>
  );
}

export default AppRoutes;
