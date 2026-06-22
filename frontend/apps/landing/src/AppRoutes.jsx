import { Route, Routes } from "react-router-dom";

import RouteSeo from "@smarterp/ui/seo/RouteSeo";
import RequireNoSession from "@smarterp/auth/guards/RequireNoSession";
import FallbackRoute from "@smarterp/auth/guards/FallbackRoute";
import PublicLayout from "./features/public/layout/PublicLayout";
import LandingPage from "./features/public/pages/LandingPage";
import LocalizedLandingPage from "./features/public/pages/LocalizedLandingPage";
import RegisterPage from "./features/public/pages/RegisterPage";
import LoginPage from "./features/public/pages/LoginPage";

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
        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </>
  );
}

export default AppRoutes;
