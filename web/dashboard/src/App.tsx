
import { RouterProvider } from 'react-router-dom'
import { router } from './router'
import './App.css'

function App() {

  return (
      <main className="h-screen">
         <div className="px-10 flex items-center justify-start bg-gradient-to-l from-slate-700 to-slate-800 text-slate-100 h-24 font-medium">
        <h1 className="text-2xl">Simple feature toggles</h1>
      </div>
        <RouterProvider router={router} />
      </main>


  )
}

export default App
