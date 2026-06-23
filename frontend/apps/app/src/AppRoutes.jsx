import { Navigate, Route, Routes } from "react-router-dom";

import RouteSeo from "@smarterp/ui/seo/RouteSeo";
import RequireClientAuth from "@smarterp/auth/guards/RequireClientAuth";
import RequireNoSession from "@smarterp/auth/guards/RequireNoSession";
import FallbackRoute from "@smarterp/auth/guards/FallbackRoute";
import PublicLayout from "./features/public/layout/PublicLayout";
import LandingPage from "./features/public/pages/LandingPage";
import LocalizedLandingPage from "./features/public/pages/LocalizedLandingPage";
import LoginPage from "./features/public/pages/LoginPage";
import RegisterPage from "./features/public/pages/RegisterPage";
import ClientLayout from "./features/dashboard/layout/ClientLayout";
import DashboardPage from "./features/dashboard/pages/DashboardPage";
import CustomersPage from "./features/dashboard/pages/CustomersPage";
import DocumentsPage from "./features/dashboard/pages/DocumentsPage";
import ReportsPage from "./features/dashboard/pages/ReportsPage";
import SettingsPage from "./features/dashboard/pages/SettingsPage";
import WarehousesPage from "./features/dashboard/pages/WarehousesPage";
import ProductCatalogPage from "./features/products/pages/ProductCatalogPage";
import ProductsPage from "./features/products/pages/ProductsPage";
import BundlesPage from "./features/bundles/pages/BundlesPage";

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
        <Route element={<RequireClientAuth><ClientLayout /></RequireClientAuth>}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/dashboard/products" element={<ProductCatalogPage />}>
            <Route index element={<ProductsPage />} />
            <Route path="bundles" element={<BundlesPage />} />
          </Route>
          <Route path="/dashboard/customers" element={<CustomersPage />} />
          <Route path="/dashboard/documents" element={<DocumentsPage />} />
          <Route path="/dashboard/warehouses" element={<WarehousesPage />} />
          <Route path="/dashboard/bundles" element={<Navigate to="/dashboard/products/bundles" replace />} />
          <Route path="/dashboard/reports" element={<ReportsPage />} />
          <Route path="/dashboard/settings" element={<SettingsPage />} />
          <Route path="/dashboard/groups" element={<Navigate to="/dashboard/products/bundles" replace />} />
          <Route path="/dashboard/docs/:type" element={<Navigate to="/dashboard/documents" replace />} />
        </Route>
        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </>
  );
}

export default AppRoutes;
