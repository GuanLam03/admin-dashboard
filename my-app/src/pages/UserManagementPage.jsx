import { useEffect, useRef, useState } from "react";
import DataTable from "datatables.net-dt";
import api from "../api/axios";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import dataTableLocales from "../utils/i18n/datatableLocales";
import { useCentrifuge } from "../contexts/CentrifugeContext";

function UserManagementPage() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const tableRef = useRef(null);
  const [error, setError] = useState("");

  // global online user list
  const { onlineUsers } = useCentrifuge();
 
  /** ------------------------------------------
   * 1️⃣ Initialize DataTable immediately (empty)
   * -------------------------------------------*/
  useEffect(() => {
    const dt = new DataTable("#usersTable", {
      data: [],               // initially empty
      destroy: true,
      language: dataTableLocales[i18n.language],
      columns: [
        { data: "role", title: t("userManagement.fields.role") },
        { data: "name", title: t("userManagement.fields.name") },
        { data: "email", title: t("userManagement.fields.email") },
        { data: "created_at", title: t("userManagement.fields.createdAt") },
        {
          data: "online",
          title: t("userManagement.status"),
          render: (online) =>
            online
              ? '<span class="inline-block w-3 h-3 bg-green-500 rounded-full"></span>'
              : '<span class="inline-block w-3 h-3 bg-red-500 rounded-full"></span>',
        },
        {
          data: "id",
          title: t("common.labels.action"),
          render: (id) =>
            `<button class="edit-btn btn btn-primary" data-id="${id}">${t(
              "common.buttons.edit"
            )}</button>`,
        },
      ],
      responsive: true,
          createdRow: (row, rowData) => {
              const btn = row.querySelector(".edit-btn");
              btn.addEventListener("click", () => {
              navigate(`/user-management/edit/${rowData.id}`);
              });
          },
    });

    tableRef.current = dt;

    return () => dt.destroy();
  }, [i18n.language, t]);

  /** ------------------------------------------
   * 2️⃣ Fetch users and fill table
   * -------------------------------------------*/
  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await api.get("/users");
        const data = res.data.users.map((u) => ({
          id: u.id,
          role: u.role || "",
          name: u.name,
          email: u.email,
          created_at: u.created_at,
          online: onlineUsers.includes(String(u.id)),
        }));

        const dt = tableRef.current;
        dt.clear();
        dt.rows.add(data);
        dt.draw();
      } catch (err) {
        setError(`Error fetching users: ${err.message}`);
      }
    };

    fetchUsers();
  }, [onlineUsers]);

  /** ------------------------------------------
   * 3️⃣ Update online status live
   * -------------------------------------------*/
  useEffect(() => {
    if (!tableRef.current) return;

    const dt = tableRef.current;
    dt.rows().every(function () {
      const rowData = this.data();
      rowData.online = onlineUsers.includes(String(rowData.id));
      this.data(rowData, false);
    });

    dt.rows().invalidate();
    dt.draw(false);
  }, [onlineUsers]);

  return (
    <>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}

      <div className="w-full">
        <h2 className="text-2xl font-bold mb-4">{t("userManagement.userManagement")}</h2>

        <div className="bg-white shadow-md rounded-lg p-4">
          <table id="usersTable" className="display" style={{ width: "100%" }}></table>
        </div>
      </div>
    </>
  );
}

export default UserManagementPage;
