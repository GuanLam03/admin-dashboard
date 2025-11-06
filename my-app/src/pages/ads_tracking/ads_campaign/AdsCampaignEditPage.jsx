import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../../../api/axios";

const eventOptions = [
  { value: "PURCHASE", label: "Purchase"},
  { value: "COMPLETE_REGISTRATION", label: "Registration"},
  { value: "FORM_SUBMIT", label: "Form Submit"}
]

function AdsCampaignEditPage() {
  const { id } = useParams(); // get :id from URL
  const [formData, setFormData] = useState({});
  const [postbacks, setPostbacks] = useState([]);

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const [statusOptions, setStatusOptions] = useState([]);


 
  useEffect(() => {
    const fetchAdsCampaign = async () => {
      try {
        const res = await api.get(`/edit-ads-campaign/${id}`);
        const result = res.data.ads_campaign;
        const statusList = res.data.status || [];
        const campaignPostbacks = res.data.ads_campaign_postbacks || [];

        setFormData({
            name: result.name || "",
            target_url: result.target_url || "",
            tracking_link: result.tracking_link || "",
            postback_link: result.postback_link || "",
            status: result.status || "",
        });

        setStatusOptions(statusList);
        setPostbacks(campaignPostbacks);

      } catch (err) {
        console.error(err);
        setError("Failed to load ads campaign details.");
      } finally {
        setLoading(false);
      }
    };
    fetchAdsCampaign();
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
      await api.post(`/edit-ads-campaign/${id}`, formData);
      setSuccess("Ads Campaign updated successfully!");

      // Redirect back to list after short delay
      setTimeout(() => {
        navigate("/ads-tracking/campaign");
      }, 1000);
    } catch (err) {
      console.error(err);
      setError(
        err.response?.data?.error || "Failed to update Ads Campaign. Please try again."
      );
    }
  };

  if (loading) {
    return <p>Loading Ads Campaign...</p>;
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

      <h2 className="text-xl font-bold mb-4">Edit Ads Campaign</h2>

      <form
        onSubmit={handleSubmit}
        className="bg-white p-4 rounded shadow-sm w-full"
      >
        <h4 className="text-lg font-semibold mb-4">Ads Campaign Info</h4>
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
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Target Url</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="target_url"
                  value={formData.target_url}
                  onChange={handleChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Tracking Link</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="tracking_link"
                  value={formData.tracking_link}
                  onChange={handleChange}
                  className="border rounded p-2 w-full bg-gray-100"
                  disabled
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Postback Link</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="postback_link"
                  value={formData.postback_link}
                  onChange={handleChange}
                  className="border rounded p-2 w-full bg-gray-100"
                  disabled
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
                  required
                >
                  {Object.entries(statusOptions).map(([value, label]) => (
                    <option key={value} value={value}>
                      {label.charAt(0).toUpperCase() + label.slice(1)}
                    </option>
                  ))}
                  
                </select>
              </td>
            </tr>

            {/* Postback Events Section */}
            <tr>
              <th className="text-left p-2 border align-top">Postback Events</th>
              <td className="p-2 border">
                {postbacks?.length > 0 ? (
                  <div className="space-y-2">
                    {postbacks.map((pb, index) => (
                      <div key={index} className="border rounded p-2 bg-gray-50">
                        <div className="flex flex-col sm:flex-row sm:items-center">
                          <span className="font-medium w-48">
                            {eventOptions.find(e => e.value === pb.event_name)?.label || pb.event_name}
                          </span>
                          <span className="text-blue-700 break-all sm:ml-4">
                            {pb.postback_url || "-"}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-gray-500 text-sm">No postbacks configured</div>
                )}
              </td>
            </tr>


              

          </tbody>
        </table>

        


        <div className="mt-4 flex justify-end">
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded"
          >
            Update Ads Campaign
          </button>
        </div>
      </form>
    </div>
  );
}

export default AdsCampaignEditPage;



