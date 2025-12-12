import { useLocalization } from '../../../hooks/useLocalization.jsx'
import './gabarites.css'
import '../PrinterDetails/PrinterDetails.css'

export default function PrinterGabarites({ PrinterDataRef }) {
  const { t } = useLocalization()
  
  // Поддержка как ref, так и обычного объекта
  const printerData = PrinterDataRef?.current || PrinterDataRef || {}
  
  // Функция для форматирования числовых значений
  const formatValue = (value) => {
    if (value === null || value === undefined) return '-'
    const numValue = typeof value === 'number' ? value : parseFloat(value)
    if (isNaN(numValue)) return value
    return numValue.toFixed(2)
  }
  
  // Получаем габариты (проверяем разные варианты структуры данных)
  const widthRaw = printerData.immutable.width
  const lengthRaw =  printerData.immutable.length 
  const heightRaw =  printerData.immutable.height
  const width = formatValue(widthRaw)
  const length = formatValue(lengthRaw)
  const height = formatValue(heightRaw)
  
  // Получаем позицию (проверяем разные варианты структуры данных)
  const posXRaw = printerData.position?.X ?? printerData.position?.x ?? printerData.X ?? printerData.x
  const posYRaw = printerData.position?.Y ?? printerData.position?.y ?? printerData.Y ?? printerData.y
  const posZRaw = printerData.position?.Z ?? printerData.position?.z ?? printerData.Z ?? printerData.z
  const posX = formatValue(posXRaw)
  const posY = formatValue(posYRaw)
  const posZ = formatValue(posZRaw)
  
  // // Получаем температуру и статус печати
  // const nozzleTempRaw = printerData.NozzleTemp ?? printerData.nozzleTemp ?? printerData.NowTempNozzle ?? printerData.nowTempNozzle
  // const bedTempRaw = printerData.BedTemp ?? printerData.bedTemp ?? printerData.NowTempBed ?? printerData.nowTempBed
  // const nozzleTemp = nozzleTempRaw !== null && nozzleTempRaw !== undefined ? `${formatValue(nozzleTempRaw)}°C` : '-'
  // const bedTemp = bedTempRaw !== null && bedTempRaw !== undefined ? `${formatValue(bedTempRaw)}°C` : '-'
  // const isPrinting = printerData.executingTask === true || printerData.Flags?.ExecutingTask === true
  
  return (
    <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
      <div className="PrinterGabarites">
        <p><strong>{t('printers.details.dimensions.width')}:</strong> <span id="PrinterWidth">{width}</span></p>
        <p><strong>{t('printers.details.dimensions.length')}:</strong> <span id="PrinterLength">{length}</span></p>
        <p><strong>{t('printers.details.dimensions.height')}:</strong> <span id="PrinterHeight">{height}</span></p>
      </div>
      <div className="PrinterPosotion">
        <p><strong>{t('printers.details.position.x')}:</strong> <span id="XPos">{posX}</span></p>
        <p><strong>{t('printers.details.position.y')}:</strong> <span id="YPos">{posY}</span></p>
        <p><strong>{t('printers.details.position.z')}:</strong> <span id="ZPos">{posZ}</span></p>
      </div>
{/* 
      <div className="PrinterState">
        <p><strong>{t('printers.details.state.nozzleTemp')}:</strong> <span id="NozzleTemp">{nozzleTemp}</span></p>
        <p><strong>{t('printers.details.state.bedTemp')}:</strong> <span id="BedTemp">{bedTemp}</span></p>
        <p><strong>{t('printers.details.state.isPrinting')}:</strong> <span id="IsPrinting">{isPrinting ? t('common.yes') : t('common.no')}</span></p>
      </div> */}
    </div>
  )
}