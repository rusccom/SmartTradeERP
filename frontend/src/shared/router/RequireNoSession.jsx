import { Navigate, useLocation } from "react-router-dom";

import { hasAnySession, resolveHomeRoute } from "../auth/session";

function RequireNoSession({ children }) {
  const { pathname } = useLocation();
  if (!hasAnySession()) {
    return children;
  }
  return <Navigate to={resolveHomeRoute(pathname)} replace />;
}

export default RequireNoSession;
