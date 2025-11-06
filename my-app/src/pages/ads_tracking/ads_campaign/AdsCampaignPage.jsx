import { useEffect, useState } from "react";
import api from "../../../api/axios";
import { useNavigate } from "react-router-dom";
import DataTable from "datatables.net-dt";
import { options } from "@fullcalendar/core/preact.js";

export default function AdsCampaignPage() {

  const [adsCampaign, setAdsCamapaign] = useState([]);
  const [filters, setFilters] = useState({

      name: "",
      target_url: "",
      status:"",
      fdate: "",
      tdate: ""
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  const [statusOptions, setStatusOptions] = useState([]);


  const fetchAdsCampaign = async () => {
      try {
          const params = new URLSearchParams();
          Object.keys(filters).forEach((key) => {
            if (filters[key]) params.append(key, filters[key]);
          });
          const res = await api.get(`/ads-campaign?${params.toString()}`);
       
          console.log(res);
          setAdsCamapaign(res.data.ads_campaigns || []);
          setStatusOptions(res.data.status || [])
      } catch (err) {
          setError(`Error fetching adsCampaign: ${err.message}`);
      }
  };

  useEffect(() => {
      fetchAdsCampaign();
  }, []);

  useEffect(() => {
      if (!adsCampaign) return;

      const table = new DataTable("#adsCampaignTable", {
          data: adsCampaign,
          destroy: true,
          columns: [
              {
                  data: "id",
                  title: "Action",
                  render: (id) =>
                      `<button class="edit-btn btn btn-primary" data-id="${id}">Edit</button>
                        <button class="report-btn btn btn-info" data-id="${id}">Report</button>
                        <button class="log-btn btn btn-secondary" data-id="${id}">Log</button>

                        `,
              },
              { data: "name", title: "Name" },
    
              { data: "target_url", title: "Target Url" },
              { data: "tracking_link", title: "Tracking Link (Marketing)" },
              { data: "postback_link", title: "Postback Link (Client)" },
              { data: "status", title: "Status" },


              { data: "created_at", title: "Created At" },
              { data: "updated_at", title: "Updated At" },
          ],
      });

      const handleClick = (e) => {
          if (e.target.classList.contains("edit-btn")) {
              const id = e.target.getAttribute("data-id");
              navigate(`/ads-tracking/campaign/edit/${id}`);
          }

          if (e.target.classList.contains("report-btn")) {
              const id = e.target.getAttribute("data-id");
              navigate(`/ads-tracking/campaign/report/${id}`);
          }

          if (e.target.classList.contains("log-btn")) {
              const id = e.target.getAttribute("data-id");
              navigate(`/ads-tracking/campaign/log/${id}`);
          }
          
          

      };

      document.addEventListener("click", handleClick);

      return () => {
          table.destroy();
          document.removeEventListener("click", handleClick);
      };
  }, [adsCampaign, navigate]);


  const handleSearch = (e) => {
      e.preventDefault();
      fetchAdsCampaign();
  };

  const handleClearSearch = () => {
      setFilters({
          name: "",
          target_url: "",
          fdate: "",
          tdate: "",
      });
      fetchAdsCampaign();
  };

  const handleAddAdsCampaign = () => {
      navigate('/ads-tracking/campaign/add');
  }

  return (
    <>
        {error && <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>}
        {success && <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>}

        <h2 className="text-xl font-bold mb-4">Ads Campaign</h2>

        {/* Search Filters */}
        <form onSubmit={handleSearch} className="flex flex-col justify-between gap-4 flex-wrap bg-white p-4 rounded shadow-sm mb-4">
            <div className="flex gap-4">
                <div>
                    <label>Name</label>
                    <input type="text" value={filters.name} onChange={(e) => setFilters({ ...filters, name: e.target.value })} className="border rounded p-2 w-full" />
                </div>
                <div>
                    <label>Target Url</label>
                    <input type="text" value={filters.target_url} onChange={(e) => setFilters({ ...filters, target_url: e.target.value })} className="border rounded p-2 w-full" />
                </div>
                <div>
                    <label>Status</label>
                    {/* <input type="text" value={filters.status} onChange={(e) => setFilters({ ...filters, status: e.target.value })} className="border rounded p-2 w-full" /> */}
                    <select
                        name="status"
                        value={filters.status}
                        onChange={(e) => setFilters({ ...filters, status: e.target.value })}
                        className="border rounded p-2 w-full"
                        
                    >
                        <option value=""></option>
                        {Object.entries(statusOptions).map(([value, label]) => (
                            <option key={value} value={value}>
                                {label.charAt(0).toUpperCase() + label.slice(1)}
                            </option>
                        ))}
                        
                    </select>
                </div>

                <div>
                    <label>From Date</label>
                    <input type="date" value={filters.fdate} onChange={(e) => setFilters({ ...filters, fdate: e.target.value })} className="border rounded p-2 w-full" />
                </div>
                <div>
                    <label>To Date</label>
                    <input type="date" value={filters.tdate} onChange={(e) => setFilters({ ...filters, tdate: e.target.value })} className="border rounded p-2 w-full" />
                </div>
            </div>
            <div className="flex justify-between">
                <div className="flex items-end gap-2">
                    <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Search</button>
                    <button type="button" onClick={handleClearSearch} className="bg-red-600 text-white px-4 py-2 rounded">Clear Search</button>
                </div>
                <div>
                    <button type="button" onClick={handleAddAdsCampaign} className="bg-blue-600 text-white px-4 py-2 rounded">Add Ads Campaign</button>
                </div>
            </div>

        </form>

        {/* DataTable */}
        <div className="bg-white shadow-md rounded-lg p-4">
            <table id="adsCampaignTable" className="display" style={{ width: "100%" }}></table>
        </div>
    </>
    
  );
}
