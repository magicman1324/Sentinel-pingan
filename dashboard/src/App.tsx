import { Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import Rules from "./pages/Rules";
import Alerts from "./pages/Alerts";
import Topology from "./pages/Topology";

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<Dashboard />} />
        <Route path="/rules" element={<Rules />} />
        <Route path="/alerts" element={<Alerts />} />
        <Route path="/topology" element={<Topology />} />
      </Route>
    </Routes>
  );
}
