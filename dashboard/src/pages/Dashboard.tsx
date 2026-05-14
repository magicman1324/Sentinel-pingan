import { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from "recharts";
import { Activity, Cpu, HardDrive, Wifi } from "lucide-react";
import { api, Health } from "../api/client";

const mockData = Array.from({ length: 24 }, (_, i) => ({
  time: `${String(i).padStart(2, "0")}:00`,
  cpu: 30 + Math.random() * 40,
  memory: 50 + Math.random() * 20,
  netIn: Math.random() * 800,
  netOut: Math.random() * 400,
}));

export default function Dashboard() {
  const [health, setHealth] = useState<Health | null>(null);

  useEffect(() => {
    api.getHealth().then(setHealth).catch(() => {});
  }, []);

  const cards = [
    { icon: Cpu, label: "CPU", color: "text-blue-400", bg: "bg-blue-500/10" },
    { icon: HardDrive, label: "内存", color: "text-emerald-400", bg: "bg-emerald-500/10" },
    { icon: Wifi, label: "网络", color: "text-purple-400", bg: "bg-purple-500/10" },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">监控仪表板</h1>
        <div className="flex gap-4 text-sm">
          <span className={`flex items-center gap-1 ${health?.tidb === "ok" ? "text-green-400" : "text-red-400"}`}>
            <Activity size={14} /> TiDB: {health?.tidb ?? "—"}
          </span>
          <span className={`flex items-center gap-1 ${health?.redis === "ok" ? "text-green-400" : "text-red-400"}`}>
            <Activity size={14} /> Redis: {health?.redis ?? "—"}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        {cards.map(({ icon: Icon, label, color, bg }) => (
          <div key={label} className={`${bg} rounded-xl p-4 border border-gray-800`}>
            <div className="flex items-center gap-2 mb-2">
              <Icon className={color} size={20} />
              <span className="text-gray-400 text-sm">{label} 使用率</span>
            </div>
            <div className="text-3xl font-bold">{(30 + Math.random() * 40).toFixed(1)}%</div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-gray-900 rounded-xl p-4 border border-gray-800">
          <h3 className="text-sm text-gray-400 mb-4">CPU 使用率 (24h)</h3>
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={mockData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
              <XAxis dataKey="time" stroke="#6b7280" fontSize={12} />
              <YAxis stroke="#6b7280" fontSize={12} unit="%" />
              <Tooltip contentStyle={{ background: "#111827", border: "1px solid #374151", borderRadius: "8px" }} />
              <Area type="monotone" dataKey="cpu" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.1} />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-gray-900 rounded-xl p-4 border border-gray-800">
          <h3 className="text-sm text-gray-400 mb-4">内存使用率 (24h)</h3>
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={mockData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
              <XAxis dataKey="time" stroke="#6b7280" fontSize={12} />
              <YAxis stroke="#6b7280" fontSize={12} unit="%" />
              <Tooltip contentStyle={{ background: "#111827", border: "1px solid #374151", borderRadius: "8px" }} />
              <Area type="monotone" dataKey="memory" stroke="#10b981" fill="#10b981" fillOpacity={0.1} />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-gray-900 rounded-xl p-4 border border-gray-800 col-span-2">
          <h3 className="text-sm text-gray-400 mb-4">网络吞吐 (MB/s)</h3>
          <ResponsiveContainer width="100%" height={200}>
            <LineChart data={mockData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
              <XAxis dataKey="time" stroke="#6b7280" fontSize={12} />
              <YAxis stroke="#6b7280" fontSize={12} />
              <Tooltip contentStyle={{ background: "#111827", border: "1px solid #374151", borderRadius: "8px" }} />
              <Line type="monotone" dataKey="netIn" stroke="#a855f7" strokeWidth={2} name="入站" />
              <Line type="monotone" dataKey="netOut" stroke="#f59e0b" strokeWidth={2} name="出站" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}
