import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import DataTable from "datatables.net-dt";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  XAxis,
  YAxis,
  CartesianGrid,
  Bar,
  Legend,
} from "recharts";
import api from "../../../api/axios";

export default function AdsCampaignReportPage() {
  const { id } = useParams();
  const [logs, setLogs] = useState([]);
  const [stats, setStats] = useState({
    totalClicks: 0,
    totalConversions: 0,
    totalRevenue: 0,
    conversionRate: 0,
  });

  const colorMap = {};
  const COLORS = ["red", "green", "yellow", "orange"];

  const fetchLogs = async () => {
    try {
      const res = await api.get(`/ads-campaign/report/${id}`);
      console.log(res.data)
      // const logsData = res.data.data || [];
      // const summary = res.data.summary || {};

      // // compute conversion rate (if not from backend)
      // const conversionRate =
      //   summary.total_clicks > 0
      //     ? ((summary.total_conversions / summary.total_clicks) * 100).toFixed(2)
      //     : 0;

      // setLogs(logsData);
      // setStats({
      //   totalClicks: summary.total_clicks || 0,
      //   totalConversions: summary.total_conversions || 0,
      //   totalRevenue: summary.total_revenue || 0,
      //   conversionRate,
      // });
    } catch (err) {
      console.error("Error fetching logs:", err);
    }
  };

 
    const getColor = (country) => {
    if (!colorMap[country]) {
        const colors = ["#8884d8", "#82ca9d", "#ffc658", "#ff7f50", "#a4de6c", "#d0ed57"];
        colorMap[country] = colors[Object.keys(colorMap).length % colors.length];
    }
    return colorMap[country];
    };



  useEffect(() => {
    fetchLogs();
  }, [id]);

  useEffect(() => {
    if (!logs.length) return;

    const table = new DataTable("#adsLogsTable", {
      data: logs,
      destroy: true,
      columns: [
        { data: "id", title: "ID" },
        { data: "ip", title: "IP Address" },
        { data: "country", title: "Country" },
        { data: "city", title: "City" },
        { data: "referrer", title: "Referrer" },
        {
          data: "converted",
          title: "Converted",
          render: (v) => (v ? "Yes" : "No"),
        },
        { data: "value", title: "Value" },
        { data: "created_at", title: "Date" },
      ],
    });

    return () => table.destroy();
  }, [logs]);

  // Pie chart for conversion breakdown
  const pieData = [
    { name: "Converted", value: stats.totalConversions },
    { name: "Not Converted", value: stats.totalClicks - stats.totalConversions },
  ];

  // Bar chart by country
  const countryStats = logs.reduce((acc, log) => {
    if (!log.country) return acc;
    if (!acc[log.country]) acc[log.country] = { country: log.country, clicks: 0 };
    acc[log.country].clicks++;
    return acc;
  }, {});
  const barData = Object.values(countryStats);

  return (
    <div className="p-6 bg-gray-50 min-h-screen">
      <h2 className="text-2xl font-bold">
        Ads Campaign Report (ID: {id})
      </h2>

      {/* Summary Cards */}
      <div className="grid grid-cols-4 gap-4 my-4">
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">Total Clicks</h5>
          <p className="text-2xl font-bold">{stats.totalClicks}</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">Conversions</h5>
          <p className="text-2xl font-bold">{stats.totalConversions}</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">Conversion Rate</h5>
          <p className="text-2xl font-bold">{stats.conversionRate}%</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">Revenue</h5>
          <p className="text-2xl font-bold">RM{stats.totalRevenue.toFixed(2)}</p>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-2 gap-6 mb-8">
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="font-semibold mb-2">Conversion Breakdown</h5>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={pieData}
                dataKey="value"
                nameKey="name"
                outerRadius={80}
                label
              >
                {pieData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="font-semibold mb-2">Clicks by Country</h5>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="5 5" />
              <XAxis dataKey="country" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="clicks" name="Clicks">
                {barData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={getColor(entry.country)} />
                ))}
               </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Logs Table */}
      {/* <div className="bg-white shadow-sm rounded p-4">
        <h5 className="font-semibold mb-4">Detailed Logs</h5>
        <table
          id="adsLogsTable"
          className="display"
          style={{ width: "100%" }}
        ></table>
      </div> */}
    </div>
  );
}
