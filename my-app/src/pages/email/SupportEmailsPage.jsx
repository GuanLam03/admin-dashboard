import { useEffect, useState } from "react";
import api from "../../api/axios";
import InboxPage from "./InboxPage";

export default function SupportEmailsPage() {
  const [account, setAccount] = useState(null);

  useEffect(() => {
    api.get("/gmail/accounts").then((res) => {
      // find account with label = "technical"
      const acc = res.data.find((a) => a.team === "support");
      setAccount(acc);
    });
  }, []);

  if (!account) {
    return <p>No support account connected yet.</p>;
  }

  return (
    <>
      <h5>Support Emails ({account.email}):</h5>
      <InboxPage folder={account.team} emailAddress={account.email} />
    </>
  );
}
