import { useEffect, useRef, useState } from "react";
import DataTable from "datatables.net-dt";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";

import FullCalendar from '@fullcalendar/react'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';
import { useTranslation } from "react-i18next";
import dataTableLocales from "../../utils/i18n/datatableLocales";





function SchedulePage() {
  const {t,i18n} = useTranslation();
  const [schedules, setSchedules] = useState([]);
  const [calendarUiEvents, setCalendarUiEvents] = useState([]);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const calendarRef = useRef(null);

  const [filters, setFilters] = useState({
    name: "",
    status: "active",
    fdate: "",
    tdate: "",
  });



  function formatDateTime(isoString) {
    if (!isoString) return "-";
    const date = new Date(isoString);

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero indexed
    const day = String(date.getDate()).padStart(2, "0");

    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");
    const seconds = String(date.getSeconds()).padStart(2, "0");

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  }


  const fetchSchedules = async () => {
    try {
      const res = await api.get("/schedules", {
        params: {
          title: filters.name,
          status: filters.status,
          fdate: filters.fdate,
          tdate: filters.tdate,
        },
      });

      const data = res.data.message || [];
      setSchedules(data.map((s) => ({
        id: s.id,
        title: s.title,
        recurrence: s.recurrence || "-",
        start_at: formatDateTime(s.start_at),
        end_at: formatDateTime(s.end_at),
        status: s.status,
        google_event_id: s.google_event_id || "-",
        created_at: s.created_at,
        updated_at: s.updated_at,
      })));
      // setCalendarUiEvents(data.map((s) => ({ id: s.id, title: s.title, start: s.start_at, end: s.end_at, })));
      setCalendarUiEvents(
        data.map((s) => {
          const startDate = new Date(s.start_at);
          const endDate = new Date(s.end_at);

          const startTime = startDate.toISOString().split('T')[1].slice(0, 5); // Extract time (HH:mm)
          const endTime = endDate.toISOString().split('T')[1].slice(0, 5);     // Extract time (HH:mm)

          // Get the day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
          const startDay = startDate.getDay();
          

          if (s.recurrence === "daily") {
            return {
              title: s.title,
              daysOfWeek: [0, 1, 2, 3, 4, 5, 6], // Repeat every day
              startTime,
              endTime,
              startRecur: s.start_at.split("T")[0],
              groupId: `recurring-${s.id}`,
              
            };
          }

          

          if (s.recurrence === "weekly") {
            // If the event should repeat on multiple days (e.g., Monday to Friday), we can use a range
            const daysOfWeek = [1, 2, 3, 4, 5];  // Monday to Friday (1 to 5)
            
            return {
              title: s.title,
              daysOfWeek: [[startDay]],  // Use dynamic days, or specific logic for user selection
              startTime,
              endTime,
              startRecur: s.start_at.split("T")[0], // Use the start date for recurrence
              groupId: `recurring-${s.id}`,
             
            };
          }

      

          // For non-recurring events
          return {
            id: s.id,
            title: s.title,
            start: s.start_at,
            end: s.end_at,
          };
        })
      );


    } catch (err) { console.error("Error fetching schedules:", err); setError("Error fetching schedules: " + err.message); }
  };


  useEffect(() => {
    fetchSchedules();
  }, []);

  useEffect(() => {
    if (!schedules) return;

    const table = new DataTable("#schedulesTable", {
      data: schedules,
      destroy: true,
      language: dataTableLocales[i18n.language],
      columns: [
        { data: "title", title: t("scheduleManagement.fields.title"), className: "dt-left" },
        { data: "recurrence", title: t("scheduleManagement.fields.recurrence"), className: "dt-left" },
        { data: "start_at", title: t("scheduleManagement.fields.startAt"), className: "dt-left" },
        { data: "end_at", title: t("scheduleManagement.fields.endAt"), className: "dt-left" },
        { data: "status", title: t("scheduleManagement.fields.status"), className: "dt-left" },
        { data: "google_event_id", title: t("scheduleManagement.fields.googleEventId"), className: "dt-left" },
        {
          data: "id",
          title: t("common.labels.action"),
          render: (id) =>
            `<button class="edit-btn btn btn-primary" data-id="${id}">${t("common.buttons.edit")}</button>`,
        },
      ],
    });

    const handleClick = (e) => {
      if (e.target && e.target.classList.contains("edit-btn")) {
        const scheduleId = e.target.getAttribute("data-id");
        navigate(`/schedules/edit/${scheduleId}`);
      }
    };

    document.addEventListener("click", handleClick);

    return () => {
      table.destroy();
      document.removeEventListener("click", handleClick);
    };
  }, [schedules, navigate]);


  const handleDateClick = (info) => {
    const calendarApi = calendarRef.current.getApi();
    calendarApi.changeView('timeGridDay', info.date);
  };

  //filter
  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchSchedules(filters);
  };

  const handleClearSearch = () => {
    setFilters({
      name: "",
      status: "",
      fdate: "",
      tdate: "",
    });
    fetchSchedules(); // fetch all again
  };

  const handleAddSchedule = () => {
    navigate("/schedules/create"); // adjust route as needed
  };






  return (
    <div>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}

      <h2 className="text-xl font-bold mb-4">{t("scheduleManagement.title")}</h2>


      {/* Search Filters */}
      <form onSubmit={handleSearch} className="flex flex-col justify-between gap-4 flex-wrap bg-white p-4 rounded shadow-sm mb-4">
        <div className="flex gap-4">
          <div>
            <label>{t("scheduleManagement.fields.title")}</label>
            <input
              type="text"
              name="name"
              value={filters.name}
              onChange={handleFilterChange}
              className="border rounded p-2 w-full"
            />
          </div>

          <div>
            <label>{t("scheduleManagement.fields.status")}</label>
            <select
              name="status"
              value={filters.status}
              onChange={handleFilterChange}
              className="border rounded p-2 w-full"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="removed">Removed</option>
            </select>
          </div>

          <div>
            <label>{t("scheduleManagement.form.filters.fromDate")}</label>
            <input
              type="date"
              name="fdate"
              value={filters.fdate}
              onChange={handleFilterChange}
              className="border rounded p-2 w-full"
            />
          </div>

          <div>
            <label>{t("scheduleManagement.form.filters.toDate")}</label>
            <input
              type="date"
              name="tdate"
              value={filters.tdate}
              onChange={handleFilterChange}
              className="border rounded p-2 w-full"
            />
          </div>
        </div>

        <div className="flex justify-between">
          <div className="flex items-end gap-2">
            <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">
              {t("common.buttons.search")}
            </button>
            <button type="button" onClick={handleClearSearch} className="bg-red-600 text-white px-4 py-2 rounded">
              {t("common.buttons.clearSearch")}
            </button>
          </div>
          <div>
            <button type="button" onClick={handleAddSchedule} className="bg-blue-600 text-white px-4 py-2 rounded">
              {t("common.buttons.add")}
            </button>
          </div>
        </div>
      </form>


      <div className="bg-white shadow-md rounded-lg p-4">
        <table id="schedulesTable" className="display" style={{ width: "100%" }}></table>
      </div>


      <h2 className="mt-5">{t("scheduleManagement.calendar.fullCalendarTitle")}</h2>
      <div className="bg-white p-4 mt-4 rounded-sm shadow-sm">

        <FullCalendar
          ref={calendarRef}
          plugins={[dayGridPlugin, interactionPlugin, timeGridPlugin]}
          initialView="dayGridMonth"
          events={calendarUiEvents}
          dateClick={handleDateClick}
          headerToolbar={{
            left: 'prev,next today',
            center: 'title',
            right: 'dayGridMonth,timeGridDay',
          }}
        />
      </div>

      <h2 className="mt-5">{t("scheduleManagement.calendar.googleCalendarTitle")}</h2>
      <div className="bg-white p-4 mt-4 rounded-sm shadow-sm">
        
          <iframe
            src="https://calendar.google.com/calendar/embed?src=c8c76cfe1510e89fd1e4cc63d6dce2061d3077fc874798dd40f193b17628987f%40group.calendar.google.com&ctz=Asia/Kuala_Lumpur"
            style={{ border: 0 }}
            width="100%"
            height="600"
          />
      </div>

    </div>
  );
}

export default SchedulePage;
