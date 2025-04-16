import axios from "axios";

export const kledIo = axios.create({
  baseURL: `${window.location.protocol}//${import.meta.env.RSBUILD_BACKEND_BASE_URL || window?.location.host}`,
});
