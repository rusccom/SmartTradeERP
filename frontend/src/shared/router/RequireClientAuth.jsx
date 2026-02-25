import { Navigate } from "react-router-dom";

import { hasClientSession } from "../auth/session";

function RequireClientAuth({ children }) {
  if (!hasClientSession()) {
    return <Navigate to="/login" replace />;
  }
  return children;
}

export default RequireClientAuth;
