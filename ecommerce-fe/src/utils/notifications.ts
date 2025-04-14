import { enqueueSnackbar, VariantType } from 'notistack';
export const showNotification = (message: string, variant: VariantType = 'success') => {
  enqueueSnackbar(message, {
    variant,
    autoHideDuration: 2000,
    anchorOrigin: {
      vertical: 'bottom',
      horizontal: 'right'
    },
    preventDuplicate: true
  });
};