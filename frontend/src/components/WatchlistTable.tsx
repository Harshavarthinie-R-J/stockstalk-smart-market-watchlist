import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import {
  Loader2,
  Trash2,
} from "lucide-react";
import { useState } from "react";

import DataStatus from "./DataStatus";
import StockDetail from "./StockDetail";

import type { WatchlistStock } from "../types";

const REMOVE_STOCK = gql`
  mutation RemoveStock(
    $watchlistId: ID!
    $instrumentId: ID!
  ) {
    removeStock(
      watchlistId: $watchlistId
      instrumentId: $instrumentId
    )
  }
`;

interface WatchlistTableProps {
  watchlistId: string;
  stocks: WatchlistStock[];
  onStockRemoved: () => void;
}

export default function WatchlistTable({
  watchlistId,
  stocks,
  onStockRemoved,
}: WatchlistTableProps) {
  const [
    removeStock,
    { loading: removing },
  ] = useMutation(REMOVE_STOCK);

  /*
   * The complete selected stock object is stored here.
   * This guarantees that symbol/name/sector are available
   * when StockDetail opens.
   */
  const [selectedStock, setSelectedStock] =
    useState<WatchlistStock | null>(null);

  async function handleRemove(
    instrumentId: string,
    symbol: string
  ) {
    const confirmed = window.confirm(
      `Remove ${symbol} from this watchlist?`
    );

    if (!confirmed) {
      return;
    }

    try {
      await removeStock({
        variables: {
          watchlistId,
          instrumentId,
        },
      });

      /*
       * If the removed stock was currently open,
       * close the detail modal.
       */
      setSelectedStock(null);

      onStockRemoved();
    } catch (error) {
      console.error(
        "Unable to remove stock:",
        error
      );

      alert(
        "Unable to remove this stock."
      );
    }
  }

  if (stocks.length === 0) {
    return (
      <div className="empty-state">
        <strong>
          No stocks in this watchlist
        </strong>

        <span>
          Search above to add your first stock.
        </span>
      </div>
    );
  }

  return (
    <>
      <div className="watchlist-table-wrapper">
        <table className="watchlist-table">
          <thead>
            <tr>
              <th>Stock</th>
              <th>Price</th>
              <th>Change</th>
              <th>Volume</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>

          <tbody>
            {stocks.map((stock) => {
              const change =
                stock.changePercent ?? 0;

              return (
                <tr
                  key={stock.instrumentId}
                  className="watchlist-stock-row"
                  onClick={() => {
                    /*
                     * IMPORTANT:
                     * Store the complete stock object.
                     */
                    setSelectedStock(stock);
                  }}
                >
                  {/* STOCK */}

                  <td>
                    <div className="table-stock">
                      <strong>
                        {stock.symbol}
                      </strong>

                      <span>
                        {stock.name}
                      </span>

                      <small>
                        {stock.exchange}

                        {stock.sector
                          ? ` · ${stock.sector}`
                          : ""}
                      </small>
                    </div>
                  </td>

                  {/* PRICE */}

                  <td>
                    {stock.price != null
                      ? `₹${stock.price.toLocaleString(
                          "en-IN",
                          {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                          }
                        )}`
                      : "—"}
                  </td>

                  {/* CHANGE */}

                  <td>
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
                  </td>

                  {/* VOLUME */}

                  <td>
                    {stock.volume != null
                      ? stock.volume.toLocaleString(
                          "en-IN"
                        )
                      : "—"}
                  </td>

                  {/* STATUS */}

                  <td>
                    <DataStatus
                      status={
                        stock.reliabilityState
                      }
                    />
                  </td>

                  {/* REMOVE */}

                  <td>
                    <button
                      type="button"
                      className="table-remove-button"
                      title={`Remove ${stock.symbol}`}
                      disabled={removing}
                      onClick={(event) => {
                        /*
                         * Prevent the row click from
                         * opening StockDetail.
                         */
                        event.stopPropagation();

                        handleRemove(
                          stock.instrumentId,
                          stock.symbol
                        );
                      }}
                    >
                      {removing ? (
                        <Loader2
                          size={16}
                          className="spin"
                        />
                      ) : (
                        <Trash2 size={16} />
                      )}
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* =====================================================
          STOCK DETAIL MODAL
          ===================================================== */}

      {selectedStock && (
        <StockDetail
          /*
           * THIS IS THE IMPORTANT PART
           *
           * The symbol now comes directly from
           * the actual WatchlistStock object.
           */
          symbol={selectedStock.symbol}
          name={selectedStock.name}
          sector={selectedStock.sector}
          onClose={() =>
            setSelectedStock(null)
          }
        />
      )}
    </>
  );
}