import { gql } from "@apollo/client";
import { useMutation, useQuery } from "@apollo/client/react";
import { Check, Loader2, Plus, Search } from "lucide-react";
import { useState } from "react";

const SEARCH_INSTRUMENTS = gql`
  query SearchInstruments($query: String!) {
    searchInstruments(query: $query) {
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
  }
`;

const ADD_STOCK = gql`
  mutation AddStock(
    $watchlistId: ID!
    $instrumentId: ID!
  ) {
    addStock(
      watchlistId: $watchlistId
      instrumentId: $instrumentId
    )
  }
`;

interface Instrument {
  id: string;
  symbol: string;
  name: string;
  exchange: string;
  segment: string;
  isin: string;
  currency: string;
  sector?: string | null;
  industry?: string | null;
}

interface StockSearchProps {
  watchlistId: string;
  existingStockIds: string[];
  onStockAdded: () => void;
}

export default function StockSearch({
  watchlistId,
  existingStockIds,
  onStockAdded,
}: StockSearchProps) {
  const [search, setSearch] = useState("");
  const [addedIds, setAddedIds] = useState<string[]>([]);

  const { data, loading, error } = useQuery<{
    searchInstruments: Instrument[];
  }>(SEARCH_INSTRUMENTS, {
    variables: { query: search.trim() },
    skip: search.trim().length < 2,
    fetchPolicy: "network-only",
  });

  const [addStock, { loading: adding }] = useMutation(ADD_STOCK);

  async function handleAdd(instrumentId: string) {
    try {
      await addStock({
        variables: { watchlistId, instrumentId },
      });
      setAddedIds((current) => [...current, instrumentId]);
      onStockAdded();
    } catch (err) {
      console.error(err);
      alert("Unable to add this stock.");
    }
  }

  const results = data?.searchInstruments ?? [];

  return (
    <div className="stock-search-panel">
      <div className="stock-search-title">
        <div>
          <strong>Add Stocks</strong>
          <span>Search company or symbol to monitor</span>
        </div>
      </div>

      <div className="stock-search-input">
        <Search size={16} />
        <input
          type="text"
          value={search}
          placeholder="Search stocks (e.g. INFY, RELIANCE)..."
          onChange={(event) => setSearch(event.target.value)}
        />
        {loading && <Loader2 size={16} className="spin search-spinner" />}
      </div>

      {search.trim().length < 2 && (
        <div className="search-hint">Type at least 2 characters to search</div>
      )}

      {error && <div className="search-error">Unable to search stocks.</div>}

      {!loading && !error && search.trim().length >= 2 && results.length === 0 && (
        <div className="search-empty">No matching stocks found.</div>
      )}

      {results.length > 0 && (
        <div className="search-results">
          {results.map((instrument) => {
            const isAdded =
              existingStockIds.includes(instrument.id) ||
              addedIds.includes(instrument.id);

            return (
              <div className="search-result" key={instrument.id}>
                <div className="search-result-left">
                  <div className="search-avatar">
                    {instrument.symbol.slice(0, 3)}
                  </div>
                  <div className="search-result-info">
                    <div className="search-symbol-row">
                      <span className="search-symbol">{instrument.symbol}</span>
                      <span className="search-exchange">{instrument.exchange}</span>
                    </div>
                    <div className="search-company">{instrument.name}</div>
                    {instrument.sector && (
                      <div className="search-meta">{instrument.sector}</div>
                    )}
                  </div>
                </div>

                {isAdded ? (
                  <button type="button" className="stock-added-button" disabled>
                    <Check size={14} />
                    Added
                  </button>
                ) : (
                  <button
                    type="button"
                    className="stock-add-button"
                    disabled={adding}
                    onClick={() => handleAdd(instrument.id)}
                  >
                    {adding ? (
                      <Loader2 size={14} className="spin" />
                    ) : (
                      <Plus size={14} />
                    )}
                    {adding ? "Adding..." : "Add"}
                  </button>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}