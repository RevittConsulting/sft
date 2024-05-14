

const ToggleHeader = () => {
 
  return (
    <div className="border-2 border-slate-400 rounded-lg flex flex-wrap justify-between py-2 px-4 my-2">
      <div className="w-[20%]">
        <h2 className="mr-4 font-medium text-lg text-left">
          Feature
        </h2>
      </div>
      <div className="flex flex-col text-left w-[50%]">
        <p className="font-medium text-lg text-left">Toggle meta</p>
      </div>
      <div className="w-[8%]">
        <p className="font-medium text-lg ">Enabled</p>
      </div>
      <div className="w-[8%]">
        <p className="font-medium text-lg ">Toggle</p>
      </div>
      <div className="w-[8%]">
        <p className="font-medium text-lg ">Delete</p>
      </div>
    </div>
  );
};

export default ToggleHeader;
