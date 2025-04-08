export type ButtonType = 'primary' | 'primary-cy' | 'secondary' | 'accent' | 'gradient' | 'gray' | 'black';

export const getButtonClass = (type: ButtonType): string => {
  switch (type) {
    case "primary":
      return "button-text btn-primary text-cyprus hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:hover:text-black";
    case "primary-cy":
      return "button-text btn-primary-cy hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:text-black";
    case "secondary":
      return "button-text btn-secondary text-cyprus hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:text-black";
    case "accent":
      return "button-text btn-accent text-cyprus hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:text-black";
    case "gradient":
      return "button-text btn-gradient text-cyprus hover:shadow-md hover:scale-[1.02] transition-all duration-300 ease-out";
    case "gray":
      return "buttonText btn-gray border border-white hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:text-black";
    case "black":
      return "buttonText btn-black border border-white-200 hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out hover:text-black";
    default:
      return "buttonText border border-black-200 hover:bg-surf-crest hover:shadow-md transition-all duration-300 ease-out";
  }
};