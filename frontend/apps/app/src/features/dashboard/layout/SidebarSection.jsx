import { ChevronDown } from "lucide-react";
import SidebarLinks from "./SidebarLinks";

export default function SidebarSection({ section, onClose }) {
  const links = <SidebarLinks items={section.items} onClose={onClose} />;
  if (section.collapsible) {
    return (
      <details className="sidebar-group" open>
        <summary className="sidebar-section-label sidebar-toggle">
          {section.label}<ChevronDown size={15} aria-hidden="true" />
        </summary>
        {links}
      </details>
    );
  }
  return <><p className="sidebar-section-label">{section.label}</p>{links}</>;
}
