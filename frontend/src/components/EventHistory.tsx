import { gql } from "@apollo/client";
import { useQuery } from "@apollo/client/react";
import {
  Clock3,
  History,
  TrendingDown,
  TrendingUp,
} from "lucide-react";

import AttentionBadge from "./AttentionBadge";
import DataStatus from "./DataStatus";

const EVENT_HISTORY = gql`
  query EventHistory(
    $instrumentId: ID!
    $limit: Int
  ) {
    eventHistory(
      instrumentId: $instrumentId
      limit: $limit
    ) {
      id
      instrumentId
      symbol
      type
      severity
      previousValue
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
    }
  }
`;

interface HistoryEvent {
  id: string;
  instrumentId: string;
  symbol: string;
  type?: string | null;
  severity?: string | null;
  previousValue?: number | null;
  currentValue?: number | null;
  changePercent?: number | null;
  currentVolume?: number | null;
  averageVolume?: number | null;
  volumeRatio?: number | null;
  volumeAnomaly?: boolean | null;
  priceAnomaly?: boolean | null;
  signals?: string[] | null;
  context?: string | null;
  confidence?: string | null;
  attentionScore?: number | null;
  detectedAt?: string | null;
}

interface EventHistoryProps {
  instrumentId: string;
  symbol: string;
}

export default function EventHistory({
  instrumentId,
  symbol,
}: EventHistoryProps) {
  const { data, loading, error } =
    useQuery<{
      eventHistory: HistoryEvent[];
    }>(EVENT_HISTORY, {
      variables: {
        instrumentId,
        limit: 20,
      },
      fetchPolicy: "network-only",
    });

  const events =
    data?.eventHistory ?? [];

  return (
    <div className="event-history-card">

      <div className="event-history-header">
        <div className="event-history-title">
          <div className="event-history-icon">
            <History size={18} />
          </div>

          <div>
            <strong>
              Event history
            </strong>

            <span>
              Previous StockStalk observations
              for {symbol}
            </span>
          </div>
        </div>

        <span className="event-history-count">
          {events.length} events
        </span>
      </div>

      {loading && (
        <div className="event-history-state">
          Loading event history...
        </div>
      )}

      {error && (
        <div className="event-history-state error">
          Unable to load event history.
        </div>
      )}

      {!loading &&
        !error &&
        events.length === 0 && (
          <div className="event-history-state">
            <History size={20} />

            <span>
              No previous events have been
              recorded for this stock.
            </span>
          </div>
        )}

      {!loading &&
        !error &&
        events.length > 0 && (
          <div className="event-history-list">
            {events.map((event) => {
              const change =
                event.changePercent ?? 0;

              return (
                <div
                  key={event.id}
                  className="event-history-item"
                >
                  <div className="history-event-icon">
                    {change >= 0 ? (
                      <TrendingUp size={15} />
                    ) : (
                      <TrendingDown size={15} />
                    )}
                  </div>

                  <div className="history-event-main">
                    <div className="history-event-top">
                      <strong>
                        {event.symbol}
                      </strong>

                      <AttentionBadge
                        score={
                          event.attentionScore
                        }
                      />
                    </div>

                    <div className="history-event-data">
                      <span
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
                      </span>

                      {event.volumeRatio !=
                        null && (
                        <span>
                          {event.volumeRatio.toFixed(
                            2
                          )}
                          × volume
                        </span>
                      )}

                      <span>
                        {event.severity ??
                          "Event"}
                      </span>
                    </div>

                    {event.context && (
                      <p>
                        {event.context}
                      </p>
                    )}

                    <div className="history-event-footer">
                      <span>
                        <Clock3 size={12} />

                        {event.detectedAt
                          ? new Date(
                              event.detectedAt
                            ).toLocaleString()
                          : "Unknown time"}
                      </span>

                      <DataStatus
                        status="STALE"
                        timestamp={
                          event.detectedAt
                        }
                      />
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
    </div>
  );
}