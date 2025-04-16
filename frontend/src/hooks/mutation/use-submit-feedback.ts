import { useMutation } from "@tanstack/react-query";
import { Feedback } from "#/api/kled-io.types";
import KledIo from "#/api/kled-io";
import { useConversation } from "#/context/conversation-context";
import { displayErrorToast } from "#/utils/custom-toast-handlers";

type SubmitFeedbackArgs = {
  feedback: Feedback;
};

export const useSubmitFeedback = () => {
  const { conversationId } = useConversation();
  return useMutation({
    mutationFn: ({ feedback }: SubmitFeedbackArgs) =>
      KledIo.submitFeedback(conversationId, feedback),
    onError: (error) => {
      displayErrorToast(error.message);
    },
  });
};
