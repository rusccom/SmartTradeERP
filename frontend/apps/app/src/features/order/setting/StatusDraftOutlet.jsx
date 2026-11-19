import { useState } from "react";
import { Outlet } from "react-router-dom";

export default function StatusDraftOutlet() {
  const statusDraft = useState([]);
  return <Outlet context={statusDraft} />;
}
