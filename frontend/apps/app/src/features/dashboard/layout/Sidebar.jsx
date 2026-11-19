import { Fragment } from "react";

import SidebarHeader from "./SidebarHeader";
import SidebarSection from "./SidebarSection";

function Sidebar({ open, sections, onClose }) {
  return (
    <>
      {open && <div className="sidebar-backdrop" onClick={onClose} />}
      <nav className={`sidebar-panel ${open ? "is-open" : ""}`}>
        <SidebarHeader onClose={onClose} />
        {sections.map((section, index) => (
          <Fragment key={section.label}>
            {index > 0 && <hr className="sidebar-divider" />}
            <SidebarSection section={section} onClose={onClose} />
          </Fragment>
        ))}
      </nav>
    </>
  );
}

export default Sidebar;
