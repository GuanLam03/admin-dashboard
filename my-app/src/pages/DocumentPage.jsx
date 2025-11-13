import { useEffect, useState } from "react";
import DataTable from "datatables.net-dt";
import api from "../api/axios";
import { useTranslation } from "react-i18next";
import dataTableLocales from "../utils/i18n/datatableLocales";

function DocumentPage() {
  const [selectedFiles, setSelectedFiles] = useState([]);
  const [documents, setDocuments] = useState([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const {t,i18n} = useTranslation();

  // Fetch documents
  const fetchDocuments = async () => {
    try {
      const res = await api.get("/documents");
      setDocuments(
        res.data.documents.map((doc) => ({
          id: doc.id,
          name: doc.Filename,
          path: doc.path,
          created_at: doc.created_at,
        }))
      );
    } catch (err) {
      console.error(err);
      setError(`Error fetching documents: ${err.message}`);
    }
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  // Handle file selection
  const handleFileChange = (e) => {
    setSelectedFiles([...e.target.files]);
  };

  // Upload files
  const handleUpload = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    if (selectedFiles.length === 0) {
      alert("Please select files");
      return;
    }

    try {
      const formData = new FormData();
      selectedFiles.forEach((file) => formData.append("files", file));

      // console.log("uploading: ",formData)
      // console.log("selected file: ",selectedFiles)

      await api.post("/documents/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        timeout: 10000,
      });

      setSuccess("Files uploaded successfully!");
      setSelectedFiles([]);
      fetchDocuments(); // reload
    } catch (err) {
      console.error("Upload error", err);
      setError(`Error uploading: ${err.message}`);
    }
  };

  // Setup DataTable when documents change
  useEffect(() => {
    if (!documents) return;

    const table = new DataTable("#documentsTable", {
      data: documents,
      destroy: true,
      language: dataTableLocales[i18n.language],
      
      columns: [
        { data: "name", title: t("documentManagement.fields.filename"), className: "dt-left" },
        { data: "created_at", title: t("documentManagement.fields.createdAt"), className: "dt-left" },
        {
          data: "id",
          title: t("common.labels.action"),
          render: (id, type, row) => `
            <button 
              class="download-btn bg-white-600 text-blue-600 px-3 py-2 rounded text-sm border-1 hover:opacity-75" 
              data-id="${id}" 
              data-name="${row.name}">
              ${t("common.buttons.download")}
            </button>
          `,
        },
      ],
    });

    return () => {
      table.destroy();
    };
  }, [documents]);

  // Global click listener (attach once)
  useEffect(() => {
    const handleClick = (e) => {
      if (e.target.classList.contains("download-btn")) {
        const id = e.target.getAttribute("data-id");
        const name = e.target.getAttribute("data-name");
        handleDownload(id, name);
      }
    };

    document.addEventListener("click", handleClick);
    return () => document.removeEventListener("click", handleClick);
  }, []);

  // File download
  const handleDownload = async (id, name) => {
    try {
      const res = await api.get(`/documents/download/${id}`, {
        responseType: "blob",
      });

      console.log(res);
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", name);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (err) {
      console.error("Download error", err);
      alert("Error downloading file");
    }
  };

  return (
    <div>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">
          {success}
        </div>
      )}

      <h2 className="text-xl font-bold mb-4">{t("documentManagement.title")}</h2>

      {/* Upload Form */}
      <form
        onSubmit={handleUpload}
        className="flex items-center gap-4 mb-4 justify-between bg-white p-4 rounded shadow-sm"
      >
        <div className="flex flex-col">
          <label>{t("documentManagement.form.selectFiles")}:</label>
          <input
            type="file"
            multiple
            onChange={handleFileChange}
            className="border rounded p-2"
          />
        </div>

        <button
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          {t("common.buttons.upload")}
        </button>
      </form>

      {/* Document List Table */}
      <div className="bg-white shadow-md rounded-lg p-4">
        <table
          id="documentsTable"
          className="display"
          style={{ width: "100%" }}
        ></table>
      </div>
    </div>
  );
}

export default DocumentPage;
