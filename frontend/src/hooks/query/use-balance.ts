import { useQuery } from "@tanstack/react-query";
import { useConfig } from "./use-config";
import KledIo from "#/api/kled-io";

export const useBalance = () => {
  const { data: config } = useConfig();

  return useQuery({
    queryKey: ["user", "balance"],
    queryFn: Kledio.getBalance,
    enabled:
      config?.APP_MODE === "saas" && config?.FEATURE_FLAGS.ENABLE_BILLING,
  });
};
