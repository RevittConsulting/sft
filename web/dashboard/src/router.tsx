import { createBrowserRouter } from "react-router-dom";
import Dashboard from "./pages/dashboard";
import Welcome from "./pages/welcome";
import CreateToggle from "./pages/createToggle";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Welcome />,
  },
  {
    path: "/dashboard",
    element: <Dashboard />,
  },
  {
    path: "/create",
    element: <CreateToggle />,
  },
]);
