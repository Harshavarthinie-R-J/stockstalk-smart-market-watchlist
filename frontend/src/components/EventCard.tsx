import { Activity, ArrowDown, ArrowUp, Clock, Volume2 } from "lucide-react";
import AttentionBadge from "./AttentionBadge";
import DataStatus from "./DataStatus";

interface EventCardProps {
  symbol: string;
  name?: string | null;
  priceChangePercent?: number | null;
  currentPrice?: number | null;
  volumeRatio?: number | null;
  priceAnomaly?: boolean | null;
  volumeAnomaly?: boolean | null;
  attentionScore?: number | null;
  context?: string | null;
  confidence?: string | null;
  detectedAt?: string | null;
  reliabilityState?: string | null;
  onClick?: () => void;
}

export default function EventCard({
  symbol,
  name,
  priceChangePercent,
  currentPrice,
  volumeRatio,
  priceAnomaly,
  volumeAnomaly,
  attentionScore,
  context,
  confidence,
  detectedAt,
  reliabilityState,
  onClick,
}: EventCardProps) {
  const change = priceChangePercent ?? 0;
  const isPositive = change >= 0;

  return (
    <button type="button" className="event-card" onClick={onClick}>
      <div className="event-card-top">
        <div className="event-stock-info">
          <div className="event-symbol">{symbol}</div>
          {name && <div className="event-company">{name}</div>}
        </div>
        <AttentionBadge score={attentionScore} />
      </div>

      <div className="event-price-row">
        <div className="event-price">
          {currentPrice !== null && currentPrice !== undefined
            ? `₹${currentPrice.toLocaleString("en-IN", {
                maximumFractionDigits: 2,
              })}`
            : "—"}
        </div>

        <div
          className={`event-change ${
            isPositive ? "event-change-positive" : "event-change-negative"
          }`}
        >
          {isPositive ? <ArrowUp size={13} /> : <ArrowDown size={13} />}
          {Math.abs(change).toFixed(2)}%
        </div>
      </div>

      <div className="event-signals">
        {priceAnomaly && (
          <span className="event-signal">
            <Activity size={12} />
            Unusual price move
          </span>
        )}

        {volumeAnomaly && (
          <span className="event-signal">
            <Volume2 size={12} />
            Volume spike
          </span>
        )}

        {!priceAnomaly && !volumeAnomaly && (
          <span className="event-signal">
            <Activity size={12} />
            Meaningful trend
          </span>
        )}

        {volumeRatio !== null &&
          volumeRatio !== undefined &&
          volumeRatio > 1 && (
            <span className="event-signal">
              {volumeRatio.toFixed(1)}× avg volume
            </span>
          )}
      </div>

      {context && (
        <div className="event-context">
          <strong>Context</strong>
          <span>{context}</span>
        </div>
      )}

      <div className="event-card-footer">
        <div className="event-confidence">
          Confidence: <strong>{confidence || "Neutral"}</strong>
        </div>

        <div className="event-detected">
          <Clock size={12} />
          {detectedAt ? formatTime(detectedAt) : "Recently"}
        </div>
      </div>

      <div className="event-status">
        <DataStatus status={reliabilityState} timestamp={detectedAt} />
      </div>
    </button>
  );
}

function formatTime(timestamp: string) {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) return "Recently";
  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}