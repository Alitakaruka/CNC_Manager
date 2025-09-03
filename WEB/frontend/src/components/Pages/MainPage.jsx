import Header from '../Header'
import PrinterData from '../PrinterData/PrinterData'

export default function MainPage() {
  return (
    <div className="min-h-screen">
      <Header />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div>
          <PrinterData />
        </div>
      </main>
    </div>
  )
}