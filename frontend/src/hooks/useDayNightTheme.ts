import { useEffect } from "react";

export default function useDayNightTheme() {
  useEffect(() => {
    function updateTheme() {
      const hour = new Date().getHours();

      // Daytime: 6 AM to 6 PM
      const isDaytime = hour >= 6 && hour < 18;

      document.body.classList.toggle(
        "day-theme",
        isDaytime
      );

      document.body.classList.toggle(
        "night-theme",
        !isDaytime
      );
    }

    updateTheme();

    // Check again every minute
    const interval = window.setInterval(
      updateTheme,
      60 * 1000
    );

    return () => {
      window.clearInterval(interval);
    };
  }, []);
}