import { useState } from "react";
import { Server, Database, Router, Search } from "lucide-react";

const mockNodes = [
  { id: "prod-01", label: "prod-01", cpu: 32, mem: 55, status: "ok" },
  { id: "prod-02", label: "prod-02", cpu: 78, mem: 62, status: "warn" },
  { id: "prod-03", label: "prod-03", cpu: 45, mem: 38, status: "ok" },
  { id: "prod-04", label: "prod-04", cpu: 94, mem: 88, status: "critical" },
  { id: "prod-05", label: "prod-05", cpu: 22, mem: 41, status: "ok" },
  { id: "db-master", label: "TiDB-M", cpu: 51, mem: 70, status: "ok", icon: Database },
  { id: "db-slave-1", label: "TiDB-S1", cpu: 48, mem: 65, status: "ok", icon: Database },
  { id: "kafka-1", label: "Kafka-1", cpu: 35, mem: 42, status: "ok", icon: Router },
];

export default function Topology() {
  const [selected, setSelected] = useState<string | null>(null);

  const statusColor = (s: string) =>
    s === "critical" ? "border-red-500 shadow-red-500/20" :
    s === "warn" ? "border-yellow-500 shadow-yellow-500/20" : "border-green-500 shadow-green-500/20";

  const node = mockNodes.find(n => n.id === selected);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">集群拓扑</h1>
        <div className="flex items-center gap-2 bg-gray-900 rounded-lg px-3 py-2 border border-gray-800">
          <Search size={14} className="text-gray-500" />
          <input className="bg-transparent text-sm outline-none w-48" placeholder="搜索节点..." />
        </div>
      </div>

      <div className="grid grid-cols-4 gap-4">
        {mockNodes.map(n => {
          const Icon = n.icon || Server;
          return (
            <button
              key={n.id}
              onClick={() => setSelected(selected === n.id ? null : n.id)}
              className={`bg-gray-900 rounded-xl border-2 p-4 text-left transition-all hover:scale-105 cursor-pointer
                ${statusColor(n.status)} ${selected === n.id ? "ring-2 ring-green-400" : ""}`}
            >
              <div className="flex items-center gap-3 mb-3">
                <Icon className={n.status === "critical" ? "text-red-400" : n.status === "warn" ? "text-yellow-400" : "text-green-400"} size={24} />
                <span className="font-medium">{n.label}</span>
              </div>
              <div className="space-y-1 text-xs text-gray-400">
                <div className="flex justify-between">
                  <span>CPU</span>
                  <span className={n.cpu > 90 ? "text-red-400" : n.cpu > 70 ? "text-yellow-400" : ""}>{n.cpu}%</span>
                </div>
                <div className="w-full bg-gray-800 rounded-full h-1.5">
                  <div className={`h-1.5 rounded-full transition-all ${n.cpu > 90 ? "bg-red-500" : n.cpu > 70 ? "bg-yellow-500" : "bg-green-500"}`}
                    style={{ width: `${n.cpu}%` }} />
                </div>
                <div className="flex justify-between mt-2">
                  <span>内存</span>
                  <span className={n.mem > 85 ? "text-red-400" : n.mem > 70 ? "text-yellow-400" : ""}>{n.mem}%</span>
                </div>
                <div className="w-full bg-gray-800 rounded-full h-1.5">
                  <div className={`h-1.5 rounded-full transition-all ${n.mem > 85 ? "bg-red-500" : n.mem > 70 ? "bg-yellow-500" : "bg-green-500"}`}
                    style={{ width: `${n.mem}%` }} />
                </div>
              </div>
            </button>
          );
        })}
      </div>

      {node && (
        <div className="mt-6 bg-gray-900 rounded-xl border border-gray-800 p-4">
          <h3 className="font-bold mb-3">{node.label} — 详细信息</h3>
          <div className="grid grid-cols-3 gap-4 text-sm">
            <div><span className="text-gray-500">CPU</span><div className="text-xl font-bold">{node.cpu}%</div></div>
            <div><span className="text-gray-500">内存</span><div className="text-xl font-bold">{node.mem}%</div></div>
            <div><span className="text-gray-500">状态</span><div className="text-xl font-bold capitalize">{node.status}</div></div>
          </div>
        </div>
      )}
    </div>
  );
}
