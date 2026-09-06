import {
  Activity,
  ArrowDown,
  ArrowUp,
  CheckCircle2,
  Clock,
  ExternalLink,
  Info,
  TrendingDown,
  TrendingUp,
  Volume2,
  X,
} from "lucide-react";

import AttentionBadge from "./AttentionBadge";
import DataStatus from "./DataStatus";

interface EventDetailProps {
  event: {
    symbol: string;
    name?: string | null;
    previousPrice?: number | null;
    currentPrice?: number | null;
    priceChange?: number | null;
    priceChangePercent?: number | null;
    previousVolume?: number | null;
    currentVolume?: number | null;
    averageVolume?: number | null;
    volumeRatio?: number | null;
    volumeAnomaly?: boolean | null;
    priceAnomaly?: boolean | null;
    near52WeekHigh?: boolean | null;
    near52WeekLow?: boolean | null;
    signals?: string[] | null;
    context?: string | null;
    confidence?: string | null;
    contextType?: string | null;
    attribution?: string | null;
    attributionConfidence?: string | null;
    attributionSource?: string | null;
    attentionScore?: number | null;
    detectedAt?: string | null;
    reliabilityState?: string | null;
  };
  onClose: () => void;
}

export default function EventDetail({ event, onClose }: EventDetailProps) {
  const change = event.priceChangePercent ?? 0;
  const positive = change >= 0;

  return (
    <div className="event-detail-overlay" onClick={onClose}>
      <div
        className="event-detail-modal"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="event-detail-header">
          <div>
            <div className="event-detail-symbol">{event.symbol}</div>
            {event.name && (
              <div className="event-detail-name">{event.name}</div>
            )}
          </div>

          <button
            type="button"
            className="event-detail-close"
            onClick={onClose}
            aria-label="Close"
          >
            <X size={18} />
          </button>
        </div>

        {/* Price & Movement Overview */}
        <div className="event-detail-movement">
          <div>
            <div className="event-detail-label">Observed Price</div>
            <div className="event-detail-price">
              {formatPrice(event.currentPrice)}
            </div>
          </div>

          <div
            className={`event-detail-change ${
              positive ? "event-change-positive" : "event-change-negative"
            }`}
          >
            {positive ? <ArrowUp size={16} /> : <ArrowDown size={16} />}
            {Math.abs(change).toFixed(2)}%
          </div>

          <AttentionBadge score={event.attentionScore} />
        </div>

        {/* Metrics Grid */}
        <div className="event-detail-section">
          <SectionTitle
            icon={<TrendingUp size={16} />}
            title="Price & Movement Metrics"
          />
          <div className="event-metric-grid">
            <Metric
              label="Previous Price"
              value={formatPrice(event.previousPrice)}
            />
            <Metric
              label="Current Price"
              value={formatPrice(event.currentPrice)}
            />
            <Metric
              label="Calculated Change"
              value={`${change >= 0 ? "+" : ""}${change.toFixed(2)}%`}
            />
          </div>
        </div>

        {/* Signals */}
        <div className="event-detail-section">
          <SectionTitle
            icon={<Activity size={16} />}
            title="Trigger Factors"
          />

          <div className="event-signal-list">
            {event.priceAnomaly && (
              <Signal
                icon={<TrendingUp size={14} />}
                title="Unusual Price Velocity"
                description="Movement surpassed 2-sigma variance relative to historical movement."
              />
            )}

            {event.volumeAnomaly && (
              <Signal
                icon={<Volume2 size={14} />}
                title="Unusual Trading Volume"
                description={
                  event.volumeRatio
                    ? `Volume is ${event.volumeRatio.toFixed(1)}× higher than standard daily volume.`
                    : "Turnover differs significantly from historical benchmark."
                }
              />
            )}

            {event.near52WeekHigh && (
              <Signal
                icon={<TrendingUp size={14} />}
                title="Near 52-Week High"
                description="Stock is trading within 2% of annual ceiling."
              />
            )}

            {event.near52WeekLow && (
              <Signal
                icon={<TrendingDown size={14} />}
                title="Near 52-Week Low"
                description="Stock is trading within 2% of annual floor."
              />
            )}

            {event.signals &&
              event.signals.map((sig, idx) => (
                <Signal
                  key={`${sig}-${idx}`}
                  icon={<Activity size={14} />}
                  title={sig}
                />
              ))}
          </div>
        </div>

        {/* Volume */}
        {(event.currentVolume != null || event.averageVolume != null) && (
          <div className="event-detail-section">
            <SectionTitle icon={<Volume2 size={16} />} title="Volume Analysis" />
            <div className="event-metric-grid">
              <Metric
                label="Current Volume"
                value={formatNumber(event.currentVolume)}
              />
              <Metric
                label="Average Volume"
                value={formatNumber(event.averageVolume)}
              />
              <Metric
                label="Volume Ratio"
                value={
                  event.volumeRatio ? `${event.volumeRatio.toFixed(2)}×` : "—"
                }
              />
            </div>
          </div>
        )}

        {/* Context */}
        <div className="event-detail-section">
          <SectionTitle icon={<Info size={16} />} title="Context & Attribution" />
          {event.context ? (
            <div className="event-context-box">
              <div className="event-context-main">{event.context}</div>
              {event.contextType && (
                <div className="event-context-type">{event.contextType}</div>
              )}
            </div>
          ) : (
            <div className="event-no-context">
              <Info size={14} />
              <span>No direct news catalyst identified for this session.</span>
            </div>
          )}
        </div>

        {/* Confidence & Status */}
        <div className="event-detail-section">
          <SectionTitle icon={<CheckCircle2 size={16} />} title="Data Reliability" />
          <DataStatus
            status={event.reliabilityState}
            timestamp={event.detectedAt}
          />
        </div>

        {/* Footer */}
        <div className="event-detail-footer">
          <div>
            Detected{" "}
            {event.detectedAt
              ? formatDateTime(event.detectedAt)
              : "recently"}
          </div>
          <button
            type="button"
            className="event-detail-done"
            onClick={onClose}
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
}

function SectionTitle({ icon, title }: { icon: React.ReactNode; title: string }) {
  return (
    <div className="event-section-title">
      {icon}
      <span>{title}</span>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="event-metric">
      <div className="event-metric-label">{label}</div>
      <div className="event-metric-value">{value}</div>
    </div>
  );
}

function Signal({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description?: string;
}) {
  return (
    <div className="event-signal-detail">
      <div className="event-signal-icon">{icon}</div>
      <div>
        <div className="event-signal-detail-title">{title}</div>
        {description && (
          <div className="event-signal-detail-description">{description}</div>
        )}
      </div>
    </div>
  );
}

function formatPrice(value?: number | null) {
  if (value == null) return "—";
  return `₹${value.toLocaleString("en-IN", { maximumFractionDigits: 2 })}`;
}

function formatNumber(value?: number | null) {
  if (value == null) return "—";
  return value.toLocaleString("en-IN");
}

function formatDateTime(timestamp: string) {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) return "recently";
  return date.toLocaleString("en-IN", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}