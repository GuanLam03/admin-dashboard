import { useState } from "react";
import api from "../../api/axios";

function ReplyBox({ mailbox,messageId, email, onSent }) {
  const [body, setBody] = useState("");
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");

  const handleSend = async () => {
    if (!body.trim()) return;
    setSending(true);
    setError("");

    try {
      await api.post(`/gmail/${mailbox}/messages/${messageId}/reply`, {
        email,
        body,
      });
      setBody("");
      onSent(); // refresh thread
    } catch (err) {
      setError("Failed to send reply");
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="mt-4">
      <textarea
        className="w-full border rounded-md p-2"
        rows={4}
        placeholder="Write your reply..."
        value={body}
        onChange={(e) => setBody(e.target.value)}
      />
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
      <button
        className="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        onClick={handleSend}
        disabled={sending}
      >
        {sending ? "Sending..." : "Send Reply"}
      </button>
    </div>
  );
}

export default ReplyBox;
