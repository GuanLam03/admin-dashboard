import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import DataTable from "datatables.net-dt";
import {
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

import LanguageIcon from '@mui/icons-material/Language';
import VisibilityIcon from '@mui/icons-material/Visibility';
import DevicesIcon from '@mui/icons-material/Devices';
import AndroidIcon from '@mui/icons-material/Android';
import EventIcon from '@mui/icons-material/Event';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import { Dropdown, ButtonGroup } from "react-bootstrap";
import CircularProgress from '@mui/material/CircularProgress';
import { useTranslation } from "react-i18next";
import dataTableLocales from "../../../utils/i18n/datatableLocales";

export default function AdsCampaignReportPage() {
  const {t,i18n} = useTranslation();
  const { id } = useParams();
  const [tableLoading, setTableLoading] = useState(false);
  const [stats, setStats] = useState({
    totalClicks: 0,
    totalConversions: 0,
    totalRevenue: 0,
    conversionRate: 0,
  });
  const [countryStats, setCountryStats] = useState([]);

  const colorMap = {};

  const [tableData, setTableData] = useState([]);
  const [selectedFilter, setSelectedFilter] = useState("");

  const [filters, setFilters] = useState({
        fdate: "",
        tdate: "",
  });

  
    // const filtersRef = useRef(filters);
    //   useEffect(() => {
    //     filtersRef.current = filters;
    // }, [filters]);

  const [activeFilter, setActiveFilter] = useState(null);
  const [selectedOption, setSelectedOption] = useState({});


  // Fetch campaign summary (total clicks, revenue, etc.)
  const fetchSummaryData = async () => {
    try {
      const res = await api.get(`/ads-campaign/report/${id}`);
      const summary = res.data.summary || {};
      const countryStatsFromAPI = res.data.country_stats || {};

      // Compute conversion rate (if not from backend)
      const conversionRate =
        summary.total_clicks > 0
          ? ((summary.total_conversions / summary.total_clicks) * 100).toFixed(2)
          : 0;

      setStats({
        totalClicks: summary.total_clicks || 0,
        totalConversions: summary.total_conversions || 0,
        totalRevenue: summary.total_revenue || 0,
        conversionRate,
      });

      setCountryStats(countryStatsFromAPI); // Set the country stats
    } catch (err) {
      console.error("Error fetching campaign summary:", err);
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
    fetchSummaryData();
  }, [id]);



  // Bar chart by country
  const barData = countryStats.map(stat => ({
    country: stat.Country,
    clicks: stat.Count,
    revenue: stat.TotalRevenue, // optional, if needed in chart
  }));

  // ------------------


  const filterComponents = [
    {
      type: "dropdown",
      key: "country",
      icon: <LanguageIcon />,
      options: [
        { key: "country", label: t("adsLogDetail.field.country") },
        { key: "city", label: t("adsLogDetail.field.city") },
      ],
    },
    {
      type: "button",
      key: "ip",
      icon: <VisibilityIcon />,
      label: t("adsLogDetail.field.ip"),
    },
    {
      type: "dropdown",
      key: "device",
      icon: <DevicesIcon />,
      options: [
        { key: "device_type", label: t("adsLogDetail.field.deviceType") },
        { key: "device_name", label: t("adsLogDetail.field.deviceName") },
      ],
    },
    {
      type: "dropdown",
      key: "os",
      icon: <AndroidIcon />,
      options: [
        { key: "os_name", label: t("adsLogDetail.field.osName") },
        { key: "os_version", label: t("adsLogDetail.field.osVersion") },
      ],
    },
    {
      type: "button",
      key: "event_name",
      icon: <EventIcon />,
      label: t("adsLogDetail.event"),
    },
    {
      type: "dropdown",
      key: "date",
      icon: <CalendarMonthIcon />,
      options: [
        { key: "date", label: t("adsLogDetail.day") },
        { key: "month", label: t("adsLogDetail.month") },
      ],
    },

    {
      type: "dropdown",
      key: "Day parting",
      icon: <CalendarTodayIcon />,
      options: [
        { key: "hour_of_day", label: t("adsLogDetail.hourOfDay") },
        { key: "day_of_week", label: t("adsLogDetail.dayOfWeek") },
      ],
    },

    
  ];


  const fetchFilteredData = async (filterType, parentName = filterType) => {
    try {
      setTableLoading(true);
      // Highlight dropdown / selected button
      setActiveFilter(parentName);
      setSelectedFilter(filterType);
      setSelectedOption((prev) => ({
        ...prev,
        [parentName]: filterType,
      }));

      // Fetch filtered data from backend
      const res = await api.get(`/ads-campaign/report/${id}/filter`, {
        params: { 
          type: filterType, 
          fdate: filters.fdate,
          tdate: filters.tdate,
       },
      });

      const data = res.data.data || [];
      const results = res.data || {};
      setTableData(data);
      setFilters({
        fdate: results.start_date?.split("T")[0],
        tdate: results.end_date?.split("T")[0] 
      })

      console.log("Fetched data:", data);
    } catch (err) {
      console.error("Error fetching filtered data:", err);
    } finally {
      setTableLoading(false); // stop loading spinner
    }
  };



  useEffect(() => {
    console.log("test",tableData);
    if (tableData.length === 0){
      new DataTable("#reportsTable").clear().draw();
      return;
    };
  
    // Get keys from first row to auto-generate columns
    // const columns = Object.keys(tableData[0]).map((key) => ({
    //   title: key.replace(/_/g, " ").replace(/\b\w/g, (l) => l.toUpperCase()),
    //   data: key,
    // }));

    const columns = Object.keys(tableData[0]).map((key) => {
      // Convert xxx_bbb to xxxBbb
      const camelCaseKey = key.replace(/_([a-z])/g, (_, p1) => p1.toUpperCase());
      return {
        title: t(`adsCampaign.adsCampaignReport.reportTableColumns.${camelCaseKey}`),
        data: key,
      };
    });


    // Initialize DataTable
    const reportTable = new DataTable("#reportsTable", {
      data: tableData,
      columns,
      destroy: true,
      language: dataTableLocales[i18n.language],
      searching: true,
      paging: true,
      info: true,
      ordering: false,
      columnDefs: [
        {
          targets: "_all", // applies to all columns
          className: "text-start" // Bootstrap class for left-align
        }
      ],
      
    });

    return () => {
        reportTable.destroy();
    }
  }, [tableData]);


  const handleSearch = (e) => {
    e.preventDefault(); 

    if (!selectedFilter) {
      alert("Please select a filter first");
      return;
    }

    // Call fetchFilteredData with the currently selected filter and pass dates as query params
    fetchFilteredData(selectedFilter, selectedFilter);
  };

  const handleClearSearch = () => {
    setFilters({ fdate: "", tdate: "" });
  };



  return (
    <div className="p-6 bg-gray-50 min-h-screen">
        
      <h2 className="text-2xl font-bold">{t("adsCampaign.adsCampaignReport.adsCampaignReport")} (ID: {id})</h2>

      {/* Summary Cards */}
      <div className="grid grid-cols-4 gap-4 my-4">
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">{t("adsCampaign.adsCampaignReport.totalClicks")}</h5>
          <p className="text-2xl font-bold">{stats.totalClicks}</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">{t("adsCampaign.adsCampaignReport.conversions")}</h5>
          <p className="text-2xl font-bold">{stats.totalConversions}</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">{t("adsCampaign.adsCampaignReport.conversionRate")}</h5>
          <p className="text-2xl font-bold">{stats.conversionRate}%</p>
        </div>
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="text-gray-600 text-sm">{t("adsCampaign.adsCampaignReport.revenue")}</h5>
          <p className="text-2xl font-bold">RM{stats.totalRevenue.toFixed(2)}</p>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-2 gap-6 mb-8">
       
        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="font-semibold mb-2">{t("adsCampaign.adsCampaignReport.revenueByCountry")}</h5>
          <table className="min-w-full border-collapse border border-gray-200">
            <thead>
              <tr className="bg-gray-100">
                <th className="border border-gray-300 px-4 py-2 text-left">{t("adsCampaign.adsCampaignReport.country")}</th>
                <th className="border border-gray-300 px-4 py-2 text-right">{t("adsCampaign.adsCampaignReport.revenue")} (RM)</th>
              </tr>
            </thead>
            <tbody>
              {(() => {
                // Sort descending by revenue
                const sortedData = [...barData].sort((a, b) => b.revenue - a.revenue);
                const highestRevenue = sortedData.length > 0 ? sortedData[0].revenue : 0;

                return sortedData.map(({ country, revenue }, index) => (
                  <tr
                    key={index}
                    className="hover:bg-gray-100"
                  >
                    <td className="border border-gray-300 px-4 py-2">{country}</td>
                    <td className={`border border-gray-300 px-4 py-2 text-right ${
                      revenue === highestRevenue ? "text-green-600 font-semibold" : ""
                    }`}>
                      {revenue.toFixed(2)}
                    </td>
                  </tr>
                ));
              })()}
            </tbody>
          </table>
        </div>

        <div className="bg-white p-4 rounded shadow-sm">
          <h5 className="font-semibold mb-2">{t("adsCampaign.adsCampaignReport.clicksByCountry")}</h5>
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

       <section>
        <h4>{t("adsCampaign.adsCampaignReport.filter")}</h4>
        <div className="d-flex flex-wrap gap-3 p-4 shadow-sm bg-white">
          {filterComponents.map((filter) =>
            filter.type === "button" ? (
              <button
                key={filter.key}
                className={`btn d-flex align-items-center gap-2 ${
                  selectedFilter === filter.key ? "btn-primary" : "btn-light"
                }`}
                onClick={() => fetchFilteredData(filter.key)}
              >
                {filter.icon}
                {filter.label}
              </button>
            ) : (
              <Dropdown as={ButtonGroup} key={filter.key}>
                <Dropdown.Toggle
                  variant={
                    filter.options.some((o) => o.key === selectedFilter)
                      ? "primary"
                      : "light"
                  }
                  className="d-flex align-items-center gap-2"
                >
                  {filter.icon}
                  {filter.options.find((o) => o.key === selectedFilter)?.label ||
                    filter.options[0].label}
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  {filter.options.map((option) => (
                    <Dropdown.Item
                      key={option.key}
                      onClick={() => fetchFilteredData(option.key)}
                    >
                      {option.label}
                    </Dropdown.Item>
                  ))}
                </Dropdown.Menu>
              </Dropdown>
            )
          )}
        </div>



        <form onSubmit={handleSearch} className="p-4 shadow-sm bg-white mb-4">
            <div className="flex gap-4">
                <div>
                    <label>{t("fromDate")}</label>
                    <input
                    type="date"
                    value={filters.fdate}
                    onChange={(e) => setFilters({ ...filters, fdate: e.target.value })}
                    className="border rounded p-2 w-full"
                    />
                </div>
                <div>
                    <label> {t("toDate")}</label>
                    <input
                    type="date"
                    value={filters.tdate}
                    onChange={(e) => setFilters({ ...filters, tdate: e.target.value })}
                    className="border rounded p-2 w-full"
                    />
                </div>

                <div className="flex gap-2 mt-4">
                  <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">
                       {t("common.buttons.search")}
                  </button>
                  <button
                      type="button"
                      onClick={handleClearSearch}
                      className="bg-red-600 text-white px-4 py-2 rounded"
                  >
                      {t("common.buttons.clearSearch")}
                  </button>
                </div>
            </div>
          </form>

          {/* Report List Table */}
          <div className="bg-white shadow-md rounded-lg p-4">            
            <table
              id="reportsTable"
              className="display"
              style={{ width: "100%" }}
            ></table>
          </div>

          {tableLoading && (
          <div className="fixed top-0 left-0 w-[100vw] h-[100vh] flex items-center justify-center z-50">
            {/* Background overlay */}
            <div className="absolute inset-0 bg-white opacity-40"></div>

            {/* Spinner on top */}
            <div className="relative z-10">
              <CircularProgress sx={{ color: '#0f4bceff' }} />
            </div>
          </div>
)}

      </section>
    </div>
  );
}
