import { NavLink } from "react-router-dom";

function Sidebar({ open, sections, onClose }) {
  return (
    <>
      {open && (
        <div className="sidebar-backdrop" onClick={onClose} />
      )}
      <nav className={`sidebar-panel ${open ? "is-open" : ""}`}>
        <div className="sidebar-header">
          <span className="sidebar-title">Menu</span>
          <button
            className="sidebar-close"
            type="button"
            onClick={onClose}
            aria-label="Close menu"
          >
            &times;
          </button>
        </div>

        {sections.map((section, idx) => (
          <SidebarSection
            key={section.label}
            section={section}
            showDivider={idx > 0}
            onClose={onClose}
          />
        ))}
      </nav>
    </>
  );
}

function SidebarSection({ section, showDivider, onClose }) {
  return (
    <>
      {showDivider && <hr className="sidebar-divider" />}
      <p className="sidebar-section-label">{section.label}</p>
      <ul className="sidebar-list">
        {section.items.map((item) => (
          <li key={item.key}>
            <NavLink
              to={item.path}
              className={readLinkClass}
              onClick={onClose}
            >
              {item.title}
            </NavLink>
          </li>
        ))}
      </ul>
    </>
  );
}

function readLinkClass({ isActive }) {
  return isActive
    ? "sidebar-link is-active"
    : "sidebar-link";
}

export default Sidebar;
