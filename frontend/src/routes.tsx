import React from "react";
import { createBrowserRouter } from "react-router";
import MainApp from "./routes/root-layout";
import KledAppRoute from "./routes/kled-app-route";
import { ErrorBoundary } from "./routes/root-layout";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <MainApp />,
    errorElement: <ErrorBoundary />,
  },
  {
    path: "/kled",
    element: <KledAppRoute />,
    errorElement: <ErrorBoundary />,
  },
]);

export default router;
