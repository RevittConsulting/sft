import React from "react";
import { Input } from "../ui/input";
import { X } from "lucide-react";
import { Button } from "../ui/button";

interface MetaDataPair {
  key: string;
  value: string;
}

interface ToggleMetaInputProps {
  metaData: MetaDataPair[];
  setMetaData: React.Dispatch<React.SetStateAction<MetaDataPair[]>>;
}

const ToggleMetaInput: React.FC<ToggleMetaInputProps> = ({
  metaData,
  setMetaData,
}) => {
  const handleKeyChange = (
    index: number,
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const newMetaData = metaData.map((item, i) => {
      if (index === i) {
        return { ...item, key: event.target.value };
      }
      return item;
    });
    setMetaData(newMetaData);
  };

  const handleValueChange = (
    index: number,
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const newMetaData = metaData.map((item, i) => {
      if (index === i) {
        return { ...item, value: event.target.value };
      }
      return item;
    });
    setMetaData(newMetaData);
  };

  const handleAddPair = () => {
    setMetaData([...metaData, { key: "", value: "" }]);
  };

  const handleRemovePair = (index: number) => {
    setMetaData(metaData.filter((_, i) => i !== index));
  };

  return (
    <>
      {metaData.map((item: any, index: any) => (
        <div key={index} className="mt-4">
          <div className="flex flex-row gap-2">
            <Input
              type="text"
              placeholder="Key"
              value={item.key}
              onChange={(e) => handleKeyChange(index, e)}
              className=""
            />
            <Input
              type="text"
              placeholder="Value"
              value={item.value}
              onChange={(e) => handleValueChange(index, e)}
              className=""
            />
            <button type="button" onClick={() => handleRemovePair(index)}>
              <X />
            </button>
          </div>
        </div>
      ))}
      <Button type="button" onClick={handleAddPair} variant="outline" className="mt-2">
        Add another pair
      </Button>
    </>
  );
};

export default ToggleMetaInput;
