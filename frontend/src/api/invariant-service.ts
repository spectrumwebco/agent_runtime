import { kledIo } from "./kled-io-axios";

class InvariantService {
  static async getPolicy() {
    const { data } = await kledIo.get("/api/security/policy");
    return data.policy;
  }

  static async getRiskSeverity() {
    const { data } = await kledIo.get("/api/security/settings");
    return data.RISK_SEVERITY;
  }

  static async getTraces() {
    const { data } = await kledIo.get("/api/security/export-trace");
    return data;
  }

  static async updatePolicy(policy: string) {
    await kledIo.post("/api/security/policy", { policy });
  }

  static async updateRiskSeverity(riskSeverity: number) {
    await kledIo.post("/api/security/settings", {
      RISK_SEVERITY: riskSeverity,
    });
  }
}

export default InvariantService;
