import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";

import { BarChart3, LineChart as LineIcon, X } from "lucide-react";
import { useState } from "react";

interface StockData {
  symbol: string;
  price?: number | null;
  previousClose?: number | null;
  changePercent?: number | null;
  volume?: number | null;
}

interface MarketVisualizationProps {
  stocks: StockData[];
}

type ViewType = "price" | "market";

export default function MarketVisualization({
  stocks,
}: MarketVisualizationProps) {
  const [open, setOpen] = useState(false);

  const [view, setView] =
    useState<ViewType>("market");

  const movementData = stocks
    .map((stock) => ({
      symbol: stock.symbol,
      change: Number(
        stock.changePercent ?? 0
      ),
    }))
    .filter((stock) =>
      Number.isFinite(stock.change)
    );

  const trendStocks = stocks.filter(
    (stock) =>
      stock.previousClose != null &&
      stock.price != null
  );

  /*
   * Price comparison:
   * Previous close -> Current price
   */
  const priceTrendData = [
    {
      point: "Previous Close",
      ...Object.fromEntries(
        trendStocks.map((stock) => [
          stock.symbol,
          stock.previousClose,
        ])
      ),
    },
    {
      point: "Current Price",
      ...Object.fromEntries(
        trendStocks.map((stock) => [
          stock.symbol,
          stock.price,
        ])
      ),
    },
  ];

  if (movementData.length === 0) {
    return (
      <section className="market-visualization-card">
        <div className="market-visualization-card-icon">
          <BarChart3 size={22} />
        </div>

        <div className="market-visualization-card-content">
          <div className="visualization-kicker">
            MARKET VISUALIZATION
          </div>

          <h2>
            Explore your market movement
          </h2>

          <p>
            Add stocks to your watchlist to
            visualize price and market trends.
          </p>
        </div>
      </section>
    );
  }

  return (
    <>
      {/* =================================================
          COMPACT CARD
      ================================================= */}

      <button
        type="button"
        className="market-visualization-card"
        onClick={() => setOpen(true)}
      >
        <div className="market-visualization-card-icon">
          <BarChart3 size={22} />
        </div>

        <div className="market-visualization-card-content">
          <div className="visualization-kicker">
            MARKET VISUALIZATION
          </div>

          <h2>
            Explore your market movement
          </h2>

          <p>
            View price trends and market
            movement across your watchlist.
          </p>
        </div>

        <div className="market-visualization-card-action">
          <span>
            {movementData.length}{" "}
            {movementData.length === 1
              ? "stock"
              : "stocks"}
          </span>

          <strong>
            View graphs →
          </strong>
        </div>
      </button>

      {/* =================================================
          MODAL
      ================================================= */}

      {open && (
        <div
          className="visualization-modal-overlay"
          onClick={() => setOpen(false)}
        >
          <div
            className="visualization-modal"
            onClick={(event) =>
              event.stopPropagation()
            }
          >
            {/* HEADER */}

            <div className="visualization-modal-header">
              <div>
                <div className="visualization-kicker">
                  MARKET VISUALIZATION
                </div>

                <h2>
                  Understand your watchlist
                </h2>

                <p>
                  Compare price and market movement
                  using the latest available data.
                </p>
              </div>

              <button
                type="button"
                className="visualization-close"
                onClick={() => setOpen(false)}
                aria-label="Close visualization"
              >
                <X size={20} />
              </button>
            </div>

            {/* TABS */}

            <div className="visualization-tabs">

              <button
                type="button"
                className={
                  view === "market"
                    ? "visualization-tab active"
                    : "visualization-tab"
                }
                onClick={() =>
                  setView("market")
                }
              >
                <BarChart3 size={17} />
                Market Trend
              </button>

              <button
                type="button"
                className={
                  view === "price"
                    ? "visualization-tab active"
                    : "visualization-tab"
                }
                onClick={() =>
                  setView("price")
                }
              >
                <LineIcon size={17} />
                Price Trend
              </button>

            </div>

            {/* =================================================
                MARKET TREND
            ================================================= */}

            {view === "market" && (
              <div className="visualization-panel">

                <div className="visualization-panel-heading">
                  <div>
                    <h3>
                      Market Trend
                    </h3>

                    <p>
                      Current percentage movement
                      across your watchlist.
                    </p>
                  </div>
                </div>

                <div className="modal-chart">
                  <ResponsiveContainer
                    width="100%"
                    height={400}
                  >
                    <BarChart
                      data={movementData}
                      margin={{
                        top: 25,
                        right: 25,
                        left: 40,
                        bottom: 25,
                      }}
                      barCategoryGap="25%"
                    >
                      <CartesianGrid
                        strokeDasharray="4 4"
                        vertical={false}
                      />

                      <XAxis
                        dataKey="symbol"
                        tick={{
                          fontSize: 12,
                          fontWeight: 700,
                        }}
                        axisLine={false}
                        tickLine={false}
                        dy={10}
                      />

                      <YAxis
                        width={50}
                        tick={{
                          fontSize: 11,
                        }}
                        axisLine={false}
                        tickLine={false}
                        tickFormatter={(value) =>
                          `${Number(value).toFixed(1)}%`
                        }
                        domain={["auto", "auto"]}
                      />

                      <Tooltip
                        cursor={{
                          fill:
                            "rgba(120, 140, 170, 0.08)",
                        }}
                        contentStyle={{
                          background: "#111824",
                          border:
                            "1px solid #2b3a50",
                          borderRadius: "12px",
                          color: "#ffffff",
                        }}
                        formatter={(value) => [
                          `${Number(
                            value
                          ).toFixed(2)}%`,
                          "24h Move",
                        ]}
                      />

                      <Bar
                        dataKey="change"
                        radius={[
                          7,
                          7,
                          0,
                          0,
                        ]}
                        maxBarSize={70}
                      >
                        {movementData.map(
                          (entry) => (
                            <Cell
                              key={
                                entry.symbol
                              }
                              fill={
                                entry.change >=
                                0
                                  ? "#00c896"
                                  : "#ef5b61"
                              }
                            />
                          )
                        )}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                </div>

                <div className="visualization-explanation">
                  <div>
                    <span className="explanation-dot positive" />
                    Positive movement
                  </div>

                  <div>
                    <span className="explanation-dot negative" />
                    Negative movement
                  </div>
                </div>

              </div>
            )}

            {/* =================================================
                PRICE TREND
            ================================================= */}

            {view === "price" && (
              <div className="visualization-panel">

                <div className="visualization-panel-heading">
                  <div>
                    <h3>
                      Price Trend
                    </h3>

                    <p>
                      Previous closing price compared
                      with the current price.
                    </p>
                  </div>
                </div>

                {trendStocks.length === 0 ? (
                  <div className="visualization-no-data">
                    <LineIcon size={30} />

                    <strong>
                      Price trend unavailable
                    </strong>

                    <span>
                      The backend has not provided
                      previous-close data for these
                      stocks yet.
                    </span>
                  </div>
                ) : (
                  <>
                    <div className="modal-chart">
                      <ResponsiveContainer
                        width="100%"
                        height={400}
                      >
                        <LineChart
                          data={priceTrendData}
                          margin={{
                            top: 25,
                            right: 30,
                            left: 45,
                            bottom: 25,
                          }}
                        >
                          <CartesianGrid
                            strokeDasharray="4 4"
                            vertical={false}
                          />

                          <XAxis
                            dataKey="point"
                            tick={{
                              fontSize: 12,
                              fontWeight: 700,
                            }}
                            axisLine={false}
                            tickLine={false}
                            dy={10}
                          />

                          <YAxis
                            width={70}
                            tick={{
                              fontSize: 11,
                            }}
                            axisLine={false}
                            tickLine={false}
                            tickFormatter={(value) =>
                              `₹${Number(
                                value
                              ).toLocaleString(
                                "en-IN"
                              )}`
                            }
                            domain={["auto", "auto"]}
                          />

                          <Tooltip
                            contentStyle={{
                              background:
                                "#111824",
                              border:
                                "1px solid #2b3a50",
                              borderRadius:
                                "12px",
                              color: "#ffffff",
                            }}
                            formatter={(
                              value,
                              name
                            ) => [
                              value != null
                                ? `₹${Number(
                                    value
                                  ).toLocaleString(
                                    "en-IN",
                                    {
                                      minimumFractionDigits: 2,
                                      maximumFractionDigits: 2,
                                    }
                                  )}`
                                : "—",
                              String(name),
                            ]}
                          />

                          {trendStocks.map(
                            (stock, index) => {
                              const colors = [
                                "#5b6fff",
                                "#00c896",
                                "#ef5b61",
                                "#f5a623",
                                "#a56eff",
                                "#28a7e8",
                              ];

                              return (
                                <Line
                                  key={
                                    stock.symbol
                                  }
                                  type="monotone"
                                  dataKey={
                                    stock.symbol
                                  }
                                  name={
                                    stock.symbol
                                  }
                                  stroke={
                                    colors[
                                      index %
                                        colors.length
                                    ]
                                  }
                                  strokeWidth={3}
                                  dot={{
                                    r: 5,
                                  }}
                                  activeDot={{
                                    r: 7,
                                  }}
                                />
                              );
                            }
                          )}

                        </LineChart>
                      </ResponsiveContainer>
                    </div>

                    <div className="visualization-explanation">
                      <span>
                        Previous Close
                      </span>

                      <span>
                        →
                      </span>

                      <span>
                        Current Price
                      </span>
                    </div>
                  </>
                )}

              </div>
            )}

            {/* FOOTER */}

            <div className="visualization-modal-footer">
              <span>
                Based on market data provided by
                StockStalk's backend.
              </span>

              <span>
                This visualization describes
                movement and does not provide
                trading recommendations.
              </span>
            </div>

          </div>
        </div>
      )}
    </>
  );
}