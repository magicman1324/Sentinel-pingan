import { useEffect, useState } from "react";
import { Plus, Trash2, ToggleLeft, ToggleRight } from "lucide-react";
import { api, Rule } from "../api/client";

const emptyRule: Rule = {
  id: 0, name: "", description: "", rule_type: "atomic",
  metric: "cpu.percent_used", operator: ">", threshold: 90,
  duration_sec: 0, severity: "warning", enabled: true,
};

export default function Rules() {
  const [rules, setRules] = useState<Rule[]>([]);
  const [editing, setEditing] = useState<Rule | null>(null);
  const [showForm, setShowForm] = useState(false);

  const load = () => api.getRules().then(setRules);

  useEffect(() => { load(); }, []);

  const save = async () => {
    if (!editing) return;
    if (editing.id) await api.updateRule(editing.id, editing);
    else await api.createRule(editing);
    setShowForm(false); setEditing(null); load();
  };

  const del = async (id: number) => {
    await api.deleteRule(id);
    load();
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">告警规则</h1>
        <button
          onClick={() => { setEditing({ ...emptyRule }); setShowForm(true); }}
          className="flex items-center gap-2 bg-green-600 hover:bg-green-700 px-4 py-2 rounded-lg text-sm transition-colors"
        >
          <Plus size={16} /> 新建规则
        </button>
      </div>

      {/* Rule Form Modal */}
      {showForm && editing && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl p-6 w-full max-w-lg border border-gray-700 space-y-4">
            <h2 className="text-lg font-bold">{editing.id ? "编辑规则" : "新建规则"}</h2>
            <input className="w-full bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700" placeholder="规则名称"
              value={editing.name} onChange={e => setEditing({ ...editing, name: e.target.value })} />
            <div className="grid grid-cols-2 gap-3">
              <select className="bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700"
                value={editing.rule_type} onChange={e => setEditing({ ...editing, rule_type: e.target.value as "atomic" | "composite" })}>
                <option value="atomic">原子规则</option>
                <option value="composite">复合规则</option>
              </select>
              <select className="bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700"
                value={editing.severity} onChange={e => setEditing({ ...editing, severity: e.target.value })}>
                <option value="critical">critical</option>
                <option value="warning">warning</option>
                <option value="info">info</option>
              </select>
            </div>
            <input className="w-full bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700" placeholder="指标 e.g. cpu.percent_used"
              value={editing.metric} onChange={e => setEditing({ ...editing, metric: e.target.value })} />
            <div className="grid grid-cols-3 gap-3">
              <select className="bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700"
                value={editing.operator} onChange={e => setEditing({ ...editing, operator: e.target.value })}>
                {[">", ">=", "<", "<=", "==", "!="].map(op => <option key={op} value={op}>{op}</option>)}
              </select>
              <input type="number" className="bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700" placeholder="阈值"
                value={editing.threshold} onChange={e => setEditing({ ...editing, threshold: +e.target.value })} />
              <input type="number" className="bg-gray-800 rounded px-3 py-2 text-sm border border-gray-700" placeholder="持续(秒)"
                value={editing.duration_sec} onChange={e => setEditing({ ...editing, duration_sec: +e.target.value })} />
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <button onClick={() => { setShowForm(false); setEditing(null); }}
                className="px-4 py-2 text-sm rounded-lg border border-gray-700 hover:bg-gray-800">取消</button>
              <button onClick={save} className="px-4 py-2 text-sm rounded-lg bg-green-600 hover:bg-green-700">保存</button>
            </div>
          </div>
        </div>
      )}

      {/* Rule List Table */}
      <div className="bg-gray-900 rounded-xl border border-gray-800 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-800 text-gray-400">
            <tr>
              <th className="text-left p-3">名称</th>
              <th className="text-left p-3">指标</th>
              <th className="text-left p-3">条件</th>
              <th className="text-left p-3">等级</th>
              <th className="text-left p-3">状态</th>
              <th className="p-3 w-20"></th>
            </tr>
          </thead>
          <tbody>
            {rules.map(r => (
              <tr key={r.id} className="border-t border-gray-800 hover:bg-gray-800/50 cursor-pointer"
                onClick={() => { setEditing({ ...r }); setShowForm(true); }}>
                <td className="p-3 font-medium">{r.name}</td>
                <td className="p-3 text-gray-400 font-mono text-xs">{r.metric}</td>
                <td className="p-3 font-mono text-xs">
                  {r.operator} {r.threshold}{r.duration_sec > 0 && ` (${r.duration_sec}s)`}
                </td>
                <td className="p-3">
                  <span className={`px-2 py-0.5 rounded text-xs font-medium ${
                    r.severity === "critical" ? "bg-red-500/20 text-red-400" :
                    r.severity === "warning" ? "bg-yellow-500/20 text-yellow-400" : "bg-blue-500/20 text-blue-400"
                  }`}>{r.severity}</span>
                </td>
                <td className="p-3">
                  {r.enabled ? <ToggleRight className="text-green-400" size={18} /> : <ToggleLeft className="text-gray-600" size={18} />}
                </td>
                <td className="p-3">
                  <button onClick={e => { e.stopPropagation(); del(r.id); }}
                    className="p-1 hover:bg-red-500/20 rounded text-gray-500 hover:text-red-400">
                    <Trash2 size={16} />
                  </button>
                </td>
              </tr>
            ))}
            {rules.length === 0 && (
              <tr><td colSpan={6} className="p-8 text-center text-gray-600">暂无规则，点击右上角新建</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
