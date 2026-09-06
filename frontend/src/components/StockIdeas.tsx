import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { Plus, TrendingUp, Activity, Check, Loader2 } from "lucide-react";
import type { WatchlistStock } from "../types";

const ADD_STOCK = gql`
  mutation AddStock($watchlistId: ID!, $instrumentId: ID!) {
    addStock(
      watchlistId: $watchlistId
      instrumentId: $instrumentId
    ) {
      instrumentId
      position
      addedAt
    }
  }
`;

interface StockIdeasProps {
  watchlistId: string;
  stocks: WatchlistStock[];
  allStocks: WatchlistStock[];
  onStockAdded: () => void;
}

export default function StockIdeas({
  watchlistId,
  stocks,
  allStocks,
  onStockAdded,
}: StockIdeasProps) {
  const [addStock, { loading }] = useMutation(ADD_STOCK);

  const existingIds = new Set(
    stocks.map((stock) => stock.instrumentId)
  );

  /*
   * We use actual market movement from the backend.
   * No prediction or fake recommendation is generated.
   *
   * Stocks with the strongest absolute movement are
   * surfaced as "Stocks to Consider".
   */
  const ideas = allStocks
    .filter((stock) => !existingIds.has(stock.instrumentId))
    .filter((stock) => stock.changePercent != null)
    .sort(
      (a, b) =>
        Math.abs(b.changePercent ?? 0) -
        Math.abs(a.changePercent ?? 0)
    )
    .slice(0, 5);

  async function handleAdd(stock: WatchlistStock) {
    try {
      await addStock({
        variables: {
          watchlistId,
          instrumentId: stock.instrumentId,
        },
      });

      onStockAdded();
    } catch (error) {
      console.error("Unable to add stock:", error);
      alert(`Unable to add ${stock.symbol}.`);
    }
  }

  return (
    <aside className="stock-ideas-panel">
      <div className="stock-ideas-header">
        <div className="stock-ideas-heading">
          <div className="stock-ideas-icon">
            <TrendingUp size={18} />
          </div>

          <div>
            <div className="section-kicker">
              DISCOVERY
            </div>

            <h2>Stocks to Consider</h2>

            <p>
              Unusual market activity worth exploring.
            </p>
          </div>
        </div>
      </div>

      <div className="stock-ideas-disclaimer">
        <Activity size={13} />

        <span>
          Ranked by observed market movement, not a
          buy recommendation.
        </span>
      </div>

      {ideas.length === 0 ? (
        <div className="stock-ideas-empty">
          <div className="stock-ideas-empty-icon">
            <Activity size={20} />
          </div>

          <strong>No new ideas right now</strong>

          <span>
            Stocks outside this watchlist will appear
            here when meaningful movement is observed.
          </span>
        </div>
      ) : (
        <div className="stock-ideas-list">
          {ideas.map((stock) => {
            const change = stock.changePercent ?? 0;
            const positive = change >= 0;

            return (
              <div
                className="stock-idea-card"
                key={stock.instrumentId}
              >
                <div className="stock-idea-top">
                  <div>
                    <strong>{stock.symbol}</strong>

                    <span className="stock-idea-name">
                      {stock.name}
                    </span>
                  </div>

                  <div
                    className={
                      positive
                        ? "stock-idea-change positive-value"
                        : "stock-idea-change negative-value"
                    }
                  >
                    {positive ? "+" : ""}
                    {change.toFixed(2)}%
                  </div>
                </div>

                <div className="stock-idea-meta">
                  <span>
                    {stock.exchange}
                  </span>

                  {stock.sector && (
                    <span>
                      {stock.sector}
                    </span>
                  )}
                </div>

                <div className="stock-idea-bottom">
                  <div className="stock-idea-price">
                    {stock.price != null
                      ? `₹${stock.price.toLocaleString(
                          "en-IN",
                          {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                          }
                        )}`
                      : "Price unavailable"}
                  </div>

                  <button
                    type="button"
                    className="stock-idea-add"
                    disabled={loading}
                    onClick={() => handleAdd(stock)}
                  >
                    {loading ? (
                      <Loader2
                        size={14}
                        className="spin"
                      />
                    ) : (
                      <Plus size={14} />
                    )}

                    Add
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <div className="stock-ideas-footer">
        <span>
          Suggestions use current observed market data.
        </span>
      </div>
    </aside>
  );
}