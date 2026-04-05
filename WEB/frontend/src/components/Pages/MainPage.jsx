import { useState } from 'react'
import Header from '../Header'
import PrinterData from '../PrinterData/PrinterData'
import MachinesRegistry from '../MachinesRegistry/MachinesRegistry'
import ReportsPlaceholder from '../ReportsPlaceholder/ReportsPlaceholder'

export default function MainPage() {
  const [activeNav, setActiveNav] = useState('home')

  return (
    <div className="min-h-screen">
      <Header activeNav={activeNav} onNavigate={setActiveNav} />
      <main className="w-full px-6 lg:px-10 py-6">
        {activeNav === 'home' && <PrinterData />}
        {activeNav === 'machines' && <MachinesRegistry />}
        {activeNav === 'reports' && <ReportsPlaceholder />}
      </main>
    </div>
  )
}