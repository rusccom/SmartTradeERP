import { NavLink } from "react-router-dom";

export default function SidebarLinks({ items, onClose }) {
  return (
    <ul className="sidebar-list">
      {items.map((item) => (
        <li key={item.key}>
          <NavLink to={item.path} end={item.end} className={readLinkClass} onClick={onClose}>
            {item.title}
          </NavLink>
        </li>
      ))}
    </ul>
  );
}

function readLinkClass({ isActive }) {
  return isActive ? "sidebar-link is-active" : "sidebar-link";
}
