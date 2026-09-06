export interface Instrument {
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

export interface Quote {
  instrumentId: string;
  symbol: string;
  price: number;
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
  reliabilityState:
    | "LIVE"
    | "DELAYED"
    | "STALE"
    | "UNAVAILABLE"
    | "CONFLICTING";
}

export interface WatchItem {
  instrumentId: string;
  position: number;
  addedAt: string;
}

export interface Watchlist {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  stocks: WatchItem[];
}

export interface WatchlistStock {
  instrumentId: string;
  symbol: string;
  name: string;
  exchange: string;
  sector?: string | null;
  industry?: string | null;
  position: number;
  addedAt: string;
  price?: number | null;
  previousClose?: number | null;
  change?: number | null;
  changePercent?: number | null;
  volume?: number | null;
  reliabilityState?:
    | "LIVE"
    | "DELAYED"
    | "STALE"
    | "UNAVAILABLE"
    | "CONFLICTING"
    | null;
}

export interface MarketChange {
  id?: string | null;
  symbol: string;
  name?: string | null;
  previousPrice?: number | null;
  currentPrice?: number | null;
  priceChange?: number | null;
  priceChangePercent?: number | null;
  previousVolume?: number | null;
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
  reliabilityState?:
    | "LIVE"
    | "DELAYED"
    | "STALE"
    | "UNAVAILABLE"
    | "CONFLICTING"
    | null;
}

export interface ChangeSummary {
  watchlistId: string;
  changes: MarketChange[];
  checkpoint?: {
    id: string;
    createdAt: string;
  } | null;
}

export interface User {
  id: string;
  name: string;
  email: string;
}

export interface AuthPayload {
  token: string;
  user: User;
}

export interface Checkpoint {
  id: string;
  watchlistId: string;
  createdAt: string;
}