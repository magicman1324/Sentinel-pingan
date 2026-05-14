import { useEffect, useState } from "react";
import { CheckCircle, Filter } from "lucide-react";
import { api, Alert } from "../api/client";

export default function Alerts() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState("all");

  useEffect(() => {
    api.getAlerts(page).then(res => {
      setAlerts(res.data);
      setTotal(res.total);
    });
  }, [page]);

  const filtered = filter === "all" ? alerts : alerts.filter(a => a.severity === filter);

  const resolve = async (id: number) => {
    await api.resolveAlert(id);
    api.getAlerts(page).then(res => setAlerts(res.data));
  };

  const severityClass = (s: string) =>
    s === "critical" ? "border-l-red-500" :
    s === "warning" ? "border-l-yellow-500" : "border-l-blue-500";

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">告警历史</h1>
        <div className="flex items-center gap-3">
          <Filter size={14} className="text-gray-500" />
          {["all", "critical", "warning", "info"].map(f => (
            <button key={f}
              onClick={() => setFilter(f)}
              className={`text-xs px-3 py-1 rounded-full transition-colors ${
                filter === f ? "bg-gray-700 text-white" : "text-gray-500 hover:text-gray-300"
              }`}
            >{f === "all" ? "全部" : f}</button>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        {filtered.map(a => (
          <div key={a.id}
            className={`bg-gray-900 rounded-xl border border-gray-800 border-l-4 ${severityClass(a.severity)} p-4 flex items-center justify-between`}
          >
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-1">
                <span className={`px-2 py-0.5 rounded text-xs font-medium ${
                  a.severity === "critical" ? "bg-red-500/20 text-red-400" :
                  a.severity === "warning" ? "bg-yellow-500/20 text-yellow-400" : "bg-blue-500/20 text-blue-400"
                }`}>{a.severity}</span>
                <span className="font-medium">{a.hostname}</span>
                <span className="text-gray-500 text-xs">{a.metric}</span>
                {a.status === "firing" && (
                  <span className="bg-red-600 text-white text-xs px-2 py-0.5 rounded-full animate-pulse">FIRING</span>
                )}
              </div>
              <p className="text-sm text-gray-400">{a.message}</p>
              <div className="flex gap-4 mt-2 text-xs text-gray-600">
                <span>当前值: <code className="text-gray-400">{a.value.toFixed(2)}</code></span>
                <span>阈值: <code className="text-gray-400">{a.threshold.toFixed(2)}</code></span>
                <span>{new Date(a.created_at).toLocaleString()}</span>
              </div>
            </div>
            {a.status === "firing" && (
              <button onClick={() => resolve(a.id)}
                className="flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg bg-green-600 hover:bg-green-700 transition-colors ml-4">
                <CheckCircle size={14} /> 确认处理
              </button>
            )}
          </div>
        ))}
        {filtered.length === 0 && (
          <div className="text-center text-gray-600 py-12">无告警记录</div>
        )}
      </div>

      {total > 20 && (
        <div className="flex justify-center gap-4 mt-6">
          <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}
            className="px-4 py-2 text-sm rounded-lg border border-gray-700 disabled:opacity-30 hover:bg-gray-800">上一页</button>
          <span className="py-2 text-sm text-gray-400">第 {page} 页</span>
          <button disabled={page * 20 >= total} onClick={() => setPage(p => p + 1)}
            className="px-4 py-2 text-sm rounded-lg border border-gray-700 disabled:opacity-30 hover:bg-gray-800">下一页</button>
        </div>
      )}
    </div>
  );
}
