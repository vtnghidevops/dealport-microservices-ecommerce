import AppRouters from '@/routers/routers';
import { AppProvider } from '@/providers/AppProvider';
import { BrowserRouter } from 'react-router-dom';
import { SnackbarProvider } from 'notistack';
import { AuthProvider } from '@/context/AuthContext';
import { Toaster } from '@/components/ui/toaster';
import '@/index.css'

function App() {
  return (
    <>
      <AuthProvider>
        <div className="w-full max-w-[1440px] mx-auto px-4 min-w-[1440px]">
          <SnackbarProvider>
            <BrowserRouter>
              <AppProvider>
                <AppRouters />
                <Toaster />
              </AppProvider>
            </BrowserRouter>
          </SnackbarProvider>
        </div>
      </AuthProvider>
    </>
  );
}

export default App
