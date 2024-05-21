import { useEffect, useState } from "react";
import { fetchToggles } from "../api/toggles/toggles";
import { Toggle } from "../types/toggles";
import ToggleItem from "../components/toggles/toggleItem";
import ToggleHeader from "@/components/toggles/toggleHeader";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";

export default function Dashboard() {
  const navigate = useNavigate();
  const [toggles, setToggles] = useState<Toggle[]>();

  const getToggles = async () => {
    const response = await fetchToggles();
    if (response) {
      setToggles(response.data);
    }
  };

  useEffect(() => {
    getToggles();
  }, []);

  const handleAddToggle = () => {
    navigate("/create");
  };

  const handleDeleteToggleDash = async (toggleId: string) => {
    const newToggles = toggles?.filter((toggle) => toggle.id !== toggleId);
    setToggles(newToggles);
  }



  return (
    <div className="">
      <div className="px-10 flex items-center justify-start bg-gradient-to-l from-slate-700 to-slate-800 text-slate-100 h-24 font-medium">
        <h1 className="text-2xl">Simple feature toggles</h1>
      </div>

      <div className="flex flex-col items-center">
        <div className="px-8 pt-6 lg:w-[1000px]">
          <div className="hidden lg:block"><ToggleHeader/></div>
          {toggles &&
            toggles.map((toggle) => {
              return <ToggleItem key={toggle.id} toggle={toggle} deleteToggleDash={handleDeleteToggleDash}/>;
            })}
        </div>
        <Button onClick={handleAddToggle} className="w-48 h-20 text-2xl text-slate-100 fixed bottom-10 right-10 lg:right-20 ">
          Add toggle
        </Button>
      </div>
    </div>
  );
}
