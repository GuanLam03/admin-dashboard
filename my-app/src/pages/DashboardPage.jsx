import { useTranslation } from "react-i18next";
import { useAuth } from "../contexts/AuthContext";




function DashboardPage() {
  const {t} = useTranslation();
  const { user } = useAuth();
  return (
    <div className="w-full p-6">
      <h2 className="text-2xl font-bold mb-4">{t("home.welcome", { name: user?.name})}</h2>
    </div>
  );
}

export default DashboardPage;
