import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../../api/axios";
import { useTranslation } from "react-i18next";

function GoogleDocumentEditPage() {
  const {t} = useTranslation();
  const { id } = useParams(); // get :id from URL
  const [formData, setFormData] = useState({
    name: "",
    original_link: "",
    status: "",
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // fetch existing document by id
  useEffect(() => {
    const fetchDocument = async () => {
      try {
        const res = await api.get(`/edit-google-documents/${id}`);
        setFormData({
          name: res.data.document.name,
          original_link: res.data.document.original_link,
          status: res.data.document.status,
        });
      } catch (err) {
        console.error(err);
        setError("Failed to load document details.");
      } finally {
        setLoading(false);
      }
    };
    fetchDocument();
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
      await api.post(`/edit-google-documents/${id}`, formData);
      setSuccess("Google Document updated successfully!");

      // Redirect back to list after short delay
      setTimeout(() => {
        navigate("/google-documents");
      }, 1000);
    } catch (err) {
      console.error(err);
      setError(
        err.response?.data?.error || "Failed to update Google Document. Please try again."
      );
    }
  };

  if (loading) {
    return <p>Loading document...</p>;
  }

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

      <h2 className="text-xl font-bold mb-4">{t("googleDocumentManagement.editPage.editGoogleDocument")}</h2>

      <form
        onSubmit={handleSubmit}
        className="bg-white p-4 rounded shadow-sm w-full"
      >
        <h4 className="text-lg font-semibold mb-4">{t("googleDocumentManagement.editPage.googleDocumentInfo")}</h4>
        <table className="w-full border-collapse">
          <tbody>
            <tr>
              <th className="text-left p-2 border">{t("googleDocumentManagement.fields.name")}</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">{t("googleDocumentManagement.fields.originalLink")}</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="original_link"
                  value={formData.original_link}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">{t("googleDocumentManagement.fields.status")}</th>
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
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded"
          >
            {t("common.buttons.update")}
          </button>
        </div>
      </form>
    </div>
  );
}

export default GoogleDocumentEditPage;
