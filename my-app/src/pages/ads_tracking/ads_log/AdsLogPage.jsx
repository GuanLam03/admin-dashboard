import { useEffect, useState } from "react";
import api from "../../../api/axios";
import { useNavigate, useParams } from "react-router-dom";
import DataTable from "datatables.net-dt";

export default function AdsLogPage() {
  const navigate = useNavigate();
  const { campaign_id } = useParams(); // optional route param if you show per-campaign logs

  const [adsLogs, setAdsLogs] = useState([]);
  const [filters, setFilters] = useState({
    ip: "",
    country: "",
    region: "",
    city: "",
    converted: "",
    fdate: "",
    tdate: "",
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Fetch ads logs from backend
  const fetchAdsLogs = async () => {
    try {
      const params = new URLSearchParams();
      Object.keys(filters).forEach((key) => {
        if (filters[key]) params.append(key, filters[key]);
      });

      if (campaign_id) params.append("ads_campaign_id", campaign_id);

      const res = await api.get(`/ads-logs?${params.toString()}`);
      setAdsLogs(res.data.ads_logs || []);
      console.log(res.data.ads_logs)
    } catch (err) {
      setError(`Error fetching ads logs: ${err.message}`);
    }
  };

  useEffect(() => {
    fetchAdsLogs();
  }, [campaign_id]);

  // Build DataTable
  useEffect(() => {
    if (!adsLogs) return;

    const table = new DataTable("#adsLogsTable", {
      data: adsLogs,
      destroy: true,
      columns: [
        { data: "id", title: "ID" },
        { data: "ads_campaign_id", title: "Campaign ID" },
        { data: "ip", title: "IP" },
        { data: "country", title: "Country" },
        { data: "region", title: "Region" },
        { data: "city", title: "City" },
        { data: "referrer", title: "Referrer" },
        { data: "user_agent", title: "User Agent" },
        {
          data: "converted",
          title: "Converted",
          render: (val) => (val ? " Yes" : " No"),
        },
        { data: "client_product_id", title: "Product ID" },
        { data: "value", title: "Value" },
        { data: "created_at", title: "Created At" },
      ],
    });

    return () => table.destroy();
  }, [adsLogs]);

  const handleSearch = (e) => {
    e.preventDefault();
    fetchAdsLogs();
  };

  const handleClearSearch = () => {
    setFilters({
      ip: "",
      country: "",
      region: "",
      city: "",
      converted: "",
      fdate: "",
      tdate: "",
    });
    fetchAdsLogs();
  };

  return (
    <>
      {error && <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>}
      {success && <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>}

      <h2 className="text-xl font-bold mb-4">
        Ads Logs {campaign_id ? `(Campaign #${campaign_id})` : ""}
      </h2>

      {/* Filters */}
      <form onSubmit={handleSearch} className="bg-white p-4 rounded shadow-sm mb-4 flex flex-wrap gap-4">
        <div>
          <label>IP</label>
          <input
            type="text"
            value={filters.ip}
            onChange={(e) => setFilters({ ...filters, ip: e.target.value })}
            className="border rounded p-2 w-full"
          />
        </div>
        <div>
          <label>Country</label>
          <input
            type="text"
            value={filters.country}
            onChange={(e) => setFilters({ ...filters, country: e.target.value })}
            className="border rounded p-2 w-full"
          />
        </div>
        <div>
          <label>Converted</label>
          <select
            value={filters.converted}
            onChange={(e) => setFilters({ ...filters, converted: e.target.value })}
            className="border rounded p-2 w-full"
          >
            <option value="">All</option>
            <option value="1">Converted</option>
            <option value="0">Not Converted</option>
          </select>
        </div>
        <div className="flex items-end gap-2">
          <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">
            Search
          </button>
          <button
            type="button"
            onClick={handleClearSearch}
            className="bg-red-600 text-white px-4 py-2 rounded"
          >
            Clear
          </button>
        </div>
      </form>

      {/* DataTable */}
      <div className="bg-white shadow-md rounded-lg p-4">
        <table id="adsLogsTable" className="display" style={{ width: "100%" }}></table>
      </div>
    </>
  );
}
