import { Navigate } from "react-router-dom";

import { hasAdminSession } from "@smarterp/auth/session";

function RequireAdminAuth({ children }) {
  if (!hasAdminSession()) {
    return <Navigate to="/admin" replace />;
  }
  return children;
}

export default RequireAdminAuth;
