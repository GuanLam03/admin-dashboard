import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import PhonelinkLockIcon from '@mui/icons-material/PhonelinkLock';

export default function SettingsPage() {
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false);
  const navigate = useNavigate();
  const { user } = useAuth();

  useEffect(() => {
    setTwoFactorEnabled(user.two_factor_enabled);
  }, [user]);

  const handleToggle = () => {
    if (!twoFactorEnabled) {
      navigate("/settings/twofactor"); // enable flow
    } else {
      navigate("/settings/twofactor/disable"); // disable flow
    }
  };

  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Settings</h2>

      <div className="flex flex-col gap-4 bg-white shadow rounded-lg p-6">
        <div className="flex items-center justify-between py-2 border-b border-gray-300">
          <div className="flex items-center gap-4">
            <PhonelinkLockIcon fontSize="large"/>
            <span className="text-gray-700">Two-Factor Authentication</span>
          </div>
          
          <label className="relative flex items-center cursor-pointer">
            <input
              type="checkbox"
              className="sr-only peer"
              checked={twoFactorEnabled}
              onChange={handleToggle}
            />
        
            <div className="w-11 h-6 bg-gray-300 rounded-full peer-checked:bg-blue-600 transition-all"></div>
            <span className="ml-3 text-sm text-gray-500">
              {twoFactorEnabled ? "On" : "Off"}
            </span>
        
            
          </label>
        </div>

        <div className="flex items-center justify-between py-2 border-b border-gray-300">

          <span className="text-gray-700">Comming Soon ...</span>
          
        </div>
        
        <div className="flex items-center justify-between py-2 border-b border-gray-300">

          <span className="text-gray-700">Comming Soon ...</span>
          
        </div>

        <div className="flex items-center justify-between py-2 border-b border-gray-300">

          <span className="text-gray-700">Comming Soon ...</span>
          
        </div>
      </div>
    </div>
  );
}
