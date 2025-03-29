import { toast } from 'react-toastify'; // Ensure correct import from react-toastify

export const showSuccess = (message: string) => {
  toast.success(message);
};

export const showError = (message: string) => {
  toast.error(message);
};

