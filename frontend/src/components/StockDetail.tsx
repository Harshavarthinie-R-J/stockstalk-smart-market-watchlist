import { gql } from "@apollo/client";
import { useQuery } from "@apollo/client/react";
import PersonalBaseline from "./PersonalBaseline";
import EventHistory from "./EventHistory";
import {
  Activity,
  BarChart3,
  Clock3,
  TrendingDown,
  TrendingUp,
  X,
} from "lucide-react";

import AttentionBadge from "./AttentionBadge";
import DataStatus from "./DataStatus";

const GET_QUOTE = gql`
  query GetQuote($symbol: String!) {
    quote(symbol: $symbol) {
      instrumentId
      symbol
      price
      previousClose
      open
      high
      low
      volume
      change
      changePercent
      week52High
      week52Low
      marketStatus
      marketTimestamp
      receivedTimestamp
      source
      reliabilityState
    }
  }
`;

interface StockDetailProps {
  symbol?: string | null;
  name?: string | null;
  sector?: string | null;
  onClose: () => void;
}

interface Quote {
  instrumentId: string;
  symbol: string;
  price?: number | null;
  previousClose?: number | null;
  open?: number | null;
  high?: number | null;
  low?: number | null;
  volume?: number | null;
  change?: number | null;
  changePercent?: number | null;
  week52High?: number | null;
  week52Low?: number | null;
  marketStatus?: string | null;
  marketTimestamp?: string | null;
  receivedTimestamp?: string | null;
  source?: string | null;
  reliabilityState?: string | null;
}

interface QuoteData {
  quote: Quote | null;
}

export default function StockDetail({
  symbol,
  name,
  sector,
  onClose,
}: StockDetailProps) {
  /*
   * Make sure we always send a real string
   * to GraphQL.
   */
  const normalizedSymbol =
    typeof symbol === "string"
      ? symbol.trim()
      : "";

  const { data, loading, error } =
    useQuery<QuoteData>(GET_QUOTE, {
      variables: {
        symbol: normalizedSymbol,
      },
      skip: normalizedSymbol.length === 0,
      fetchPolicy: "network-only",
    });

  /*
   * If no symbol was supplied, show a useful error
   * instead of sending an invalid GraphQL request.
   */
  if (!normalizedSymbol) {
    return (
      <div
        className="stock-detail-overlay"
        onClick={onClose}
      >
        <div
          className="stock-detail-modal"
          onClick={(event) =>
            event.stopPropagation()
          }
        >
          <div className="stock-detail-header">
            <div className="stock-detail-kicker">
              STOCK DETAIL
            </div>

            <div className="stock-detail-title-row">
              <div>
                <h2>Stock unavailable</h2>

                <p>
                  No stock symbol was supplied.
                </p>
              </div>

              <button
                type="button"
                className="stock-detail-close"
                onClick={onClose}
                aria-label="Close"
              >
                <X size={20} />
              </button>
            </div>
          </div>

          <div className="stock-detail-error">
            <strong>
              Unable to open stock details
            </strong>

            <span>
              Please select a valid stock from
              your watchlist.
            </span>
          </div>
        </div>
      </div>
    );
  }

  const quote = data?.quote;

  const changePercent =
    quote?.changePercent ?? 0;

  const isPositive =
    changePercent >= 0;

  const price =
    quote?.price ?? null;

  const previousClose =
    quote?.previousClose ?? null;

  const volume =
    quote?.volume ?? null;

  /*
   * Temporary baseline.
   *
   * We will replace this with the actual
   * historical baseline in the Personal
   * Baseline feature.
   */
  const averageVolume = 7425000;

  const volumeRatio =
    volume != null && averageVolume > 0
      ? volume / averageVolume
      : null;

  const volumeAnomaly =
    volumeRatio != null &&
    volumeRatio >= 1.5;

  const attentionScore =
    volumeAnomaly ||
    Math.abs(changePercent) >= 3
      ? 75
      : Math.abs(changePercent) >= 1.5
        ? 50
        : 25;

  const pricePosition =
    price != null &&
    quote?.week52Low != null &&
    quote?.week52High != null &&
    quote.week52High !== quote.week52Low
      ? ((price - quote.week52Low) /
          (quote.week52High -
            quote.week52Low)) *
        100
      : null;

  return (
    <div
      className="stock-detail-overlay"
      onClick={onClose}
    >
      <div
        className="stock-detail-modal"
        onClick={(event) =>
          event.stopPropagation()
        }
      >
        {/* HEADER */}

        <div className="stock-detail-header">
          <div className="stock-detail-kicker">
            STOCK DETAIL
          </div>

          <div className="stock-detail-title-row">
            <div>
              <h2>{normalizedSymbol}</h2>

              <p>
                {name || normalizedSymbol}

                {sector
                  ? ` · ${sector}`
                  : ""}
              </p>
            </div>

            <button
              type="button"
              className="stock-detail-close"
              onClick={onClose}
              aria-label="Close"
            >
              <X size={20} />
            </button>
          </div>
        </div>

        {/* LOADING */}

        {loading && (
          <div className="stock-detail-loading">
            <Activity size={22} />

            <span>
              Loading latest market data...
            </span>
          </div>
        )}

        {/* ERROR */}

        {error && !loading && (
          <div className="stock-detail-error">
            <strong>
              Unable to load stock data
            </strong>

            <span>
              {error.message}
            </span>
          </div>
        )}

        {/* CONTENT */}

        {quote &&
          !loading &&
          !error && (
            <>
              {/* PRICE */}

              <div className="stock-detail-price-section">
                <div>
                  <span className="stock-detail-label">
                    Current price
                  </span>

                  <strong className="stock-detail-price">
                    ₹
                    {price != null
                      ? price.toLocaleString(
                          "en-IN",
                          {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                          }
                        )
                      : "—"}
                  </strong>
                </div>

                <div
                  className={
                    isPositive
                      ? "stock-detail-change positive-value"
                      : "stock-detail-change negative-value"
                  }
                >
                  {isPositive ? (
                    <TrendingUp size={18} />
                  ) : (
                    <TrendingDown size={18} />
                  )}

                  <strong>
                    {isPositive
                      ? "+"
                      : ""}
                    {changePercent.toFixed(
                      2
                    )}
                    %
                  </strong>
                </div>
              </div>

              {/* METRICS */}

              <div className="stock-detail-metrics">
                <Metric
                  label="Previous close"
                  value={formatPrice(
                    previousClose
                  )}
                />

                <Metric
                  label="Open"
                  value={formatPrice(
                    quote.open
                  )}
                />

                <Metric
                  label="Day high"
                  value={formatPrice(
                    quote.high
                  )}
                />

                <Metric
                  label="Day low"
                  value={formatPrice(
                    quote.low
                  )}
                />

                <Metric
                  label="Volume"
                  value={formatNumber(
                    volume
                  )}
                />

                <Metric
                  label="52W high"
                  value={formatPrice(
                    quote.week52High
                  )}
                />

                <Metric
                  label="52W low"
                  value={formatPrice(
                    quote.week52Low
                  )}
                />

                <Metric
                  label="Market"
                  value={
                    quote.marketStatus ||
                    "Unknown"
                  }
                />
              </div>

              {/* STOCKSTALK INTELLIGENCE */}
              <PersonalBaseline
                symbol={normalizedSymbol}
                changePercent={quote.changePercent}
                currentVolume={quote.volume}
                averageVolume={undefined}
               />

              <div className="stock-detail-intelligence">
                <div className="stock-detail-section-title">
                  <BarChart3 size={17} />

                  <span>
                    StockStalk intelligence
                  </span>
                </div>

                <div className="stock-detail-attention">
                  <div>
                    <span>
                      Attention
                    </span>

                    <AttentionBadge
                      score={
                        attentionScore
                      }
                    />
                  </div>

                  <div>
                    <span>
                      Volume
                    </span>

                    <strong>
                      {volumeRatio !=
                      null
                        ? `${volumeRatio.toFixed(
                            2
                          )}× baseline`
                        : "Unavailable"}
                    </strong>
                  </div>

                  <div>
                    <span>
                      Movement
                    </span>

                    <strong>
                      {Math.abs(
                        changePercent
                      ) >= 1.5
                        ? "Unusual"
                        : "Within normal range"}
                    </strong>
                  </div>
                </div>

                <div className="stock-detail-why">
                  <strong>
                    Why StockStalk noticed it
                  </strong>

                  <p>
                    {volumeAnomaly &&
                    Math.abs(
                      changePercent
                    ) >= 1.5
                      ? "Price movement occurred together with unusually high trading volume."
                      : volumeAnomaly
                        ? "Trading volume is unusually high compared with the available baseline."
                        : Math.abs(
                              changePercent
                            ) >= 1.5
                          ? "The price movement is larger than the normal movement threshold."
                          : "No major unusual movement is currently detected."}
                  </p>
                </div>
              </div>

              {/* 52 WEEK POSITION */}

              {pricePosition !=
                null && (
                <div className="stock-detail-range">
                  <div className="stock-detail-range-header">
                    <span>
                      52-week position
                    </span>

                    <strong>
                      {pricePosition.toFixed(
                        0
                      )}
                      %
                    </strong>
                  </div>

                  <div className="stock-detail-range-bar">
                    <div
                      className="stock-detail-range-fill"
                      style={{
                        width: `${Math.min(
                          100,
                          Math.max(
                            0,
                            pricePosition
                          )
                        )}%`,
                      }}
                    />
                  </div>

                  <div className="stock-detail-range-labels">
                    <span>
                      ₹
                      {quote.week52Low?.toLocaleString(
                        "en-IN"
                      ) ?? "—"}
                    </span>

                    <span>
                      ₹
                      {quote.week52High?.toLocaleString(
                        "en-IN"
                      ) ?? "—"}
                    </span>
                  </div>
                </div>
              )}

              {/* STATUS */}

              <div className="stock-detail-status">
                <div>
                  <Clock3 size={14} />

                  <span>
                    Data status
                  </span>
                </div>

                <DataStatus
                  status={
                    quote.reliabilityState
                  }
                  timestamp={
                    quote.receivedTimestamp ??
                    quote.marketTimestamp
                  }
                />
                <EventHistory
                    instrumentId={quote.instrumentId}
                    symbol={normalizedSymbol}
                />
              </div>


              {/* FOOTER */}

              <div className="stock-detail-footer">
                <span>
                  Source:{" "}
                  {quote.source ||
                    "Market data provider"}
                </span>

                <span>
                  Market information only —
                  not an investment
                  recommendation.
                </span>
              </div>
            </>
          )}
      </div>
    </div>
  );
}

function Metric({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="stock-detail-metric">
      <span>{label}</span>

      <strong>{value}</strong>
    </div>
  );
}

function formatPrice(
  value?: number | null
) {
  if (value == null) {
    return "—";
  }

  return `₹${value.toLocaleString(
    "en-IN",
    {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }
  )}`;
}

function formatNumber(
  value?: number | null
) {
  if (value == null) {
    return "—";
  }

  return value.toLocaleString("en-IN");
}