import { Toggle } from "@/types/toggles";
import { Switch } from "@/components/ui/switch";
import { useEffect, useState } from "react";
import { toggleFeature } from "@/api/toggles/toggles";
import { Axe } from "lucide-react";
import { deleteToggle } from "@/api/toggles/toggles";
import { DeleteToggleDialog } from "./dialogs/deleteToggleDialog";

interface ToggleItemProps {
  toggle: Toggle;
  deleteToggleDash: (toggleId: string) => void;
}

const ToggleItem: React.FC<ToggleItemProps> = ({toggle, deleteToggleDash}) => {
  const [localToggle, setLocalToggle] = useState<Toggle>(toggle);

  useEffect(() => {
    setLocalToggle(toggle);
  }, [toggle]);

  const handleToggle = async () => {
    const newEnabledState = !localToggle.enabled;
    setLocalToggle((prevLocalToggle) => ({ ...prevLocalToggle, enabled: newEnabledState }));
    try {
      const response = await toggleFeature(localToggle.id);
      if (response) {
        console.log("feature toggle toggled!");
      }
    } catch (error) {
      setLocalToggle((prevLocalToggle) => ({ ...prevLocalToggle, enabled: !newEnabledState }));
      console.log("Error toggling feature:", error);
    }
  };

  const handleDeleteToggle = async () => {
    console.log("trying to delete");
    try {
      const response = await deleteToggle(localToggle.id);
      if (response) { 
        console.log(response);
      }
      deleteToggleDash(localToggle.id);
    } catch (error) {
      console.log("Error deleting feature:", error);
    }
  };

  // Dialog box
  const [openDialogBox, setOpenDialogBox] = useState(false);

  return (
    <>
      <div className="border-2 border-slate-400 rounded-lg flex flex-col gap-2 lg:gap-0 lg:flex-row lg:flex-wrap justify-between py-2 px-4 my-2">
        
          <h2 className="mr-4 font-medium text-lg text-left lg:w-[20%]">
            {toggle.feature_name}
          </h2>
        
        <div className="flex flex-col text-left lg:w-[50%] ml-4 lg:ml-0">
          {" "}
          {toggle.toggle_meta &&
            Object.keys(toggle.toggle_meta).map((key) => {
              return (
                <ul key={key} className="list-disc">
                  <li>
                    {key}: {String(toggle.toggle_meta[key])}
                  </li>
                </ul>
              );
            })}
        </div>
        <div className="flex flex-wrap w-[240px] justify-between">
          <div className="">
            <p className="">{localToggle.enabled ? "Enabled" : "Disabled"}</p>
          </div>
          <div className="">
            <Switch
              checked={localToggle.enabled}
              onCheckedChange={() => {
                handleToggle();
              }}
            />
          </div>
          <div className="">
            <button
              onClick={() => {
                setOpenDialogBox(true);
              }}
            >
              <Axe />
            </button>
          </div>
        </div>
      </div>
      <DeleteToggleDialog 
        isOpen={openDialogBox} 
        setIsOpen={setOpenDialogBox} 
        deleteToggle={() => handleDeleteToggle()}
      />
    </>
  );
};

export default ToggleItem;
