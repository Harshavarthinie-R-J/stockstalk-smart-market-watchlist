import { gql } from "@apollo/client";
import { useQuery } from "@apollo/client/react";
import {
  Activity,
  AlertCircle,
  CheckCircle2,
  Clock3,
  Flame,
  RefreshCw,
  TrendingUp,
  Volume2,
  Zap,
} from "lucide-react";
import { useMemo, useState } from "react";

import AttentionBadge from "./AttentionBadge";
import DataStatus from "./DataStatus";
import EventDetail from "./EventDetail";

const SINCE_LAST_CHECKED = gql`
  query SinceLastChecked($watchlistId: ID!) {
    changesSinceLastVisit(watchlistId: $watchlistId) {
      currentTime
      meaningfulChanges
      notableChanges
      totalStocks
      unchanged
      previousCheckpoint

      changes {
        id
        instrumentId
        symbol
        type
        severity
        currentValue
        changePercent
        currentVolume
        averageVolume
        volumeRatio
        volumeAnomaly
        priceAnomaly
        near52WeekHigh
        near52WeekLow
        signals
        context
        confidence
        contextType
        attribution
        attributionConfidence
        attributionSource
        attentionScore
        detectedAt
        reliabilityState
      }
    }
  }
`;

interface ChangeData {
  id?: string | null;
  instrumentId: string;
  symbol: string;

  type?: string | null;
  severity?: string | null;

  currentValue?: number | null;
  changePercent?: number | null;

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
}

interface SummaryData {
  currentTime?: string | null;

  meaningfulChanges: number;
  notableChanges: number;

  totalStocks: number;
  unchanged: number;

  previousCheckpoint?: string | null;

  changes: ChangeData[];
}

interface QueryData {
  changesSinceLastVisit: SummaryData | null;
}

interface SinceLastCheckedProps {
  watchlistId: string;
}

type EventFilter =
  | "ALL"
  | "MEANINGFUL"
  | "NOTABLE"
  | "HIGH_ATTENTION";

export default function SinceLastChecked({
  watchlistId,
}: SinceLastCheckedProps) {
  const [selectedEvent, setSelectedEvent] =
    useState<ChangeData | null>(null);

  const [activeFilter, setActiveFilter] =
    useState<EventFilter>("ALL");

  const { data, loading, error, refetch } =
    useQuery<QueryData>(SINCE_LAST_CHECKED, {
      variables: {
        watchlistId,
      },
      fetchPolicy: "network-only",
      skip: !watchlistId,
    });

  const summary = data?.changesSinceLastVisit;

  /*
   * Always keep the highest-attention events first.
   * The score comes from the backend.
   */
  const sortedChanges = useMemo(() => {
    const changes = summary?.changes ?? [];

    return [...changes].sort(
      (a, b) =>
        (b.attentionScore ?? 0) -
        (a.attentionScore ?? 0)
    );
  }, [summary?.changes]);

  /*
   * Top 3 events that deserve attention first.
   */
  const topEvents = useMemo(() => {
    return sortedChanges.slice(0, 3);
  }, [sortedChanges]);

  /*
   * Apply the selected event filter.
   */
  const filteredChanges = useMemo(() => {
    switch (activeFilter) {
      case "MEANINGFUL":
        return sortedChanges.filter(
          (event) =>
            event.severity?.toUpperCase() ===
            "SIGNIFICANT"
        );

      case "NOTABLE":
        return sortedChanges.filter(
          (event) =>
            event.severity?.toUpperCase() ===
            "NOTABLE"
        );

      case "HIGH_ATTENTION":
        return sortedChanges.filter(
          (event) =>
            (event.attentionScore ?? 0) >= 75
        );

      default:
        return sortedChanges;
    }
  }, [sortedChanges, activeFilter]);

  async function handleRefresh() {
    try {
      await refetch();
    } catch (err) {
      console.error(
        "Unable to refresh changes:",
        err
      );
    }
  }

  function openEvent(event: ChangeData) {
    setSelectedEvent(event);
  }

  if (loading) {
    return (
      <section className="since-checked-section">
        <div className="since-checked-header">
          <div>
            <div className="section-kicker">
              SINCE YOU LAST CHECKED
            </div>

            <h2>What changed?</h2>
          </div>

          <div className="since-loading">
            Checking market activity...
          </div>
        </div>
      </section>
    );
  }

  if (error) {
    const checkpointMissing =
      error.message
        ?.toLowerCase()
        .includes("checkpoint not found");

    if (!checkpointMissing) {
      return (
        <section className="since-checked-section">
          <div className="since-checked-header">
            <div>
              <div className="section-kicker">
                SINCE YOU LAST CHECKED
              </div>

              <h2>What changed?</h2>

              <p>
                Meaningful movements detected
                in your watchlist.
              </p>
            </div>
          </div>

          <div className="since-error">
            <AlertCircle size={18} />

            <div>
              <strong>
                Unable to load market changes
              </strong>

              <span>
                {error.message}
              </span>
            </div>

            <button
              type="button"
              onClick={handleRefresh}
            >
              Retry
            </button>
          </div>
        </section>
      );
    }
  }

  /*
   * First visit / no checkpoint yet.
   */
  if (
    !summary ||
    summary.previousCheckpoint == null
  ) {
    return (
      <section className="since-checked-section">
        <div className="since-checked-header">
          <div>
            <div className="section-kicker">
              SINCE YOU LAST CHECKED
            </div>

            <h2>What changed?</h2>

            <p>
              StockStalk is ready to start
              tracking this watchlist.
            </p>
          </div>
        </div>

        <div className="first-check-state">
          <div className="first-check-icon">
            <CheckCircle2 size={25} />
          </div>

          <div>
            <strong>
              This is your first check
            </strong>

            <span>
              StockStalk will start tracking
              meaningful movements from this
              watchlist.
            </span>
          </div>
        </div>

        <div className="since-checked-footer">
          <div>
            <Clock3 size={13} />

            <span>
              Last checkpoint:
            </span>

            <strong>
              Not available — first check
            </strong>
          </div>

          <span>
            {summary?.totalStocks ?? 0} stocks tracked
          </span>
        </div>
      </section>
    );
  }

  return (
    <>
      <section className="since-checked-section">

        {/* =====================================================
            HEADER
        ===================================================== */}

        <div className="since-checked-header">
          <div>
            <div className="section-kicker">
              SINCE YOU LAST CHECKED
            </div>

            <h2>
              What changed?
            </h2>

            <p>
              StockStalk highlights the movements
              worth your attention.
            </p>
          </div>

          <div className="since-summary">

            <div className="since-summary-item">
              <strong>
                {summary.meaningfulChanges}
              </strong>

              <span>
                Meaningful
              </span>
            </div>

            <div className="since-summary-item">
              <strong>
                {summary.notableChanges}
              </strong>

              <span>
                Notable
              </span>
            </div>

            <div className="since-summary-item">
              <strong>
                {summary.unchanged}
              </strong>

              <span>
                Unchanged
              </span>
            </div>

          </div>
        </div>

        {/* =====================================================
            TOP THINGS YOU MISSED
        ===================================================== */}

        {topEvents.length > 0 && (
          <div className="top-missed-section">

            <div className="top-missed-header">
              <div className="top-missed-title">
                <div className="top-missed-icon">
                  <Flame size={17} />
                </div>

                <div>
                  <strong>
                    Focus first
                  </strong>

                  <span>
                    Top things you missed
                  </span>
                </div>
              </div>

              <span className="top-missed-count">
                {topEvents.length}{" "}
                {topEvents.length === 1
                  ? "priority"
                  : "priorities"}
              </span>
            </div>

            <div className="top-missed-list">
              {topEvents.map(
                (event, index) => {
                  const change =
                    event.changePercent ?? 0;

                  const score =
                    event.attentionScore ?? 0;

                  return (
                    <button
                      type="button"
                      key={
                        event.id ??
                        `${event.symbol}-${index}`
                      }
                      className="top-missed-item"
                      onClick={() =>
                        openEvent(event)
                      }
                    >
                      <div className="top-missed-rank">
                        {index === 0 ? (
                          <Flame size={16} />
                        ) : (
                          <Zap size={16} />
                        )}
                      </div>

                      <div className="top-missed-stock">
                        <strong>
                          {event.symbol}
                        </strong>

                        <span>
                          {change >= 0 ? "+" : ""}
                          {change.toFixed(2)}%
                          {event.volumeRatio != null
                            ? ` · ${event.volumeRatio.toFixed(
                                2
                              )}× volume`
                            : ""}
                        </span>
                      </div>

                      <AttentionBadge
                        score={score}
                      />
                    </button>
                  );
                }
              )}
            </div>
          </div>
        )}

        {/* =====================================================
            NO CHANGES
        ===================================================== */}

        {sortedChanges.length === 0 && (
          <div className="no-changes-card">

            <div className="no-changes-icon">
              <Activity size={22} />
            </div>

            <div>
              <strong>
                No meaningful changes
              </strong>

              <span>
                Your watchlist has not shown
                any significant movement since
                your last check.
              </span>
            </div>

          </div>
        )}

        {/* =====================================================
            FILTERS
        ===================================================== */}

        {sortedChanges.length > 0 && (
          <div className="event-filter-bar">

            <div className="event-filter-label">
              <Activity size={15} />
              <span>Filter events</span>
            </div>

            <div className="event-filter-buttons">

              <button
                type="button"
                className={
                  activeFilter === "ALL"
                    ? "event-filter-button active"
                    : "event-filter-button"
                }
                onClick={() =>
                  setActiveFilter("ALL")
                }
              >
                All
                <span>
                  {sortedChanges.length}
                </span>
              </button>

              <button
                type="button"
                className={
                  activeFilter === "MEANINGFUL"
                    ? "event-filter-button active"
                    : "event-filter-button"
                }
                onClick={() =>
                  setActiveFilter("MEANINGFUL")
                }
              >
                Meaningful
                <span>
                  {summary.meaningfulChanges}
                </span>
              </button>

              <button
                type="button"
                className={
                  activeFilter === "NOTABLE"
                    ? "event-filter-button active"
                    : "event-filter-button"
                }
                onClick={() =>
                  setActiveFilter("NOTABLE")
                }
              >
                Notable
                <span>
                  {summary.notableChanges}
                </span>
              </button>

              <button
                type="button"
                className={
                  activeFilter === "HIGH_ATTENTION"
                    ? "event-filter-button active"
                    : "event-filter-button"
                }
                onClick={() =>
                  setActiveFilter(
                    "HIGH_ATTENTION"
                  )
                }
              >
                High attention
                <span>
                  {
                    sortedChanges.filter(
                      (event) =>
                        (event.attentionScore ??
                          0) >= 75
                    ).length
                  }
                </span>
              </button>

            </div>
          </div>
        )}

        {/* =====================================================
            FILTERED EVENTS
        ===================================================== */}

        {sortedChanges.length > 0 && (
          <div className="change-event-list">

            {filteredChanges.length === 0 ? (
              <div className="filtered-empty-state">
                <Activity size={20} />

                <strong>
                  No events in this filter
                </strong>

                <span>
                  Try another attention level.
                </span>
              </div>
            ) : (
              filteredChanges.map(
                (event, index) => {

                  const change =
                    event.changePercent ?? 0;

                  const price =
                    event.currentValue;

                  const volumeRatio =
                    event.volumeRatio;

                  return (
                    <button
                      type="button"
                      key={
                        event.id ??
                        `${event.instrumentId}-${index}`
                      }
                      className="change-event-card"
                      onClick={() =>
                        openEvent(event)
                      }
                    >

                      {/* TOP */}

                      <div className="change-event-top">

                        <div className="change-event-stock">

                          <div className="change-event-icon">
                            <TrendingUp size={18} />
                          </div>

                          <div>
                            <strong>
                              {event.symbol}
                            </strong>

                            <span>
                              {event.type
                                ?.replaceAll(
                                  "_",
                                  " "
                                ) ??
                                "Market event"}
                            </span>
                          </div>

                        </div>

                        <AttentionBadge
                          score={
                            event.attentionScore
                          }
                        />

                      </div>

                      {/* MAIN METRICS */}

                      <div className="change-event-metrics">

                        <div>
                          <span>
                            Price movement
                          </span>

                          <strong
                            className={
                              change >= 0
                                ? "positive-value"
                                : "negative-value"
                            }
                          >
                            {change >= 0
                              ? "+"
                              : ""}
                            {change.toFixed(2)}%
                          </strong>
                        </div>

                        <div>
                          <span>
                            Current price
                          </span>

                          <strong>
                            {price != null
                              ? `₹${price.toLocaleString(
                                  "en-IN",
                                  {
                                    minimumFractionDigits: 2,
                                    maximumFractionDigits: 2,
                                  }
                                )}`
                              : "—"}
                          </strong>
                        </div>

                        <div>
                          <span>
                            Volume
                          </span>

                          <strong>
                            {volumeRatio != null
                              ? `${volumeRatio.toFixed(
                                  2
                                )}×`
                              : "—"}
                          </strong>

                          <small>
                            baseline
                          </small>
                        </div>

                      </div>

                      {/* SIGNALS */}

                      <div className="change-event-signals">

                        {event.priceAnomaly && (
                          <span className="event-signal">
                            <TrendingUp size={12} />
                            Price anomaly
                          </span>
                        )}

                        {event.volumeAnomaly && (
                          <span className="event-signal">
                            <Volume2 size={12} />
                            Unusual volume
                          </span>
                        )}

                        {event.near52WeekHigh && (
                          <span className="event-signal">
                            Near 52W high
                          </span>
                        )}

                        {event.near52WeekLow && (
                          <span className="event-signal">
                            Near 52W low
                          </span>
                        )}

                      </div>

                      {/* CONTEXT */}

                      {event.context && (
                        <div className="change-event-context">

                          <strong>
                            Why it matters
                          </strong>

                          <p>
                            {event.context}
                          </p>

                        </div>
                      )}

                      {/* FOOTER */}

                      <div className="change-event-footer">

                        <div className="event-confidence">

                          <span>
                            Confidence
                          </span>

                          <strong>
                            {event.confidence ??
                              "Unknown"}
                          </strong>

                        </div>

                        <DataStatus
                          status={
                            event.reliabilityState
                          }
                          timestamp={
                            event.detectedAt
                          }
                        />

                      </div>

                    </button>
                  );
                }
              )
            )}

          </div>
        )}

        {/* =====================================================
            CHECKPOINT
        ===================================================== */}

        <div className="since-checked-footer">

          <div>
            <Clock3 size={13} />

            <span>
              Last checkpoint:
            </span>

            <strong>
              {summary.previousCheckpoint
                ? new Date(
                    summary.previousCheckpoint
                  ).toLocaleString()
                : "Not available"}
            </strong>
          </div>

          <span>
            {summary.totalStocks} stocks tracked
          </span>

        </div>

      </section>

      {/* EVENT DETAIL */}

      {selectedEvent && (
        <EventDetail
          event={{
            symbol:
              selectedEvent.symbol,

            name:
              selectedEvent.symbol,

            priceChangePercent:
              selectedEvent.changePercent,

            currentPrice:
              selectedEvent.currentValue,

            volumeRatio:
              selectedEvent.volumeRatio,

            priceAnomaly:
              selectedEvent.priceAnomaly,

            volumeAnomaly:
              selectedEvent.volumeAnomaly,

            attentionScore:
              selectedEvent.attentionScore,

            context:
              selectedEvent.context,

            confidence:
              selectedEvent.confidence,

            detectedAt:
              selectedEvent.detectedAt,

            reliabilityState:
              selectedEvent.reliabilityState,
          }}
          onClose={() =>
            setSelectedEvent(null)
          }
        />
      )}
    </>
  );
}