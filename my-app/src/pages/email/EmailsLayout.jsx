// EmailsLayout.jsx
import { Outlet } from "react-router-dom";

export default function EmailsLayout() {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Emails</h2>
      <Outlet />
    </div>
  );
}
