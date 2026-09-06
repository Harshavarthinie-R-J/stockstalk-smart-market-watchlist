import {
  Activity,
  BarChart3,
  CheckCircle2,
  TrendingDown,
  TrendingUp,
} from "lucide-react";

interface PersonalBaselineProps {
  symbol: string;
  changePercent?: number | null;
  currentVolume?: number | null;
  averageVolume?: number | null;
}

export default function PersonalBaseline({
  symbol,
  changePercent,
  currentVolume,
  averageVolume,
}: PersonalBaselineProps) {
  const change = changePercent ?? 0;

  const volumeRatio =
    currentVolume != null &&
    averageVolume != null &&
    averageVolume > 0
      ? currentVolume / averageVolume
      : null;

  const priceUnusual = Math.abs(change) >= 1.5;
  const volumeUnusual =
    volumeRatio != null && volumeRatio >= 1.5;

  const unusual =
    priceUnusual || volumeUnusual;

  return (
    <div className="personal-baseline-card">
      <div className="personal-baseline-header">
        <div className="personal-baseline-icon">
          <BarChart3 size={18} />
        </div>

        <div>
          <strong>
            Personal baseline
          </strong>

          <span>
            Is {symbol} behaving unusually?
          </span>
        </div>
      </div>

      <div className="personal-baseline-grid">
        <div className="baseline-metric">
          <span>Today's movement</span>

          <strong
            className={
              change >= 0
                ? "positive-value"
                : "negative-value"
            }
          >
            {change >= 0 ? "+" : ""}
            {change.toFixed(2)}%
          </strong>

          <small>
            {priceUnusual
              ? "Above normal movement"
              : "Within normal movement"}
          </small>
        </div>

        <div className="baseline-metric">
          <span>Volume</span>

          <strong>
            {volumeRatio != null
              ? `${volumeRatio.toFixed(2)}×`
              : "—"}
          </strong>

          <small>
            {volumeUnusual
              ? "Above historical baseline"
              : "Near historical baseline"}
          </small>
        </div>
      </div>

      <div
        className={
          unusual
            ? "baseline-result unusual"
            : "baseline-result normal"
        }
      >
        {unusual ? (
          <>
            <Activity size={18} />

            <div>
              <strong>
                Unusual for {symbol}
              </strong>

              <span>
                The current behaviour is
                outside the normal pattern
                detected by StockStalk.
              </span>
            </div>
          </>
        ) : (
          <>
            <CheckCircle2 size={18} />

            <div>
              <strong>
                Within normal range
              </strong>

              <span>
                No significant deviation from
                the available baseline.
              </span>
            </div>
          </>
        )}
      </div>
    </div>
  );
}