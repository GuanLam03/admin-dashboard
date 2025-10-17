import { useState, useEffect, useRef } from "react";
import api from "../../api/axios";
import { HiOutlineTemplate } from "react-icons/hi";
import { Tooltip } from "bootstrap"; 
import CloseIcon from '@mui/icons-material/Close';
import Editor from "../../tools/Editor";


// const templates = [
//   { id: 1, label: "Forward to Support", content: "Hi Support,\n\nPlease check the forwarded email below. Thanks!\n\n" },
//   { id: 2, label: "Customer Apology", content: "Dear Customer,\n\nWe apologize for the inconvenience caused. We are checking this.\n\n" },
//   { id: 3, label: "Forward to Technical", content: "Hi Technical Team,\n\nPlease investigate this forwarded message. Details below.\n\n" },
// ];



export default function MessageActionBox({ mailbox, messageId, email, onSent, type, originalMessage }) {
  const [body, setBody] = useState(""); 
  const [to, setTo] = useState(""); 
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");
  const [forwardedHTML, setForwardedHTML] = useState("");

  const [showTemplateModal, setShowTemplateModal] = useState(false);

  const textareaRef = useRef(null);
  const [templates, setTemplates] = useState([]);

  useEffect(() => {
      if (mailbox) {
        fetchTemplates();
      }
  }, [mailbox]);

    
  useEffect(() => {
    if (type === "forward" && originalMessage) {
      // Escape < > to read <xxx@gmail.com>
      const safeFrom = originalMessage.from
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");

      // Build forwarded block with headers + original HTML body
      const forwarded = `
        <div style="margin-top:16px; padding-top:8px;">
          <p>--- Forwarded message ---</p>
          <p><b>From:</b> ${safeFrom}</p>
          <p><b>Date:</b> ${originalMessage.date}</p>
          <p><b>Subject:</b> ${originalMessage.subject}</p>
          ${originalMessage.body}
        </div>
      `;
      setForwardedHTML(forwarded);
    }
  }, [type, originalMessage]);

  //tooltip initialization bootstrap
  useEffect(() => {
    const triggers = document.querySelectorAll('[data-bs-toggle="tooltip"]');
    const instances = Array.from(triggers).map(el => {
      // disable "click" trigger, keep only hover + focus
      return new Tooltip(el, { trigger: "hover focus" });
    });

    return () => {
      instances.forEach(i => i.dispose());
    };
  }, []);

  function adjustHeight(){
    if (textareaRef.current) {
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  };


  async function handleSend(){
    if (!body.trim() || (type === "forward" && !to.trim())) return;
    setSending(true);
    setError("");

    try {
      if (type === "reply") {
        // Reply needs messageId for threading
        await api.post(`/gmail/${mailbox}/messages/${messageId}/reply`, {
          email,
          body: body + forwardedHTML,
        });
      } else if (type === "forward") {
        // Forward is a brand new message, no messageId
        await api.post(`/gmail/${mailbox}/messages/forward`, {
          email,
          to,
          body: body + forwardedHTML,
          subject: `Fwd: ${originalMessage.subject}`,
        });
      }

      setBody("");
      setTo("");
      setForwardedHTML("");
      onSent();
    } catch (err) {
      setError("Failed to send message");
    } finally {
      setSending(false);
    }
  };


  function hideAllTooltips(){
    document.querySelectorAll('[data-bs-toggle="tooltip"]').forEach(el => {
      const inst = Tooltip.getInstance(el);
      if (inst) inst.hide();
    });
  };
  

  function handleTemplate(){
    hideAllTooltips();
    setShowTemplateModal(true);
  };


  async function fetchTemplates() {
      try {
        const res = await api.get(`/gmail/templates?team=${mailbox}`);
        setTemplates(res.data); // backend returns templates for that team
      } catch (err) {
        console.error("Failed to load templates", err);
      }
    }

  
     

  return (
    <div className="border px-2 py-4 rounded-lg shadow-sm">
      {type === "forward" && (
        <div className="mb-4 flex items-center gap-2 sticky top-0 bg-white">
          <label htmlFor="to" className="block text-sm font-semibold">To:</label>
          <input
            id="to"
            type="email"
            className="w-full rounded-md p-2 focus:outline-none"
            placeholder="Enter recipient's email"
            value={to}
            onChange={(e) => setTo(e.target.value)}
          />
        </div>
      )}

      {/* User’s text */}
      {/* <textarea
        ref={textareaRef}
        className="w-full resize-none rounded-md p-2 focus:outline-none overflow-hidden"
        placeholder={type === "forward" ? "Write your message..." : "Write your reply..."}
        value={body}
        onChange={(e) => {
          setBody(e.target.value);
          adjustHeight();
        }}
        style={{ minHeight: "100px" }}
      /> */}

      <Editor
        value={body}
        onChange={setBody}
      />

      {/* Forwarded HTML preview */}
      {type === "forward" && forwardedHTML && (
        <div
          className="mt-4 pt-2 text-sm text-gray-800"
          dangerouslySetInnerHTML={{ __html: forwardedHTML }}
        />
      )}

      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}

      <div className="flex items-center gap-2 mt-4">
        <button
          className="px-4 py-2 bg-blue-600 text-white !rounded-full hover:bg-blue-700"
          onClick={handleSend}
          disabled={sending}
        >
          {sending ? "Sending..." : type === "forward" ? "Send Forward" : "Send Reply"}
        </button>
        <button
          type="button"
          className="btn "
          data-bs-toggle="tooltip"
          data-bs-placement="bottom"
          data-bs-title="Template"
          onClick={
            handleTemplate
          }

        >
          <HiOutlineTemplate size={20}/>
        </button>

       
      </div>

      {showTemplateModal && (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg shadow-lg w-96 p-4">
          <div className="flex justify-between items-center mb-2">
            <h3 className="text-lg font-semibold">Choose a Template</h3>
            <button onClick={() => setShowTemplateModal(false)}><CloseIcon /></button>
          </div>
          
          <div className="space-y-2">
            {templates.map((t) => (
              <div key={t.id}>
                <button
                  className="w-full text-left p-2 border rounded hover:bg-gray-100"
                  onClick={() => {
                    setBody(t.content);  // apply template to textarea
                    setShowTemplateModal(false); // close modal
                  }}
                >
                  {t.name}
                </button>
              </div>
            ))}
          </div>
          <button
            className="mt-4 px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
            onClick={() => setShowTemplateModal(false)}
          >
            Close
          </button>
        </div>
      </div>
    )}


      
    
    </div>
  );
}
