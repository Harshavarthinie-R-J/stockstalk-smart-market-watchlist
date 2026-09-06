import React from "react";

interface AttentionBadgeProps {
  score?: number | null;
}

export default function AttentionBadge({ score }: AttentionBadgeProps) {
  if (score === null || score === undefined) {
    return (
      <span className="attention-badge attention-unknown">
        No score
      </span>
    );
  }

  let label = "Low";
  let className = "attention-low";

  if (score >= 75) {
    label = "High";
    className = "attention-high";
  } else if (score >= 45) {
    label = "Medium";
    className = "attention-medium";
  }

  return (
    <span className={`attention-badge ${className}`}>
      <span className="attention-dot" />
      <span className="attention-label">{label}</span>
      <span className="attention-score">{Math.round(score)}</span>
    </span>
  );
}