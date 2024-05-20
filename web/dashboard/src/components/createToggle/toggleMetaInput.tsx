import React from "react";


interface MetaDataPair {
    key: string;
    value: string;
}

interface ToggleMetaInputProps {
    metaData: MetaDataPair[];
    setMetaData: React.Dispatch<React.SetStateAction<MetaDataPair[]>>;
}


const ToggleMetaInput: React.FC<ToggleMetaInputProps> = ({ metaData, setMetaData }) => {
  const handleKeyChange = (index: number, event: React.ChangeEvent<HTMLInputElement>) => {
    const newMetaData = metaData.map((item, i) => {
      if (index === i) {
        return { ...item, key: event.target.value };
      }
      return item;
    });
    setMetaData(newMetaData);
  };

  const handleValueChange = (index: number, event: React.ChangeEvent<HTMLInputElement>) => {
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
            <div key={index}>
                <input
                    type="text"
                    placeholder="Key"
                    value={item.key}
                    onChange={(e) => handleKeyChange(index, e)}
                    className="border-2 rounded p-2 m-2"
                />
                <input
                    type="text"
                    placeholder="Value"
                    value={item.value}
                    onChange={(e) => handleValueChange(index, e)}
                    className="border-2 rounded p-2 m-2"
                />
                <button type="button" onClick={() => handleRemovePair(index)}>
                    Remove
                </button>
            </div>
        ))}
        <button type="button" onClick={handleAddPair}>
            Add More
        </button>
    </>
);
};

export default ToggleMetaInput;
