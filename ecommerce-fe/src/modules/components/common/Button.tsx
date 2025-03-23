import React from 'react';

const getButtonClass = (type: string) => {
  switch(type){
    case "primary":
      return "button-text btn-primary text-cyprus"
    case "primary-cy":
      return "button-text btn-primary-cy"
    case "secondary":
      return "button-text btn-secondary text-cyprus"
    case "accent":
      return "button-text btn-accent text-cyprus"
    case "gradient":
      return "button-text btn-gradient text-cyprus"
    case "gray":
      return "buttonText btn-gray border border-white"
    case "black":
      return "buttonText btn-black border border-white-200"
    default:
      // view all
      return "buttonText border border-black-200"
  }
}
export function Button({type, text}) {
  return (
    <div className="w-[12rem] h-[4rem] rounded-3xl ">
      <button className={getButtonClass(type)}>{text}</button>
    </div>
  );
}
export { getButtonClass };