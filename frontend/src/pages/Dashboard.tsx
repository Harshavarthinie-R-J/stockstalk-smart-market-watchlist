import { gql } from "@apollo/client";
import { useQuery } from "@apollo/client/react";
import { Activity, Layers, RefreshCw } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import Header from "../components/Header";
import WatchlistManager from "../components/WatchlistManager";
import WatchlistTable from "../components/WatchlistTable";
import StockSearch from "../components/StockSearch";
import SinceLastChecked from "../components/SinceLastChecked";
import CheckpointButton from "../components/CheckpointButton";
import MarketVisualization from "../components/MarketVisualization";
import StockIdeas from "../components/StockIdeas";

import type {
  Watchlist,
  WatchlistStock,
  Instrument,
  Quote,
} from "../types";

const WATCHLISTS = gql`
  query Watchlists {
    watchlists {
      id
      name
      createdAt
      updatedAt
      stocks {
        instrumentId
        position
        addedAt
      }
    }
  }
`;

const MARKET = gql`
  query Market {
    market {
      market
      lastUpdated

      instruments {
        id
        symbol
        name
        exchange
        segment
        isin
        currency
        sector
        industry
      }

      quotes {
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
  }
`;

interface WatchlistsQueryData {
  watchlists: Watchlist[];
}

interface MarketQueryData {
  market: {
    market: string;
    lastUpdated: string;
    instruments: Instrument[];
    quotes: Quote[];
  };
}

export default function Dashboard({
  onLogout,
}: {
  onLogout: () => void;
}) {
  const {
    data: watchlistData,
    loading: watchlistsLoading,
    error: watchlistsError,
    refetch: refetchWatchlists,
  } = useQuery<WatchlistsQueryData>(WATCHLISTS, {
    fetchPolicy: "network-only",
  });

  const {
    data: marketData,
    loading: marketLoading,
    error: marketError,
    refetch: refetchMarket,
  } = useQuery<MarketQueryData>(MARKET, {
    fetchPolicy: "network-only",
  });

  const watchlists = watchlistData?.watchlists ?? [];

  const [activeWatchlistId, setActiveWatchlistId] =
    useState<string>("");

  useEffect(() => {
    if (
      watchlists.length > 0 &&
      !watchlists.some(
        (watchlist) =>
          watchlist.id === activeWatchlistId
      )
    ) {
      setActiveWatchlistId(watchlists[0].id);
    }
  }, [watchlists, activeWatchlistId]);

  const activeWatchlist = useMemo(() => {
    return (
      watchlists.find(
        (watchlist) =>
          watchlist.id === activeWatchlistId
      ) ?? watchlists[0]
    );
  }, [watchlists, activeWatchlistId]);

  const instrumentMap = useMemo(() => {
    const map = new Map<string, Instrument>();

    for (const instrument of
      marketData?.market.instruments ?? []) {
      map.set(instrument.id, instrument);
    }

    return map;
  }, [marketData]);

  const quoteMap = useMemo(() => {
    const map = new Map<string, Quote>();

    for (const quote of
      marketData?.market.quotes ?? []) {
      map.set(quote.instrumentId, quote);
    }

    return map;
  }, [marketData]);

  /*
   * Stocks belonging to the active watchlist.
   */
  const watchlistStocks = useMemo<WatchlistStock[]>(() => {
    if (!activeWatchlist) {
      return [];
    }

    return [...activeWatchlist.stocks]
      .sort((a, b) => a.position - b.position)
      .map((item) => {
        const instrument =
          instrumentMap.get(item.instrumentId);

        const quote =
          quoteMap.get(item.instrumentId);

        return {
          instrumentId: item.instrumentId,
          symbol:
            instrument?.symbol ??
            quote?.symbol ??
            "UNKNOWN",
          name:
            instrument?.name ??
            "Unknown instrument",
          exchange:
            instrument?.exchange ??
            "NSE",
          sector:
            instrument?.sector ?? null,
          industry:
            instrument?.industry ?? null,
          position: item.position,
          addedAt: item.addedAt,

          price: quote?.price ?? null,
          previousClose:
            quote?.previousClose ?? null,
          change: quote?.change ?? null,
          changePercent:
            quote?.changePercent ?? null,
          volume: quote?.volume ?? null,
          reliabilityState:
            quote?.reliabilityState ?? null,
        };
      });
  }, [
    activeWatchlist,
    instrumentMap,
    quoteMap,
  ]);

  /*
   * All available market instruments.
   * This is used by StockIdeas to discover stocks
   * outside the current watchlist.
   */
  const allMarketStocks = useMemo<WatchlistStock[]>(() => {
    return (marketData?.market.instruments ?? [])
      .map((instrument) => {
        const quote =
          quoteMap.get(instrument.id);

        return {
          instrumentId: instrument.id,
          symbol: instrument.symbol,
          name: instrument.name,
          exchange: instrument.exchange,
          sector: instrument.sector ?? null,
          industry: instrument.industry ?? null,
          position: 0,
          addedAt: "",

          price: quote?.price ?? null,
          previousClose:
            quote?.previousClose ?? null,
          change: quote?.change ?? null,
          changePercent:
            quote?.changePercent ?? null,
          volume: quote?.volume ?? null,
          reliabilityState:
            quote?.reliabilityState ?? null,
        };
      })
      .filter(
        (stock) =>
          stock.changePercent != null
      );
  }, [marketData, quoteMap]);

  const existingStockIds = useMemo(() => {
    return new Set(
      watchlistStocks.map(
        (stock) => stock.instrumentId
      )
    );
  }, [watchlistStocks]);

  async function refreshDashboard() {
    await Promise.all([
      refetchWatchlists(),
      refetchMarket(),
    ]);
  }

  if (watchlistsLoading && !watchlistData) {
    return (
      <div className="dashboard-loading">
        <RefreshCw
          size={24}
          className="spin"
        />
        Loading StockStalk...
      </div>
    );
  }

  if (watchlistsError) {
    return (
      <div className="dashboard-error">
        <h2>Unable to load StockStalk</h2>

        <p>
          We couldn't retrieve your watchlists.
        </p>

        <button
          type="button"
          onClick={refreshDashboard}
        >
          Try again
        </button>
      </div>
    );
  }

  return (
    <div className="stockstalk-app">
      <Header onLogout={function (): void {
        throw new Error("Function not implemented.");
      } } />

      <div className="stockstalk-layout">

        {/* =================================================
            LEFT COLUMN
            ================================================= */}

        <aside className="stockstalk-sidebar">
          <WatchlistManager
            watchlists={watchlists}
            activeWatchlistId={
              activeWatchlist?.id ?? ""
            }
            onSelect={setActiveWatchlistId}
            onChanged={refreshDashboard}
          />
        </aside>


        {/* =================================================
            CENTER COLUMN
            ================================================= */}

        <main className="dashboard-main">

          {!activeWatchlist ? (
            <div className="empty-dashboard">
              <div className="empty-dashboard-icon">
                <Layers size={32} />
              </div>

              <h2>
                Create your first watchlist
              </h2>

              <p>
                Create a watchlist from the left
                sidebar to start tracking market
                activity.
              </p>
            </div>
          ) : (
            <>
              <SinceLastChecked
                watchlistId={activeWatchlist.id}
              />

              <section className="watchlist-summary">

                <div className="watchlist-summary-left">

                  <div className="watchlist-summary-icon">
                    <Activity size={18} />
                  </div>

                  <div>
                    <div className="summary-kicker">
                      CURRENT LIST
                    </div>

                    <h2>
                      {activeWatchlist.name}
                    </h2>
                  </div>

                </div>

                <div className="watchlist-summary-right">

                  <div className="watchlist-count">
                    <strong>
                      {watchlistStocks.length}
                    </strong>

                    <span>
                      {watchlistStocks.length === 1
                        ? "stock"
                        : "stocks"}
                    </span>
                  </div>

                  <CheckpointButton
                    watchlistId={
                      activeWatchlist.id
                    }
                    onCheckpointCreated={
                      refreshDashboard
                    }
                  />

                </div>

              </section>

              <StockSearch
                watchlistId={
                  activeWatchlist.id
                }
                existingStockIds={
                  Array.from(existingStockIds)
                }
                onStockAdded={
                  refreshDashboard
                }
              />

              <section className="watchlist-section">

                <div className="section-header">

                  <div>
                    <div className="section-kicker">
                      WATCHLIST
                    </div>

                    <h2>
                      Tracked Securities
                    </h2>

                    <p>
                      Current market state and
                      daily performance.
                    </p>
                  </div>

                  {marketLoading && (
                    <div className="section-loading">
                      <RefreshCw
                        size={14}
                        className="spin"
                      />
                      Updating
                    </div>
                  )}

                </div>

                {marketError ? (
                  <div className="inline-error">
                    Unable to load current market
                    quotes.
                  </div>
                ) : (
                  <WatchlistTable
                    watchlistId={
                      activeWatchlist.id
                    }
                    stocks={watchlistStocks}
                    onStockRemoved={
                      refreshDashboard
                    }
                  />
                )}

              </section>

              <MarketVisualization
                stocks={watchlistStocks}
              />

              <section className="market-snapshot">

                <div className="market-snapshot-header">

                  <div>
                    <div className="section-kicker">
                      MARKET SNAPSHOT
                    </div>

                    <h2>
                      Major Watchlist Movers
                    </h2>

                    <p>
                      Largest current percentage
                      swings in this list.
                    </p>
                  </div>

                  <div className="market-status">
                    <span className="market-status-dot" />

                    <span>
                      {marketData?.market?.market ??
                        "NSE / BSE"}
                    </span>
                  </div>

                </div>

                <div className="snapshot-grid">

                  {watchlistStocks
                    .filter(
                      (stock) =>
                        stock.changePercent != null
                    )
                    .sort(
                      (a, b) =>
                        Math.abs(
                          b.changePercent ?? 0
                        ) -
                        Math.abs(
                          a.changePercent ?? 0
                        )
                    )
                    .slice(0, 3)
                    .map((stock) => {

                      const change =
                        stock.changePercent ?? 0;

                      return (
                        <div
                          className="market-snapshot-card"
                          key={
                            stock.instrumentId
                          }
                        >
                          <span>
                            {stock.symbol}
                          </span>

                          <strong>
                            {stock.price != null
                              ? `₹${stock.price.toLocaleString(
                                  "en-IN",
                                  {
                                    minimumFractionDigits: 2,
                                    maximumFractionDigits: 2,
                                  }
                                )}`
                              : "—"}
                          </strong>

                          <div
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
                          </div>
                        </div>
                      );
                    })}

                </div>

              </section>
            </>
          )}

          <footer className="dashboard-footer">
            <span>
              StockStalk
            </span>

            <span>
              Market information only ·
              Not investment advice
            </span>
          </footer>

        </main>


        {/* =================================================
            RIGHT COLUMN
            ================================================= */}

        {activeWatchlist && (
          <StockIdeas
            watchlistId={
              activeWatchlist.id
            }
            stocks={watchlistStocks}
            allStocks={allMarketStocks}
            onStockAdded={
              refreshDashboard
            }
          />
        )}

      </div>
    </div>
  );
}