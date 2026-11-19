import { useState } from "react";
import { Outlet } from "react-router-dom";
import "./orders.css";

export default function OrdersLayout() {
  const statusDraft = useState([]);
  return (
    <div className="orders-workspace" lang="uk">
      <Outlet context={statusDraft} />
    </div>
  );
}
