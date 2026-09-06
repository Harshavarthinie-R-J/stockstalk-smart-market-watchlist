import React from "react";

interface DataStatusProps {
  status?: string | null;
  timestamp?: string | null;
}

export default function DataStatus({ status, timestamp }: DataStatusProps) {
  const normalizedStatus = status?.toUpperCase() || "UNAVAILABLE";

  const labels: Record<string, string> = {
    LIVE: "Live",
    DELAYED: "Delayed",
    STALE: "Stale",
    UNAVAILABLE: "Unavailable",
    CONFLICTING: "Conflicting",
  };

  const label = labels[normalizedStatus] || "Unknown";

  return (
    <div className="data-status">
      <span
        className={`data-status-dot status-${normalizedStatus.toLowerCase()}`}
      />
      <span className="data-status-label">{label}</span>
      {timestamp && (
        <span className="data-status-time">{formatTime(timestamp)}</span>
      )}
    </div>
  );
}

function formatTime(timestamp: string) {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) {
    return "";
  }
  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}