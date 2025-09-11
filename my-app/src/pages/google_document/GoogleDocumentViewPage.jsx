import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../api/axios";

function GoogleDocumentViewPage() {
  const { id } = useParams();
  const [document, setDocument] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDocument = async () => {
      try {
        const res = await api.get(`/google-documents/${id}`);
        setDocument(res.data.document);
      } catch (err) {
        console.error(err);
        setError("Failed to load document.");
      } finally {
        setLoading(false);
      }
    };
    fetchDocument();
  }, [id]);

  if (loading) {
    return <p>Loading document...</p>;
  }

  if (error) {
    return <div className="bg-red-100 text-red-600 px-4 py-2 rounded">{error}</div>;
  }

  if (!document) {
    return <p>No document found.</p>;
  }

  return (
    <div className="h-full">
      <h2 className="text-xl font-bold mb-4">{document.name}</h2>
      <iframe
        src={document.link}
        title={document.name}
        className="w-full h-[95%] border rounded"
      />
    </div>
  );
}

export default GoogleDocumentViewPage;
