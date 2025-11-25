
 export default function PrinterGabarites ({PrinterDataRef}) {
    return (<div style={{ display: "flex" }}>
      <div className="PrinterGabarites">
        {/* <div style={{display:"flex", justifyContent:"center",margin:"5px"}}>
          <h4>Gabarites</h4>
        </div> */}
        <p><strong>Width:</strong> <span id="PrinterWidth">{PrinterDataRef.current.Width}</span></p>
        <p><strong>Length:</strong> <span id="PrinterLength">{PrinterDataRef.current.Length}</span></p>
        <p><strong>Height:</strong> <span id="PrinterHeight">{PrinterDataRef.current.Height}</span></p>
      </div>
      <div className="PrinterPosotion">
        {/* <div style={{display:"flex", justifyContent:"center"}}>
        <h4>Position</h4>
        </div> */}
        <p><strong>X:</strong> <span id="XPos">{PrinterDataRef.current.X}</span></p>
        <p><strong>Y:</strong> <span id="YPos">{PrinterDataRef.current.Y}</span></p>
        <p><strong>Z:</strong> <span id="ZPos">{PrinterDataRef.current.Z}</span></p>
      </div>

      <div className="PrinterState">
        <p><strong>Nozzle temp:</strong> <span id="NozzleTemp">{PrinterDataRef.current.NozzleTemp}</span></p>
        <p><strong>Bed temp:</strong> <span id="BedTemp">{PrinterDataRef.current.BedTemp}</span></p>
        <p><strong>Is printing:</strong> <span id="IsPrinting">{PrinterDataRef.current.executingTask === true ? "Yes" : "No"}</span></p>
      </div>
    </div>
    )
  }