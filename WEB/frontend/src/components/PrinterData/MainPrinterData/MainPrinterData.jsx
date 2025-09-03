export default function MainPrinterData({PrinterData}) {
    return (<div className="MainPrinterData">
        <p><strong>Type:</strong> <span id="detailType">{PrinterData.current.Type}</span></p>
        <p><strong>Version:</strong> <span id="detailVersion">{PrinterData.current.Version}</span></p>
        <p><strong>Working:</strong> <span id="detailWorking">{PrinterData.current.Working === true ? "Yes" : "No"}</span></p>
        <p><strong>Connection:</strong> <span id="detailConnection">{PrinterData.current.Connection}</span></p>
    </div>)
}