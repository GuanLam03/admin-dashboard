import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Select from "react-select"; // ✅ fancy multi-select
import api from "../../api/axios";
import { useTranslation } from "react-i18next";

function ScheduleAddPage() {
  const {t} = useTranslation();
  const [formData, setFormData] = useState({
    title: "",
    recurrence: "",
    start_at: "",
    end_at: "",
    status: "active",
  });

  const [roles, setRoles] = useState([]); // list of roles from backend
  const [selectedNotifyRoles, setSelectedNotifyRoles] = useState([]); // array of string IDs

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  // Fetch roles from backend
  useEffect(() => {
    const fetchRoles = async () => {
      try {
        const res = await api.get("/roles"); // adjust endpoint
        // Expect: res.data.message = [{ id: 1, name: "Developer" }, ...]
        setRoles(res.data.message || []);
      } catch (err) {
        console.error("Failed to fetch roles", err);
      }
    };
    fetchRoles();
  }, []);

  // Handle normal inputs
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  // Submit form
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      await api.post("/add-schedules", {
        ...formData,
        // send as array of IDs (strings is fine; backend can cast to int if needed)
        notify_roles: selectedNotifyRoles,
      });

      setSuccess("Schedule added successfully!");
      setFormData({
        title: "",
        recurrence: "",
        start_at: "",
        end_at: "",
        status: "active",
      });
      setSelectedNotifyRoles([]);
      // navigate("/schedules"); // optional redirect
    } catch (err) {
      console.error(err);
      setError(err?.response?.data?.error || "Failed to add schedule. Please try again.");
    }
  };

  // Build react-select options and value
  const roleOptions = roles.map((role) => ({
    value: role.id,
    label: role.name || role.Name, // handle either key casing
  }));

  const selectedRoleOptions = roleOptions.filter((opt) =>
    selectedNotifyRoles.includes(opt.value)
  );

  return (
    <div>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}
      {success && (
        <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">
          {success}
        </div>
      )}

      <h2 className="text-xl font-bold mb-4">{t("scheduleManagement.addPage.title")}</h2>

      <form onSubmit={handleSubmit} className="bg-white p-4 rounded shadow-sm w-full">
        <h4 className="text-lg font-semibold mb-4">{t("scheduleManagement.sectionTitle")}</h4>
        <table className="w-full border-collapse">
          <tbody>
            <tr>
              <th className="text-left p-2 border">{t("scheduleManagement.fields.title")}</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="title"
                  value={formData.title}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  placeholder={t("scheduleManagement.addPage.form.scheduleTitlePlaceholder")}
                  required
                />
              </td>
            </tr>

            <tr>
              <th className="text-left p-2 border">{t("scheduleManagement.fields.recurrence")}</th>
              <td className="p-2 border">
                <select
                  name="recurrence"
                  value={formData.recurrence}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                >
                  <option value="">None</option>
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </td>
            </tr>

            <tr>
              <th className="text-left p-2 border">{t("scheduleManagement.fields.startAt")}</th>
              <td className="p-2 border">
                <input
                  type="datetime-local"
                  name="start_at"
                  value={formData.start_at}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>

            <tr>
              <th className="text-left p-2 border">{t("scheduleManagement.fields.endAt")}</th>
              <td className="p-2 border">
                <input
                  type="datetime-local"
                  name="end_at"
                  value={formData.end_at}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>

            <tr>
              <th className="text-left p-2 border">{t("scheduleManagement.fields.notifyRoles")}</th>
              <td className="p-2 border">
                <Select
                  isMulti
                  options={roleOptions}
                  value={selectedRoleOptions}
                  onChange={(selected) => {
                    setSelectedNotifyRoles((selected || []).map((s) => s.value));
                  }}
                  className="w-full"
                  classNamePrefix="react-select"
                  placeholder={t("scheduleManagement.addPage.form.scheduleNotifyRolesPlaceholder")}
                  isClearable
                />
                <small className="text-gray-500 block mt-1">
                  You can select multiple roles, or leave it blank to notify no one.
                </small>
              </td>
            </tr>
          </tbody>
        </table>

        <div className="mt-4 flex justify-end">
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded"
          >
            Save Schedule
          </button>
        </div>
      </form>
    </div>
  );
}

export default ScheduleAddPage;
