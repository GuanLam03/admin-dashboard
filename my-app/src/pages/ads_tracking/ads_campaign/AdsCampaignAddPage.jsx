import { useEffect, useState } from "react";
import api from "../../../api/axios";
import Tooltip from "@mui/material/Tooltip";

const eventOptions = [
  { value: "PURCHASE", label: "Purchase"},
  { value: "COMPLETE_REGISTRATION", label: "Registration"},
  { value: "FORM_SUBMIT", label: "Form Submit"}
]

export default function AdsCampaignAddPage() {
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState("");
  const [formData, setFormData] = useState({
    name: "",
    target_url: ""
  });

  const [generatedLinks, setGeneratedLinks] = useState(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState({ tracking: false, postback: false });

  const [postbackEnabled, setPostbackEnabled] = useState(false);
  const [selectedEvents, setSelectedEvents] = useState([]);
  const [postbackUrls, setPostbackUrls] = useState({});

  const [parameters,setParameters] = useState([])

  useEffect(() => {
    const parameters = async () => {
      try{
        const res = await api.get("/add-ads-campaign/support-parameters");
        setParameters(res.data.support_parameter || {});

      }catch(err){
        setError(err);
      }
    }
    parameters();
  },[])

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
      const payload = {
        ...formData,
        postback_enabled: postbackEnabled,
        postback_events: selectedEvents.map(event => ({
          event_name: event.value,
          url: postbackUrls[event.value]?.url || "",
          include_click_params: postbackUrls[event.value]?.include_click_params || false,
        })),
      };

     

      const res = await api.post("/add-ads-campaign", payload);
      setGeneratedLinks({
        trackingLink: res.data.tracking_link,
        postbackLink: res.data.postback_link,
      });
      setSuccess(res.data.status_name);
    } catch (err) {
      console.error("Error generating campaign:", err);
      setGeneratedLinks(null);//reset 
      if (err.response?.status === 422) {
        setError( err.response.data.errors);
      } else {
        console.error("Error generating campaign:", err);
      }
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


 const toggleEvent = (index) => {
    const eventValue = eventOptions[index];

    setSelectedEvents((prev) => {
      if (prev.includes(eventValue)) {
        return prev.filter((v) => v !== eventValue);
      } else {
        return [...prev, eventValue];
      }
    });
  };

  return (
    <>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {Object.entries(error).map(([field, messages]) => (
            <div key={field}>
              {Object.values(messages).join(", ")}
            </div>
          ))}
        </div>
      )}
      {success && <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>}
      <h2 className="text-xl font-bold mb-4">Ads Campaign</h2>
      <div className="flex gap-4">

        <div className="p-6 bg-white shadow-sm rounded-sm w-full">
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


            <div className="font-medium flex items-center justify-between">
              <label>Postback</label>
              <label className="inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" checked={postbackEnabled}
                  onChange={(e) => {
                      const isChecked = e.target.checked;
                      setPostbackEnabled(isChecked);

                      if (!isChecked) {
                        setSelectedEvents([]);
                      }
                    }}

                />
                <div className="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600 dark:peer-checked:bg-blue-600"></div>
              </label>
            </div>

            {
              postbackEnabled && (
                <div>
                  <label className="block font-medium">Choose the event to postback:</label>
                  <div className="flex flex-col gap-2 mt-2 ml-2">
                    {eventOptions.map((event,index) => {
                      return ( 
                        <div key={index} className="flex items-center justify-between">
                          <label className="block">{index + 1}. {event.label}</label>
                          <input type="checkbox" onChange={() => toggleEvent(index)}></input>
                        </div>
                      )
                      
                    })}
                  </div>
                  
                
                </div>
              )
            }

            {selectedEvents.map(event => (
              <div key={event.value} className="mb-2 bg-gray-50 rounded-sm p-4">
                <label className="block font-medium">Enter {event.label} Postback URL</label>
                <input
                  type="url"
                  value={postbackUrls[event.value]?.url || ""}
                  onChange={(e) =>
                    setPostbackUrls((prev) => ({
                      ...prev,
                      [event.value]: {
                        ...prev[event.value],
                        url: e.target.value,
                      },
                    }))
                  }
                  className="w-full border p-2 rounded"
                  required
                />


                <div className="flex items-center mt-2 gap-2">
                  <input
                    id={`${event.value}_include_params`}
                    type="checkbox"
                    checked={postbackUrls[event.value]?.include_click_params || false}
                    onChange={(e) =>
                      setPostbackUrls(prev => ({
                        ...prev,
                        [event.value]: {
                          ...prev[event.value],
                          include_click_params: e.target.checked,
                        },
                      }))
                    }
                    className="mr-2"
                  />

                  <Tooltip title="Send UTM and click data (source, campaign, etc.) for this postback.">
                    <label
                      htmlFor={`${event.value}_include_params`}
                      className="text-sm text-gray-700"
                    >
                      Include original tracking parameters
                    </label>
                  </Tooltip>
                  
                </div>

                
              </div>
            ))}

            





            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 text-white mt-4 px-4 py-2 rounded hover:bg-blue-700"
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
                    {!copied.tracking ? <span>Copy</span> : (
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
                    {!copied.postback ? <span>Copy</span> : (
                      <span className="text-green-600 text-sm">Copied ✓</span>
                    )}
                  </button>

                </div>
              </div>
            </div>
          )}
        </div>

        {(postbackEnabled && selectedEvents.length > 0) && (
          <section className="min-w-[30%] p-6 bg-white shadow-sm rounded-sm">
              <h5>Supported Parameter</h5>
             
              <ul className="mt-4 grid grid-cols-2 gap-2">
                {
                  parameters.map( (p,index) => {
                    return <li key={index}>{p}</li>
                  })
                }

                
              </ul>
            
              
              
            </section>
          )}
      </div>
    
    </>

  );
}
