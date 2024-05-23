import ToggleMetaInput from "@/components/createToggle/toggleMetaInput";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { ToggleDto } from "@/types/toggles";
import { ArrowLeft } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createToggle } from "@/api/toggles/toggles";

export default function CreateToggle() {
  // navigation
  const navigate = useNavigate();
  const handleGoBack = () => {
    navigate(-1);
  };

  const [featureName, setFeatureName] = useState<string>("");
  const [enabled, setEnabled] = useState<boolean>(false);

  // toggle_meta
  const [metaData, setMetaData] = useState([{ key: "", value: "" }]);

  // submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // create toggle_meta object from metaData, ignoring blank entries
    const formattedMetaData: Record<string, string> = metaData.reduce(
      (acc: Record<string, string>, { key, value }) => {
        if (key != "") {
          acc[key] = value;
        }
        return acc;
      },
      {}
    );

    const payload: ToggleDto = {
      feature_name: featureName,
      enabled,
      toggle_meta: formattedMetaData,
    };

    try {
      const response = await createToggle(payload);
      if (response) {
        console.log(response);
        navigate("/dashboard");
      }
    } catch (err) {
      console.error("Error creating toggle:", err);
    }

    console.log("submit");
  };

  useEffect(() => {
    console.log(enabled);
  }, [enabled]);

  return (
    <form
      className="flex flex-col h-full w-full"
      onSubmit={(e) => {
        handleSubmit(e);
      }}
    >
      <div className="px-2 gap-2 h-16 flex items-center justify-start bg-gradient-to-l from-slate-800 to-slate-950 text-center text-slate-100 py-1 font-medium">
        <Button variant="ghost" onClick={handleGoBack}>
          <ArrowLeft />
        </Button>
        <h1 className="text-2xl">Create new toggle</h1>
      </div>

      <div className="flex justify-center">
        <div className="flex flex-col gap-4 mx-6 mt-10 w-[600px] sm:border-2 sm:px-4 sm:py-8 sm:rounded-lg">
          <div className="">
            <Label className="text-xl">Feature name</Label>
            <Input
              placeholder="Feature name"
              onChange={(e) => {
                setFeatureName(e.target.value);
              }}
              className="mt-4"
            />
          </div>
          <div>
            <Label className="text-xl">Toggle meta</Label>
            <ToggleMetaInput metaData={metaData} setMetaData={setMetaData} />
          </div>
          <RadioGroup
            defaultValue="false"
            onValueChange={(e) => {
              setEnabled(e === "true");
            }}
          >
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="true" id="enabled" />
              <Label htmlFor="enabled">Enable</Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="false" id="disable" />
              <Label htmlFor="disable">Disable</Label>
            </div>
          </RadioGroup>
          <div className="flex flex-row gap-4 justify-center">
            <Button variant="outline" onClick={handleGoBack} className="w-32">
              Cancel
            </Button>
            <Button>Create toggle</Button>
          </div>
        </div>
      </div>
    </form>
  );
}
