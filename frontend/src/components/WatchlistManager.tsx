import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { Check, Loader2, Plus, X, Layers } from "lucide-react";
import { useState } from "react";

const CREATE_WATCHLIST = gql`
  mutation CreateWatchlist($name: String!) {
    createWatchlist(name: $name) {
      id
      name
    }
  }
`;

interface Watchlist {
  id: string;
  name: string;
  stocks: unknown[];
}

interface CreateWatchlistData {
  createWatchlist: {
    id: string;
    name: string;
  };
}

interface WatchlistManagerProps {
  watchlists: Watchlist[];
  activeWatchlistId?: string;
  onSelect: (watchlistId: string) => void;
  onChanged: () => Promise<unknown> | void;
}

export default function WatchlistManager({
  watchlists,
  activeWatchlistId,
  onSelect,
  onChanged,
}: WatchlistManagerProps) {
  const [showCreate, setShowCreate] = useState(false);
  const [newName, setNewName] = useState("");

  const [createWatchlist, { loading: creating }] = useMutation<
    CreateWatchlistData,
    { name: string }
  >(CREATE_WATCHLIST);

  async function handleCreate() {
    const name = newName.trim();
    if (!name) return;

    try {
      const result = await createWatchlist({
        variables: { name },
      });

      const created = result.data?.createWatchlist;
      if (!created?.id) {
        throw new Error("Backend did not return a watchlist ID.");
      }

      setNewName("");
      setShowCreate(false);
      onSelect(created.id);

      try {
        await onChanged();
      } catch (refreshError) {
        console.error("Watchlist created, but refresh failed:", refreshError);
      }
    } catch (error) {
      console.error("CREATE WATCHLIST ERROR:", error);
      const message = error instanceof Error ? error.message : "Unknown error";
      alert(`Unable to create watchlist.\n\n${message}`);
    }
  }

  function handleCancelCreate() {
    setShowCreate(false);
    setNewName("");
  }

  return (
    <div className="watchlist-manager">
      <div className="watchlist-manager-header">
        <div>
          <strong>Watchlists</strong>
          <span>Organize the stocks you follow</span>
        </div>

        <button
          type="button"
          className="watchlist-create-button"
          onClick={() => setShowCreate((current) => !current)}
        >
          <Plus size={14} />
          New
        </button>
      </div>

      {showCreate && (
        <div className="watchlist-create-form">
          <input
            type="text"
            value={newName}
            placeholder="Watchlist name"
            autoFocus
            maxLength={50}
            disabled={creating}
            onChange={(event) => setNewName(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter") handleCreate();
              if (event.key === "Escape") handleCancelCreate();
            }}
          />

          <button
            type="button"
            className="watchlist-form-confirm"
            onClick={handleCreate}
            disabled={creating || !newName.trim()}
            title="Create watchlist"
          >
            {creating ? (
              <Loader2 size={13} className="spin" />
            ) : (
              <Check size={13} />
            )}
          </button>

          <button
            type="button"
            className="watchlist-form-cancel"
            onClick={handleCancelCreate}
            disabled={creating}
            title="Cancel"
          >
            <X size={13} />
          </button>
        </div>
      )}

      <div className="watchlist-manager-list">
        {watchlists.map((watchlist) => {
          const active = watchlist.id === activeWatchlistId;

          return (
            <div
              key={watchlist.id}
              className={`watchlist-manager-item ${
                active ? "watchlist-manager-item-active" : ""
              }`}
            >
              <button
                type="button"
                className="watchlist-select-button"
                onClick={() => onSelect(watchlist.id)}
              >
                <div className="watchlist-item-left">
                  <span className="watchlist-active-bar" />
                  <span className="watchlist-item-name">{watchlist.name}</span>
                </div>
                <span className="watchlist-item-count">
                  {watchlist.stocks.length}
                </span>
              </button>
            </div>
          );
        })}

        {watchlists.length === 0 && (
          <div className="watchlist-manager-empty">
            <Layers size={24} />
            <span>No watchlists yet</span>
          </div>
        )}
      </div>
    </div>
  );
}