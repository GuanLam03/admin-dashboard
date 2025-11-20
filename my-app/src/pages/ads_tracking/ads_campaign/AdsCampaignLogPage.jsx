import { useEffect, useRef, useState } from "react";
import DataTable from "datatables.net-dt";
import api from "../../../api/axios";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import dataTableLocales from "../../../utils/i18n/datatableLocales";

export default function AdsCampaignLogPage() {
    const {t,i18n} = useTranslation();
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
                ...dataTableLocales[i18n.language], // apply correct locale
                
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
                    recordsTotal: res.data.recordsTotal,
                    recordsFiltered: res.data.recordsFiltered,
                    data: results,
                });
            });
        },


        columns: [
            { data: "id", title: t("adsLogDetail.field.id") },
            { data: "ip", title: t("adsLogDetail.field.ip") },
            { data: "country", title: t("adsLogDetail.field.country") },
            { data: "region", title: t("adsLogDetail.field.region") },
            { data: "city", title: t("adsLogDetail.field.city") },
            { data: "user_agent", title: t("adsLogDetail.field.userAgent") },
            { data: "referrer", title: t("adsLogDetail.field.referrer") },
            { data: "device_type", title: t("adsLogDetail.field.deviceType") },
            { data: "device_name", title: t("adsLogDetail.field.deviceName") },
            { data: "os_name", title: t("adsLogDetail.field.osName") },
            { data: "os_version", title: t("adsLogDetail.field.osVersion") },
            { data: "browser_name", title: t("adsLogDetail.field.browserName") },
            { data: "browser_version", title: t("adsLogDetail.field.browserVersion") },
            { data: "created_at", title: t("adsLogDetail.field.createdAt") },
            { data: "updated_at", title: t("adsLogDetail.field.updatedAt") }
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
                        <label>{t("adsLogDetail.field.ip")}</label>
                        <input
                        type="text"
                        value={filters.ip}
                        onChange={(e) => setFilters({ ...filters, ip: e.target.value })}
                        className="border rounded p-2 w-full"
                        />
                    </div>
                    <div>
                        <label>{t("adsLogDetail.field.country")}</label>
                        <input
                        type="text"
                        value={filters.country}
                        onChange={(e) => setFilters({ ...filters, country: e.target.value })}
                        className="border rounded p-2 w-full"
                        />
                    </div>
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
                        <label>{t("tromDate")}</label>
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
            </form>


            {/* Logs Table */}
            <div className="bg-white shadow-sm rounded p-4">
                <h5 className="font-semibold mb-4">{t("adsLogDetail.detailedLogs")}</h5>
                <table id="adsLogsTable" className="display" style={{ width: "100%" }}></table>
            </div>


        </div>

    
        
        
    )
}