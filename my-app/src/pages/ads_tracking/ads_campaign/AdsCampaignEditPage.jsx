import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../../../api/axios";
import Tooltip from "@mui/material/Tooltip";
import CloseIcon from '@mui/icons-material/Close';

const eventOptions = [
  { value: "PURCHASE", label: "Purchase" },
  { value: "COMPLETE_REGISTRATION", label: "Registration" },
  { value: "FORM_SUBMIT", label: "Form Submit" },
];

function AdsCampaignEditPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({});
  const [postbacks, setPostbacks] = useState([]);
  const [statusOptions, setStatusOptions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const fetchCampaign = async () => {
      try {
        const res = await api.get(`/edit-ads-campaign/${id}`);
        const result = res.data.ads_campaign;
        const campaignPostbacks = res.data.ads_campaign_postbacks || [];
        const statusList = res.data.status || [];

        setFormData({
          name: result.name || "",
          target_url: result.target_url || "",
          tracking_link: result.tracking_link || "",
          postback_link: result.postback_link || "",
          status: result.status || "",
        });
        setPostbacks(campaignPostbacks);
    
        setStatusOptions(statusList);
      } catch (err) {
        console.error(err);
        setError("Failed to load campaign details.");
      } finally {
        setLoading(false);
      }
    };
    fetchCampaign();
  }, [id]);

  const handleFormChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      await api.post(`/edit-ads-campaign/${id}`, { ...formData, postback_events: postbacks });
      setSuccess("Ads Campaign updated successfully!");
      setTimeout(() => navigate("/ads-tracking/campaign"), 1200);
    } catch (err) {
      console.error(err);
      setError("Failed to update campaign. Please try again.");
    }
  };

  if (loading) return <p>Loading Ads Campaign...</p>;

  return (
    <div>
      {error && <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>}
      {success && <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>}

      <h2 className="text-xl font-bold mb-4">Edit Ads Campaign</h2>

      <form onSubmit={handleSubmit} className="bg-white p-4 rounded shadow-sm w-full">
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
                  onChange={handleFormChange}
                  className="border rounded p-2 w-full"
                  required
                />
              </td>
            </tr>
            <tr>
              <th className="text-left p-2 border">Target URL</th>
              <td className="p-2 border">
                <input
                  type="text"
                  name="target_url"
                  value={formData.target_url}
                  onChange={handleFormChange}
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
                  value={formData.tracking_link}
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
                  value={formData.postback_link}
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
                  onChange={handleFormChange}
                  className="border rounded p-2 w-full"
                  required
                >
                  {Object.entries(statusOptions).map(([val, label]) => (
                    <option key={val} value={val}>
                      {label.charAt(0).toUpperCase() + label.slice(1)}
                    </option>
                  ))}
                </select>
              </td>
            </tr>

            {/* Editable Postback Events */}
            {/* Postback Events Section */}
            <tr>
              <th className="text-left p-2 border align-top">Postback Events</th>
              <td className="p-2 border">
                {postbacks.length > 0 ? (
                  <div className="space-y-3">
                    {postbacks.map((pb, index) => {
                      // Compute remaining available events (prevent duplicates)
                      const usedEvents = postbacks.map(p => p.event_name);
                      const availableEvents = eventOptions.filter(
                        opt => opt.value === pb.event_name || !usedEvents.includes(opt.value)
                      );

                      return (
                        <div
                          key={index}
                          className="relative border rounded p-3 bg-gray-50 space-y-2"
                        >
                          <div className="flex flex-col sm:flex-row sm:items-center gap-3">
                            <select
                              name="event_name"
                              value={pb.event_name}
                              onChange={(e) => {
                                const updated = [...postbacks];
                                updated[index].event_name = e.target.value;
                                setPostbacks(updated);
                              }}
                              className="border rounded p-2 w-full sm:w-48"
                              required
                            >
                              <option value="">-- Select Event --</option>
                              {availableEvents.map((opt) => (
                                <option key={opt.value} value={opt.value}>
                                  {opt.label}
                                </option>
                              ))}
                            </select>

                            <input
                              type="text"
                              name="postback_url"
                              placeholder="Enter Postback URL"
                              value={pb.postback_url || ""}
                              onChange={(e) => {
                                const updated = [...postbacks];
                                updated[index].postback_url = e.target.value;
                                setPostbacks(updated);
                              }}
                              className="border rounded p-2 flex-1"
                              required
                            />

                            <div className="flex items-center gap-2 text-sm text-gray-700 mr-6">
                              <input
                                type="checkbox"
                                checked={!!pb.include_click_params}
                                onChange={(e) => {
                                  const updated = [...postbacks];
                                  updated[index].include_click_params = e.target.checked;
                                  setPostbacks(updated);
                                }}
                                
                              />
                  
                               <Tooltip title="Send UTM and click data (source, campaign, etc.) for this postback.">
                                  <label
                                    className="text-sm text-gray-700"
                                  >
                                    Include original tracking parameters
                                  </label>
                                </Tooltip>
                            </div>

                            
                          </div>
                          <button
                              type="button"
                              onClick={() => {
                                const updated = postbacks.filter((_, i) => i !== index);
                                setPostbacks(updated);
                              }}
                              className="absolute right-2 top-1"
                            >
                              <CloseIcon fontSize="small"/>
                          </button>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div className="text-gray-500 text-sm">No postbacks configured</div>
                )}

                {/* Add new postback button */}
                {postbacks.length < eventOptions.length && (
                  <button
                    type="button"
                    onClick={() =>
                      setPostbacks([
                        ...postbacks,
                        {
                          id: null,
                          event_name: "",
                          postback_url: "",
                          include_click_params: false,
                        },
                      ])
                    }
                    className="mt-3 bg-green-600 text-white px-3 py-1 rounded"
                  >
                    + Add Postback Event
                  </button>
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
