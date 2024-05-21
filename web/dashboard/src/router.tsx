import { createBrowserRouter } from "react-router-dom";
import Dashboard from "./pages/dashboard";
import Welcome from "./pages/welcome";
import Register from "./pages/register";
import CreateToggle from "./pages/createToggle";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Welcome />,
  },
  {
    path: "/register",
    element: <Register />,
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
