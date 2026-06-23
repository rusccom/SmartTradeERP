import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";

import I18nProvider from "@smarterp/i18n/I18nProvider";
import AppRoutes from "./AppRoutes";
import "./tailwind.css";
import "../../../packages/ui/styles/tokens.css";
import "../../../packages/ui/styles/global.css";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <I18nProvider>
        <AppRoutes />
      </I18nProvider>
    </BrowserRouter>
  </StrictMode>,
);
