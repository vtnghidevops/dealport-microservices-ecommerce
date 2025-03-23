import Reactq from 'react'
import Header from './modules/components/layouts/Header'
import Home from './modules/pages/dashboard/Home'
import Footer from './modules/components/layouts/Footer'
import './index.css'
function App() {
  return (
    <div className='w-full max-w-[1440px] mx-auto px-4 min-w-[1440px]'>
      <Header></Header>
      <Home></Home>
      <Footer></Footer>
    </div>
  )
}

export default App
