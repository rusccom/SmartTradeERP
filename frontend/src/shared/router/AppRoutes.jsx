import { Navigate, Route, Routes } from "react-router-dom";

import RouteSeo from "../seo/RouteSeo";
import LandingPage from "../../features/public/pages/LandingPage";
import LocalizedLandingPage from "../../features/public/pages/LocalizedLandingPage";
import RegisterPage from "../../features/public/pages/RegisterPage";
import LoginPage from "../../features/public/pages/LoginPage";
import PublicLayout from "../../features/public/layout/PublicLayout";
import AdminLoginPage from "../../features/admin/pages/AdminLoginPage";
import AdminDashboardPage from "../../features/admin/pages/AdminDashboardPage";
import AdminTenantsPage from "../../features/admin/pages/AdminTenantsPage";
import AdminLayout from "../../features/admin/layout/AdminLayout";
import BundlesPage from "../../features/dashboard/pages/BundlesPage";
import CustomersPage from "../../features/dashboard/pages/CustomersPage";
import DashboardPage from "../../features/dashboard/pages/DashboardPage";
import DocumentsPage from "../../features/dashboard/pages/DocumentsPage";
import ProductsPage from "../../features/dashboard/pages/ProductsPage";
import ReportsPage from "../../features/dashboard/pages/ReportsPage";
import SettingsPage from "../../features/dashboard/pages/SettingsPage";
import WarehousesPage from "../../features/dashboard/pages/WarehousesPage";
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
          <Route path="/:locale" element={<LocalizedLandingPage />} />
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
          <Route path="/dashboard/products" element={<ProductsPage />} />
          <Route path="/dashboard/customers" element={<CustomersPage />} />
          <Route path="/dashboard/documents" element={<DocumentsPage />} />
          <Route path="/dashboard/warehouses" element={<WarehousesPage />} />
          <Route path="/dashboard/bundles" element={<BundlesPage />} />
          <Route path="/dashboard/reports" element={<ReportsPage />} />
          <Route path="/dashboard/settings" element={<SettingsPage />} />
          <Route path="/dashboard/groups" element={<Navigate to="/dashboard/bundles" replace />} />
          <Route path="/dashboard/docs/:type" element={<Navigate to="/dashboard/documents" replace />} />
        </Route>
        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </>
  );
}

export default AppRoutes;
