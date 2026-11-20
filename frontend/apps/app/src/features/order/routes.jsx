import { Route } from "react-router-dom";
import OrdersPage from "./list-order/OrdersPage";
import OrderEntryPage from "./create-order/OrderEntryPage";
import PolishOrderSettings from "../regions/pl/PolishOrderSettings";
import { orderPath, orderSettingsSegment } from "./menu";
import OrderSettingsPage from "./setting/OrderSettingsPage";
import StatusDraftOutlet from "./setting/StatusDraftOutlet";
import "./order.css";

const orderLayout = (
  <div className="orders-workspace">
    <StatusDraftOutlet />
  </div>
);

export const orderRoutes = (
  <Route path={orderPath} element={orderLayout}>
    <Route index element={<OrdersPage />} />
    <Route path="new" element={<OrderEntryPage />} />
    <Route path=":id" element={<OrderEntryPage />} />
    <Route path="settings/poland" element={<PolishOrderSettings />} />
    <Route path={orderSettingsSegment} element={<OrderSettingsPage />} />
  </Route>
);
