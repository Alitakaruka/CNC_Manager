import Header from '../Header'
import PrinterData from '../PrinterData/PrinterData'

export default function MainPage() {
  return (
    <div className="min-h-screen">
      <Header />
      <main className="w-full px-6 lg:px-10 py-6">
        <div>
          <PrinterData />
        </div>
      </main>
    </div>
  )
}