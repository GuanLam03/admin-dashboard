// EmailsLayout.jsx
import { useTranslation } from "react-i18next";
import { Outlet } from "react-router-dom";

export default function EmailsLayout() {
  const {t} = useTranslation();
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">{t("emailManagement.emails")}</h2>
      <Outlet />
    </div>
  );
}
