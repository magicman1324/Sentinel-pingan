const BASE = "/api/v1";

async function request<T>(path: string, opts?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...opts,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export interface Rule {
  id: number;
  name: string;
  description: string;
  rule_type: "atomic" | "composite";
  metric: string;
  operator: string;
  threshold: number;
  duration_sec: number;
  severity: string;
  expression?: string;
  enabled: boolean;
}

export interface Alert {
  id: number;
  rule_id: number;
  hostname: string;
  severity: string;
  metric: string;
  value: number;
  threshold: number;
  message: string;
  status: string;
  created_at: string;
}

export interface AlertPage {
  data: Alert[];
  total: number;
  page: number;
}

export interface Health {
  status: string;
  tidb: string;
  redis: string;
}

export const api = {
  getRules:      ()            => request<Rule[]>("/rules"),
  createRule:    (r: Rule)     => request<Rule>("/rules", { method: "POST", body: JSON.stringify(r) }),
  updateRule:    (id: number, r: Rule) => request<Rule>(`/rules/${id}`, { method: "PUT", body: JSON.stringify(r) }),
  deleteRule:    (id: number)  => request<void>(`/rules/${id}`, { method: "DELETE" }),

  getAlerts:     (page = 1)    => request<AlertPage>(`/alerts?page=${page}&size=20`),
  resolveAlert:  (id: number)  => request<void>(`/alerts/${id}/resolve`, { method: "POST" }),

  getHealth:     ()            => request<Health>("/health"),
};
