import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";

function GoogleDocumentAddPage() {
  const [formData, setFormData] = useState({
    name: "",
    original_link: "",
    status: "active",
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      await api.post("/add-google-documents", formData);
      setSuccess("Google Document added successfully!");
      setFormData({ name: "", original_link: "", status: "active" });

      // Redirect back to list after short delay
      setTimeout(() => {
        navigate("/google-documents");
      }, 1000);
    } catch (err) {
      console.error(err);
      setError(err.response.data.error || "Failed to add Google Document. Please try again.");
    }
  };

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

      <h2 className="text-xl font-bold mb-4">Add Google Document</h2>

      <form
        onSubmit={handleSubmit}
        className="bg-white p-4 rounded shadow-sm w-full"
      >
        <h4 className="text-lg font-semibold mb-4">Google Document Info</h4>
        <table className="w-full border-collapse">
          <tbody>
            <tr>
              <th className="text-left p-2 border">Name</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  placeholder="Enter document name"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Link</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="original_link"
                  value={formData.original_link}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  placeholder="Paste Google Doc link"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Status</th>
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
            Save Document
          </button>
        </div>
      </form>
    </div>
  );
}

export default GoogleDocumentAddPage;
