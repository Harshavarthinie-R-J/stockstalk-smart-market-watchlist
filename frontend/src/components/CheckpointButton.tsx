import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { Check, Loader2 } from "lucide-react";

const CREATE_CHECKPOINT = gql`
  mutation CreateCheckpoint($watchlistId: ID!) {
    createCheckpoint(watchlistId: $watchlistId) {
      id
      watchlistId
      createdAt
    }
  }
`;

interface CheckpointButtonProps {
  watchlistId: string;
  onCheckpointCreated?: () => void;
}

export default function CheckpointButton({
  watchlistId,
  onCheckpointCreated,
}: CheckpointButtonProps) {
  const [createCheckpoint, { loading }] = useMutation(CREATE_CHECKPOINT);

  async function handleCheckpoint() {
    try {
      await createCheckpoint({
        variables: {
          watchlistId,
        },
      });
      onCheckpointCreated?.();
    } catch (error) {
      console.error("Unable to create checkpoint:", error);
      alert("Unable to mark this watchlist as checked.");
    }
  }

  return (
    <button
      type="button"
      className="checkpoint-button"
      onClick={handleCheckpoint}
      disabled={loading}
    >
      {loading ? (
        <>
          <Loader2 size={14} className="spin" />
          <span>Saving...</span>
        </>
      ) : (
        <>
          <Check size={14} />
          <span>Mark as checked</span>
        </>
      )}
    </button>
  );
}