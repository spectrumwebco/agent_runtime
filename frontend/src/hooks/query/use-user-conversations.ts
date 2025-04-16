import { useQuery } from "@tanstack/react-query";
import KledIo from "#/api/kled-io";
import { useIsAuthed } from "./use-is-authed";

export const useUserConversations = () => {
  const { data: userIsAuthenticated } = useIsAuthed();

  return useQuery({
    queryKey: ["user", "conversations"],
    queryFn: Kled.getUserConversations,
    enabled: !!userIsAuthenticated,
  });
};
