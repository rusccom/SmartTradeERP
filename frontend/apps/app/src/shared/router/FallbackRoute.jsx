import { Navigate, useLocation } from "react-router-dom";

import { resolveHomeRoute } from "@smarterp/auth/session";

function FallbackRoute() {
  const { pathname } = useLocation();
  return <Navigate to={resolveHomeRoute(pathname)} replace />;
}

export default FallbackRoute;
