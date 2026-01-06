import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const Logout = () => {
  const navigate = useNavigate();

  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [showToast, setShowToast] = useState(false);

  const sendLogoutRequest = async () => {
    try {
      const res = await fetch(`${baseUrl}/api/v1/logout`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json();
      
      if (data.error) {
        setSnackbarMessage(data.error);
        setShowToast(true);
      }

      if (data.ok) {
        // Clear session storage
        sessionStorage.removeItem("token");
        sessionStorage.removeItem("baseUrl");

        // Or clear ALL session storage:
        // sessionStorage.clear();

        // Redirect to login page
        navigate("/login");
      }
    } catch (err) {
      setSnackbarMessage(err.messaage);
      setShowToast(true);
    }

  };

  useEffect(() => {
    sendLogoutRequest();
    // Redirect to login page
    navigate("/login");
  }, []);

  return null; // No UI needed
};

export default Logout;
