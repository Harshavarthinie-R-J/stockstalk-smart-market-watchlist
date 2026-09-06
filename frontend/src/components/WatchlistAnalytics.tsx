import {
  Activity,
  ArrowDown,
  ArrowUp,
  BarChart3,
  Flame,
} from "lucide-react";

import AttentionBadge from "./AttentionBadge";

interface AnalyticsStock {
  symbol: string;
  changePercent?: number | null;
  volume?: number | null;
}

interface WatchlistAnalyticsProps {
  stocks: AnalyticsStock[];
  events: {
    symbol: string;
    changePercent?: number | null;
    volumeRatio?: number | null;
    attentionScore?: number | null;
  }[];
}

export default function WatchlistAnalytics({
  stocks,
  events,
}: WatchlistAnalyticsProps) {
  const validStocks = stocks.filter(
    (stock) =>
      stock.changePercent != null
  );

  const averageMovement =
    validStocks.length > 0
      ? validStocks.reduce(
          (sum, stock) =>
            sum +
            (stock.changePercent ?? 0),
          0
        ) / validStocks.length
      : 0;

  const topGainer =
    [...validStocks].sort(
      (a, b) =>
        (b.changePercent ?? 0) -
        (a.changePercent ?? 0)
    )[0];

  const topLoser =
    [...validStocks].sort(
      (a, b) =>
        (a.changePercent ?? 0) -
        (b.changePercent ?? 0)
    )[0];

  const highestAttention =
    [...events].sort(
      (a, b) =>
        (b.attentionScore ?? 0) -
        (a.attentionScore ?? 0)
    )[0];

  const highestVolume =
    [...events].sort(
      (a, b) =>
        (b.volumeRatio ?? 0) -
        (a.volumeRatio ?? 0)
    )[0];

  return (
    <section className="watchlist-analytics">

      <div className="watchlist-analytics-header">
        <div>
          <div className="section-kicker">
            WATCHLIST ANALYTICS
          </div>

          <h2>
            Market at a glance
          </h2>

          <p>
            A quick view of what is moving
            across your watchlist.
          </p>
        </div>

        <div className="analytics-icon">
          <BarChart3 size={20} />
        </div>
      </div>

      <div className="analytics-grid">

        <div className="analytics-card">
          <div className="analytics-card-icon">
            <Activity size={17} />
          </div>

          <span>
            Average movement
          </span>

          <strong
            className={
              averageMovement >= 0
                ? "positive-value"
                : "negative-value"
            }
          >
            {averageMovement >= 0
              ? "+"
              : ""}
            {averageMovement.toFixed(2)}%
          </strong>
        </div>

        <div className="analytics-card">
          <div className="analytics-card-icon">
            <ArrowUp size={17} />
          </div>

          <span>
            Top mover
          </span>

          <strong>
            {topGainer?.symbol ?? "—"}
          </strong>

          {topGainer && (
            <small
              className={
                (topGainer.changePercent ??
                  0) >= 0
                  ? "positive-value"
                  : "negative-value"
              }
            >
              {(topGainer.changePercent ??
                0) >= 0
                ? "+"
                : ""}
              {(
                topGainer.changePercent ??
                0
              ).toFixed(2)}
              %
            </small>
          )}
        </div>

        <div className="analytics-card">
          <div className="analytics-card-icon">
            <ArrowDown size={17} />
          </div>

          <span>
            Largest decline
          </span>

          <strong>
            {topLoser?.symbol ?? "—"}
          </strong>

          {topLoser && (
            <small className="negative-value">
              {(topLoser.changePercent ??
                0) >= 0
                ? "+"
                : ""}
              {(
                topLoser.changePercent ??
                0
              ).toFixed(2)}
              %
            </small>
          )}
        </div>

        <div className="analytics-card">
          <div className="analytics-card-icon">
            <Flame size={17} />
          </div>

          <span>
            Highest attention
          </span>

          <strong>
            {highestAttention?.symbol ??
              "—"}
          </strong>

          {highestAttention && (
            <AttentionBadge
              score={
                highestAttention.attentionScore
              }
            />
          )}
        </div>

      </div>

      <div className="analytics-highlight">

        <div>
          <strong>
            {events.length}
          </strong>

          <span>
            detected events
          </span>
        </div>

        <div>
          <strong>
            {events.filter(
              (event) =>
                (event.volumeRatio ??
                  0) >= 1.5
            ).length}
          </strong>

          <span>
            unusual-volume events
          </span>
        </div>

        <div>
          <strong>
            {highestVolume?.symbol ??
              "—"}
          </strong>

          <span>
            highest volume anomaly
          </span>
        </div>

      </div>

    </section>
  );
}