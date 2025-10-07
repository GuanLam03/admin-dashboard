import { useState } from "react";
import api from "../../../api/axios";

export default function AdsCampaignAddPage() {
  const [formData, setFormData] = useState({
    name: "",
    target_url: ""
  });

  const [generatedLinks, setGeneratedLinks] = useState(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState({ tracking: false, postback: false });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      // Call backend API to create campaign
      console.log("FOrm: " ,formData)
      const res = await api.post("/add-ads-campaign", formData);
      setGeneratedLinks({
        trackingLink: res.data.tracking_link,
        postbackLink: res.data.postback_link,
      }); 
      // Expected backend response:
      // {
      //   trackingLink: "https://middleman.com/cd67890",
      //   postbackTemplate: "https://middleman.com/postback/cd67890?log_id={logId}&value={value}&productid={productId}"
      // }
    } catch (err) {
      console.error("Error generating campaign:", err);
    }
    setLoading(false);
  };

  

  const copyToClipboard = (type, text) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopied((prev) => ({ ...prev, [type]: true }));
      setTimeout(() => {
        setCopied((prev) => ({ ...prev, [type]: false }));
      }, 1000);
    });
  };

  return (
    <>
        <h2 className="text-xl font-bold mb-4">Ads Campaign</h2>
        <div className="p-6 bg-white shadow-sm rounded-sm">
        <h2 className="text-xl font-bold mb-4">Create Ad Campaign Link</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
            <div>
            <label className="block font-medium">Ad Name</label>
            <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className="w-full border p-2 rounded"
                required
            />
            </div>

            <div>
            <label className="block font-medium">Target URL</label>
            <input
                type="url"
                name="target_url"
                value={formData.target_url}
                onChange={handleChange}
                className="w-full border p-2 rounded"
                required
            />
            </div>

            <button
            type="submit"
            disabled={loading}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
            >
            {loading ? "Generating..." : "Generate Link"}
            </button>
        </form>

        {generatedLinks && (
            <div className="mt-6 p-4 border rounded bg-gray-50">
            <h3 className="font-semibold mb-2">Generated Links:</h3>

            {/* Tracking Link */}
            <div className="mb-3">
              <strong>Tracking Link:</strong>
              <div className="flex justify-between items-center gap-2 mt-1 flex-wrap">
                <a
                  href={generatedLinks.trackingLink}
                  target="_blank"
                  rel="noreferrer"
                  className="text-blue-600 underline break-all"
                >
                  {generatedLinks.trackingLink}
                </a>
                <button
                  onClick={() => copyToClipboard("tracking", generatedLinks.trackingLink)}
                  className="px-2 py-1 text-sm text-blue-500 rounded hover:text-blue-300 transition"
                >
                  {!copied.tracking ? <span>Copy</span> :(
                    <span className="text-green-600 text-sm">Copied ✓</span>
                  )}
                </button>
                
              </div>
            </div>

            {/* Postback Link */}
            <div>
              <strong>Postback Link:</strong>
              <div className="flex justify-between items-center gap-2 mt-1 flex-wrap">
                <code className="bg-gray-100 p-1 rounded text-sm break-all">
                  {generatedLinks.postbackLink}
                </code>
                <button
                  onClick={() => copyToClipboard("postback", generatedLinks.postbackLink)}
                  className="px-2 py-1 text-sm text-blue-500 rounded hover:text-blue-300 transition"
                >
                  {!copied.postback ? <span>Copy</span> :(
                    <span className="text-green-600 text-sm">Copied ✓</span>
                  )}
                </button>
                
              </div>
            </div>
          </div>
        )}
        </div>
    </>
    
  );
}
