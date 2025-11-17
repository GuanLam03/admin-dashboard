import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../../api/axios";
import { useTranslation } from "react-i18next";

function ScheduleEditPage() {
  const {t} = useTranslation();
  const { id } = useParams(); // get :id from URL
  const [formData, setFormData] = useState({
    title: "",
    recurrence: "daily",
    start_at: "",
    end_at: "",
    status: "active",
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // fetch existing schedule by id
  useEffect(() => {
    const fetchSchedule = async () => {
      try {
        const res = await api.get(`/edit-schedules/${id}`);
        const s = res.data.schedule;

        setFormData({
          title: s.title,
          recurrence: s.recurrence,
          start_at: s.start_at ? s.start_at.slice(0, 16) : "", // format for datetime-local
          end_at: s.end_at ? s.end_at.slice(0, 16) : "",
          status: s.status,
        });
      } catch (err) {
        console.error(err);
        setError("Failed to load schedule details.");
      } finally {
        setLoading(false);
      }
    };
    fetchSchedule();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      await api.post(`/edit-schedules/${id}`, formData);
      setSuccess("Schedule updated successfully!");

      // Redirect back to list after short delay
    //   setTimeout(() => {
    //     navigate("/schedules");
    //   }, 1000);
    } catch (err) {
      console.error(err);
      setError(t(err.response?.data?.error) || "Failed to update schedule. Please try again.");
    }
  };

  if (loading) {
    return <p>Loading schedule...</p>;
  }

  return (
    <div>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>
      )}
      {success && (
        <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>
      )}

      <h2 className="text-xl font-bold mb-4">{t("scheduleManagement.editPage.editSchedule")}</h2>

      <form onSubmit={handleSubmit} className="bg-white p-4 rounded shadow-sm w-full">
        <h4 className="text-lg font-semibold mb-4">{t("scheduleManagement.editPage.scheduleinfo")}</h4>
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
              <th className="text-left p-2 border">{t("scheduleManagement.fields.status")}</th>
              <td className="p-2 border">
                <select
                  name="status"
                  value={formData.status}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                >
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                  <option value="removed">Removed</option>
                </select>
              </td>
            </tr>
          </tbody>
        </table>

        <div className="mt-4 flex justify-end">
          <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">
            {t("common.buttons.update")}
          </button>
        </div>
      </form>
    </div>
  );
}

export default ScheduleEditPage;
