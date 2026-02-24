import { Route, Routes } from "react-router-dom";

import AppFrame from "../../app/AppFrame";
import LandingPage from "../../features/public/pages/LandingPage";
import RegisterPage from "../../features/public/pages/RegisterPage";
import LoginPage from "../../features/public/pages/LoginPage";
import AdminLoginPage from "../../features/admin/pages/AdminLoginPage";
import AdminDashboardPage from "../../features/admin/pages/AdminDashboardPage";
import AdminTenantsPage from "../../features/admin/pages/AdminTenantsPage";
import DashboardPage from "../../features/dashboard/pages/DashboardPage";
import ProductsPage from "../../features/dashboard/pages/ProductsPage";
import BundlesPage from "../../features/dashboard/pages/BundlesPage";
import WarehousesPage from "../../features/dashboard/pages/WarehousesPage";
import DocumentsPage from "../../features/dashboard/pages/DocumentsPage";
import DocumentCardPage from "../../features/dashboard/pages/DocumentCardPage";
import ReportsPage from "../../features/dashboard/pages/ReportsPage";
import SettingsPage from "../../features/dashboard/pages/SettingsPage";

function AppRoutes() {
  return (
    <Routes>
      <Route element={<AppFrame />}>
        <Route path="/" element={<LandingPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/admin" element={<AdminLoginPage />} />
        <Route path="/admin/dashboard" element={<AdminDashboardPage />} />
        <Route path="/admin/tenants" element={<AdminTenantsPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/dashboard/products" element={<ProductsPage />} />
        <Route path="/dashboard/bundles" element={<BundlesPage />} />
        <Route path="/dashboard/warehouses" element={<WarehousesPage />} />
        <Route path="/dashboard/documents" element={<DocumentsPage />} />
        <Route path="/dashboard/documents/:id" element={<DocumentCardPage />} />
        <Route path="/dashboard/reports" element={<ReportsPage />} />
        <Route path="/dashboard/settings" element={<SettingsPage />} />
      </Route>
    </Routes>
  );
}

export default AppRoutes;
