import { useEffect, useRef, useState } from "react";
import DataTable from "datatables.net-dt";
import api from "../../../api/axios";
import { useParams } from "react-router-dom";

export default function AdsCampaignLogPage() {
    const { id } = useParams();
    const [filters, setFilters] = useState({
        ip: "",
        country: "",
        fdate: "",
        tdate: "",
    });

    const filtersRef = useRef(filters);
      useEffect(() => {
        filtersRef.current = filters;
    }, [filters]);
    

    useEffect(() => {
        // Initialize DataTable
        const table = new DataTable("#adsLogsTable", {
            processing: true,
            serverSide: true,
            language: {
                info: "Showing _START_ to _END_ of _TOTAL_ entries", // hides the "filtered" part
            },
            ajax: function (data, callback) {
                const currentFilters = filtersRef.current; // always latest filters

                const params = {
                ...data,
                ip: currentFilters.ip || "",
                country: currentFilters.country || "",
                fdate: currentFilters.fdate || "",
                tdate: currentFilters.tdate || "",
                };
                

                api.get(`ads-campaign/report/ads-log-details/${id}`, { params })
                .then((res) => {
                const results = res.data.data;

                callback({
                    draw: data.draw,
                    // recordsTotal: res.data.recordsTotal,
                    recordsFiltered: res.data.recordsFiltered,
                    data: results,
                });
            });
        },


        columns: [
            { data: "id", title: "ID" },
            { data: "ip", title: "IP Address" },
            { data: "country", title: "Country" },
            { data: "region", title: "Region" },
            { data: "city", title: "City" },
            { data: "user_agent", title: "User Agent" },
            { data: "referrer", title: "Referrer" },
            { data: "device_type", title: "Device Type" },
            { data: "device_name", title: "Device Name" },
            { data: "os_name", title: "OS Name" },
            { data: "os_version", title: "OS Version" },
            { data: "browser_name", title: "Browser" },
            { data: "browser_version", title: "Browser Version" },
            { data: "created_at", title: "Created At" },
            { data: "updated_at", title: "Updated At" }
        ]
        });

        window.adsLogsTable = table;
        return () => table.destroy();
    }, [id]);


    const handleSearch = (e) => {
        e.preventDefault();
        if (window.adsLogsTable) {
            window.adsLogsTable.ajax.reload(); // Refresh table with new filters
        }
    };

    const handleClearSearch = () => {
        setFilters({ ip: "", country: "", fdate: "", tdate: "" });
        if (window.adsLogsTable) {
            window.adsLogsTable.ajax.reload();
        }
    };

    return (
        <div className="p-6 bg-gray-50 min-h-screen">
            <h2 className="text-2xl font-bold">Ads Campaign Log (ID: {id})</h2>
            {/* Search Filters */}
            <form onSubmit={handleSearch} className="bg-white p-4 rounded shadow-sm my-4">
                <div className="flex gap-4">
                    <div>
                        <label>IP Address</label>
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
                        <label>From Date</label>
                        <input
                        type="date"
                        value={filters.fdate}
                        onChange={(e) => setFilters({ ...filters, fdate: e.target.value })}
                        className="border rounded p-2 w-full"
                        />
                    </div>
                    <div>
                        <label>To Date</label>
                        <input
                        type="date"
                        value={filters.tdate}
                        onChange={(e) => setFilters({ ...filters, tdate: e.target.value })}
                        className="border rounded p-2 w-full"
                        />
                    </div>
                </div>

                <div className="flex gap-2 mt-4">
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


            {/* Logs Table */}
            <div className="bg-white shadow-sm rounded p-4">
                <h5 className="font-semibold mb-4">Detailed Logs</h5>
                <table id="adsLogsTable" className="display" style={{ width: "100%" }}></table>
            </div>


        </div>

    
        
        
    )
}