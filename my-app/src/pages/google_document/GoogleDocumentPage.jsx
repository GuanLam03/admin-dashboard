import { useEffect, useState } from "react";
import DataTable from "datatables.net-dt";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";

function GoogleDocumentPage() {
    const [documents, setDocuments] = useState([]);
    const [filters, setFilters] = useState({

        name: "",
        status: "active",
        fdate: "",
        tdate: ""
    });

    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const navigate = useNavigate();

    const fetchDocuments = async () => {
        try {
            const params = new URLSearchParams();
            Object.keys(filters).forEach((key) => {
                if (filters[key]) params.append(key, filters[key]);
            });

            const res = await api.get(`/google-documents?${params.toString()}`);
            setDocuments(res.data.documents || []);
        } catch (err) {
            console.error(err);
            setError(`Error fetching documents: ${err.message}`);
        }
    };

    useEffect(() => {
        fetchDocuments();
    }, []);

    useEffect(() => {
        if (!documents) return;

        const table = new DataTable("#documentsTable", {
            data: documents,
            destroy: true,
            columns: [
                {
                    data: "id",
                    title: "Action",
                    render: (id) =>
                        `<button class="edit-btn btn btn-primary" data-id="${id}">Edit</button>
             <button class="view-btn btn btn-info" data-id="${id}">View</button>
             `,
                },
                { data: "name", title: "Name" },
                {
                    data: "link",
                    title: "Link",
                    render: (link) => `<a href="${link}" target="_blank">${link}</a>`,
                },
                { data: "status", title: "Status" },
                { data: "created_at", title: "Created At" },
                { data: "updated_at", title: "Updated At" },
            ],
        });

        const handleClick = (e) => {
            if (e.target.classList.contains("edit-btn")) {
                const id = e.target.getAttribute("data-id");
                navigate(`/google-documents/edit/${id}`);
            }
            if (e.target.classList.contains("view-btn")) {
                const id = e.target.getAttribute("data-id");
                navigate(`/google-documents/view/${id}`);
            }

        };

        document.addEventListener("click", handleClick);

        return () => {
            table.destroy();
            document.removeEventListener("click", handleClick);
        };
    }, [documents, navigate]);


    const handleSearch = (e) => {
        e.preventDefault();
        fetchDocuments();
    };

    const handleClearSearch = () => {
        setFilters({
            name: "",
            status: "active",
            fdate: "",
            tdate: "",
        });
        fetchDocuments();
    };

    const handleAddGoogleDocument = () => {
        navigate('/google-documents/add');
    }

    return (
        <div>
            {error && <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>}
            {success && <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">{success}</div>}

            <h2 className="text-xl font-bold mb-4">Google Documents</h2>

            {/* Search Filters */}
            <form onSubmit={handleSearch} className="flex flex-col justify-between gap-4 flex-wrap bg-white p-4 rounded shadow-sm mb-4">
                <div className="flex gap-4">
                    {/* <div>
                        <label>ID</label>
                        <input type="text" value={filters.id} onChange={(e) => setFilters({ ...filters, id: e.target.value })} className="border rounded p-2 w-full" />
                    </div> */}
                    <div>
                        <label>Name</label>
                        <input type="text" value={filters.name} onChange={(e) => setFilters({ ...filters, name: e.target.value })} className="border rounded p-2 w-full" />
                    </div>
                    <div>
                        <label>Status</label>
                        <select
                            name="status"
                            value={filters.status}
                            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
                            className="border rounded p-2 w-full"
                        >
                            <option value="active">Active</option>
                            <option value="inactive">Inactive</option>
                            <option value="removed">Removed</option>
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
                        <button type="button" onClick={handleAddGoogleDocument} className="bg-blue-600 text-white px-4 py-2 rounded">Add Google Document</button>
                    </div>
                </div>

            </form>

            {/* DataTable */}
            <div className="bg-white shadow-md rounded-lg p-4">
                <table id="documentsTable" className="display" style={{ width: "100%" }}></table>
            </div>
        </div>
    );
}

export default GoogleDocumentPage;
