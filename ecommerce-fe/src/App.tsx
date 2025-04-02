import AppRouters from './routers/routers';
import { SnackbarProvider } from "notistack";
import './index.css'
function App() {
  return (
    <>
    <div className='w-full max-w-[1440px] mx-auto px-4 min-w-[1440px]'>
      <SnackbarProvider>
          <AppRouters></AppRouters>
      </SnackbarProvider>
    </div>
    </>
  )
}

export default App
