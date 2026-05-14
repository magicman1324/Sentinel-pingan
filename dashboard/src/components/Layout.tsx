import { NavLink, Outlet } from "react-router-dom";
import { LayoutDashboard, ShieldCheck, BellRing, GitGraph, Moon, Sun } from "lucide-react";
import { useState } from "react";

const nav = [
  { to: "/",         icon: LayoutDashboard, label: "仪表板" },
  { to: "/rules",    icon: ShieldCheck,      label: "规则" },
  { to: "/alerts",   icon: BellRing,         label: "告警" },
  { to: "/topology", icon: GitGraph,         label: "拓扑" },
];

export default function Layout() {
  const [dark, setDark] = useState(true);
  const toggle = () => {
    document.documentElement.classList.toggle("dark");
    setDark(!dark);
  };

  return (
    <div className="flex h-screen">
      <aside className="w-56 bg-gray-900 border-r border-gray-800 flex flex-col">
        <div className="p-4 border-b border-gray-800 flex items-center gap-2">
          <ShieldCheck className="text-green-500" size={24} />
          <span className="font-bold text-lg">哨兵</span>
        </div>
        <nav className="flex-1 p-2 space-y-1">
          {nav.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded text-sm transition-colors ${
                  isActive ? "bg-gray-800 text-green-400" : "text-gray-400 hover:text-gray-200 hover:bg-gray-800/50"
                }`
              }
            >
              <Icon size={18} />
              {label}
            </NavLink>
          ))}
        </nav>
        <button onClick={toggle} className="m-3 p-2 rounded hover:bg-gray-800 text-gray-400">
          {dark ? <Sun size={18} /> : <Moon size={18} />}
        </button>
      </aside>
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
